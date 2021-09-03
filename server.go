package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "expvar"         // to be used for monitoring, see https://github.com/divan/expvarmon
	_ "net/http/pprof" // profiler, see https://golang.org/pkg/net/http/pprof/

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	logging "github.com/vkuznet/http-logging"
)

// ServerConfiguration stores server configuration parameters
type ServerConfiguration struct {
	Port          int    `json:"port"`       // server port number
	Base          string `json:"base"`       // base URL
	Verbose       int    `json:"verbose"`    // verbose output
	ServerCrt     string `json:"serverCrt"`  // path to server crt file
	ServerKey     string `json:"serverKey"`  // path to server key file
	RootCA        string `json:"rootCA"`     // RootCA file
	CSRFKey       string `json:"csrfKey"`    // CSRF 32-byte-long-auth-key
	Production    bool   `json:"production"` // production server or not
	VaultArea     string `json:"vault_area"` // vault directory
	LimiterPeriod string `json:"rate"`       // limiter rate value
	LogFile       string `json:"log_file"`   // server log file
}

// ServerConfig variable represents configuration object
var ServerConfig ServerConfiguration

// helper function to parse configuration
func parseServerConfig(configFile string) error {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println("Unable to read", err)
		return err
	}
	err = json.Unmarshal(data, &ServerConfig)
	if err != nil {
		log.Println("Unable to parse", err)
		return err
	}
	if ServerConfig.LimiterPeriod == "" {
		ServerConfig.LimiterPeriod = "100-S"
	}
	if ServerConfig.Port == 0 {
		ServerConfig.Port = 8888
	}
	return nil
}

func basePath(api string) string {
	if ServerConfig.Base != "" {
		if strings.HasPrefix(api, "/") {
			api = strings.Replace(api, "/", "", 1)
		}
		if strings.HasPrefix(ServerConfig.Base, "/") {
			return fmt.Sprintf("%s/%s", ServerConfig.Base, api)
		}
		return fmt.Sprintf("/%s/%s", ServerConfig.Base, api)
	}
	return api
}

func handlers() *mux.Router {
	router := mux.NewRouter()
	//     router.StrictSlash(true) // to allow /route and /route/ end-points
	router.HandleFunc(basePath("/vault/{vault:[0-9a-zA-Z]+}"), VaultHandler).Methods("GET")
	router.HandleFunc(basePath("/vault/{vault:[0-9a-zA-Z]+}/{rid:[0-9a-zA-Z-]+}"), VaultRecordHandler).Methods("GET")
	router.HandleFunc(basePath("/vault/{vault:[0-9a-zA-Z]+}/{rid:[0-9a-zA-Z-]+}"), VaultDeleteHandler).Methods("DELETE")
	router.HandleFunc(basePath("/vault/{vault:[0-9a-zA-Z]+}"), VaultAddHandler).Methods("POST")
	router.HandleFunc(basePath("/token"), TokenHandler).Methods("GET")
	// for all requests
	router.Use(logging.LoggingMiddleware)
	// for all requests perform first auth/authz action
	router.Use(authMiddleware)
	// validate all input parameters
	router.Use(validateMiddleware)
	// use limiter middleware to slow down clients
	router.Use(limitMiddleware)

	return router
}

// http server implementation
func server(serverCrt, serverKey string) {
	// define server hand	// dynamic handlers
	if ServerConfig.CSRFKey != "" {
		CSRF := csrf.Protect(
			[]byte(ServerConfig.CSRFKey),
			csrf.RequestHeader("Authenticity-Token"),
			csrf.FieldName("authenticity_token"),
			csrf.Secure(ServerConfig.Production),
			csrf.ErrorHandler(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					log.Printf("### CSRF error handler: %+v\n", r)
					w.WriteHeader(http.StatusForbidden)
				},
			)),
		)

		http.Handle("/", CSRF(handlers()))
	} else {
		http.Handle("/", handlers())
	}
	// define our HTTP server
	addr := fmt.Sprintf(":%d", ServerConfig.Port)
	srv := &http.Server{
		Addr: addr,
	}

	// make extra channel for graceful shutdown
	// https://medium.com/honestbee-tw-engineer/gracefully-shutdown-in-go-http-server-5f5e6b83da5a
	httpDone := make(chan os.Signal, 1)
	signal.Notify(httpDone, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		var err error
		if serverCrt != "" && serverKey != "" {
			//start HTTPS server which require user certificates
			rootCA := x509.NewCertPool()
			caCert, _ := ioutil.ReadFile(ServerConfig.RootCA)
			rootCA.AppendCertsFromPEM(caCert)
			srv = &http.Server{
				Addr: addr,
				TLSConfig: &tls.Config{
					//                 ClientAuth: tls.RequestClientCert,
					RootCAs: rootCA,
				},
			}
			log.Println("Starting HTTPs server", addr)
			err = srv.ListenAndServeTLS(ServerConfig.ServerCrt, ServerConfig.ServerKey)
		} else {
			// Start server without user certificates
			log.Println("Starting HTTP server", addr)
			err = srv.ListenAndServe()
		}
		if err != nil {
			log.Printf("Fail to start server %v", err)
		}
	}()

	// properly stop our HTTP and Migration Servers
	<-httpDone
	log.Print("HTTP server stopped")

	// add extra timeout for shutdown service stuff
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("HTTP server exited properly")
}

// helper function to start the web server
func startServer(config string) {
	err := parseServerConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	log.SetFlags(0)
	if ServerConfig.Verbose > 0 {
		log.SetFlags(log.Lshortfile)
	}
	log.SetOutput(new(logging.LogWriter))
	if ServerConfig.LogFile != "" {
		rl, err := rotatelogs.New(ServerConfig.LogFile + "-%Y%m%d")
		if err == nil {
			rotlogs := logging.RotateLogWriter{RotateLogs: rl}
			log.SetOutput(rotlogs)
		}
	}

	initLimiter(ServerConfig.LimiterPeriod)
	if err != nil {
		log.Fatalf("Unable to parse config file %s, error: %v", config, err)
	}
	_, e1 := os.Stat(ServerConfig.ServerCrt)
	_, e2 := os.Stat(ServerConfig.ServerKey)
	if e1 == nil && e2 == nil {
		server(ServerConfig.ServerCrt, ServerConfig.ServerKey)
	} else {
		server("", "")
	}
}
