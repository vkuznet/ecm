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
	"path/filepath"
	"strings"
	"syscall"
	"time"

	_ "expvar"         // to be used for monitoring, see https://github.com/divan/expvarmon
	_ "net/http/pprof" // profiler, see https://golang.org/pkg/net/http/pprof/

	"github.com/dchest/captcha"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	kvdb "github.com/vkuznet/ecm/kvdb"
	logging "github.com/vkuznet/http-logging"
	"golang.org/x/crypto/acme/autocert"
)

// ServerConfiguration stores server configuration parameters
type ServerConfiguration struct {
	Port          int      `json:"port"`         // server port number
	Base          string   `json:"base"`         // base URL
	Verbose       int      `json:"verbose"`      // verbose output
	ServerCrt     string   `json:"serverCrt"`    // path to server crt file
	ServerKey     string   `json:"serverKey"`    // path to server key file
	RootCA        string   `json:"rootCA"`       // RootCA file
	CSRFKey       string   `json:"csrfKey"`      // CSRF 32-byte-long-auth-key
	Production    bool     `json:"production"`   // production server or not
	VaultArea     string   `json:"vault_area"`   // vault directory
	LimiterPeriod string   `json:"rate"`         // limiter rate value
	LogFile       string   `json:"log_file"`     // server log file
	LetsEncrypt   bool     `json:"lets_encrypt"` // start LetsEncrypt HTTPs server
	DomainNames   []string `json:"domain_names"` // list of domain names to use
	StaticDir     string   `json:"static"`       // location of static files
	Templates     string   `json:"templates"`    // server templates
	DBStore       string   `json:"dbstore"`      // location of dbstore
	DevelopMode   bool     `json:"develop_mode"` // development mode
}

// ServerConfig variable represents configuration object
var ServerConfig ServerConfiguration

// DBStore represents user data store
var DBStore *kvdb.Store

// helper function to parse configuration
func parseServerConfig(configFile string) error {
	path, err := os.Getwd()
	if err != nil {
		log.Println("unable to get current directory", err)
		path = "."
	}
	ServerConfig.StaticDir = fmt.Sprintf("%s/static", path)
	ServerConfig.Templates = fmt.Sprintf("%s/static/tmpl", path)

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

func srvRouter() *mux.Router {
	router := mux.NewRouter()
	//     router.StrictSlash(true) // to allow /route and /route/ end-points
	router.HandleFunc(basePath("/vault/{vault:[0-9a-zA-Z]+}/auth"), VaultAuthHandler).Methods("POST")
	router.HandleFunc(basePath("/vault/{vault:[0-9a-zA-Z]+}/records"), VaultRecordsHandler).Methods("GET")
	router.HandleFunc(basePath("/vault/{vault:[0-9a-zA-Z]+}"), VaultHandler).Methods("GET")
	router.HandleFunc(basePath("/vault/{vault:[0-9a-zA-Z]+}/{rid:[0-9a-zA-Z-\\.]+}"), VaultRecordHandler).Methods("GET", "POST")
	router.HandleFunc(basePath("/vault/{vault:[0-9a-zA-Z]+}/{rid:[0-9a-zA-Z-]+}"), VaultDeleteHandler).Methods("DELETE")
	router.HandleFunc(basePath("/vault/{vault:[0-9a-zA-Z]+}"), VaultAddHandler).Methods("POST")
	router.HandleFunc(basePath("/token"), TokenHandler).Methods("GET")
	router.HandleFunc(basePath("/favicon.ico"), FaviconHandler)

	router.HandleFunc(basePath("/signup"), SignUpHandler).Methods("GET")
	router.HandleFunc(basePath("/login"), LoginHandler).Methods("GET")
	router.HandleFunc(basePath("/logout"), LogoutHandler).Methods("GET")
	router.HandleFunc("/authenticate", AuthHandler).Methods("POST")
	router.HandleFunc("/verify", VerifyHandler).Methods("POST")
	if ServerConfig.DevelopMode {
		router.HandleFunc(basePath("/main"), MainHandler).Methods("GET", "POST")
	} else {
		router.HandleFunc(basePath("/main"), ValidateMiddleware(MainHandler)).Methods("GET", "POST")
	}
	router.HandleFunc(basePath("/user"), UserHandler).Methods("GET", "POST")
	router.HandleFunc(basePath("/qrcode"), QRHandler).Methods("GET", "POST")
	router.HandleFunc(basePath("/"), HomeHandler).Methods("GET")

	// this is for displaying the QR code on /qr end point
	// and static area which holds user's images
	log.Println("server static area", ServerConfig.StaticDir)
	fileServer := http.StripPrefix("/static/", http.FileServer(http.Dir(ServerConfig.StaticDir)))
	router.PathPrefix(basePath("/static/{user:[0-9a-zA-Z-]+}/{file:[0-9a-zA-Z-\\.]+}")).Handler(fileServer)

	// static css content
	router.PathPrefix(basePath("/css/{file:[0-9a-zA-Z-\\.]+}")).Handler(fileServer)
	router.PathPrefix(basePath("/js/{file:[0-9a-zA-Z-\\.]+}")).Handler(fileServer)
	router.PathPrefix(basePath("/images/{file:[0-9a-zA-Z-\\.]+}")).Handler(fileServer)
	router.PathPrefix(basePath("/fonts/{file:[0-9a-zA-Z-\\.]+}")).Handler(fileServer)

	// add captcha server
	captchaServer := captcha.Server(captcha.StdWidth, captcha.StdHeight)
	router.PathPrefix(basePath("/captcha/")).Handler(captchaServer)

	// for all requests
	router.Use(logging.LoggingMiddleware)
	// for all requests perform first auth/authz action
	router.Use(authMiddleware)
	// validate all input parameters
	router.Use(validateMiddleware)
	// use limiter middleware to slow down clients
	router.Use(limitMiddleware)
	// use cors middleware
	//     router.Use(corsMiddleware)

	return router
}

// http server implementation
func server(serverCrt, serverKey string) {

	// setup our DB store
	_, err := os.Stat(ServerConfig.DBStore)
	if os.IsNotExist(err) {
		err = os.MkdirAll(ServerConfig.DBStore, 0755)
		if err != nil {
			log.Fatalf("unable to create new KV store, error %v", err)
		}
	}
	store, err := kvdb.NewStore(ServerConfig.DBStore)
	if err != nil {
		log.Fatalf("unable to init KV store, error %v", err)
	}
	defer store.Close()
	DBStore = store

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

		http.Handle("/", CSRF(srvRouter()))
	} else {
		http.Handle("/", srvRouter())
	}
	// define our HTTP server
	srv := getServer()

	// make extra channel for graceful shutdown
	// https://medium.com/honestbee-tw-engineer/gracefully-shutdown-in-go-http-server-5f5e6b83da5a
	httpDone := make(chan os.Signal, 1)
	signal.Notify(httpDone, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		var err error
		if serverCrt != "" && serverKey != "" {
			//start HTTPS server
			log.Printf("Starting HTTPs server :%d", ServerConfig.Port)
			err = srv.ListenAndServeTLS(ServerConfig.ServerCrt, ServerConfig.ServerKey)
		} else if ServerConfig.LetsEncrypt {
			//start LetsEncrypt HTTPS server
			log.Printf("Starting LetsEncrypt HTTPs server :%d", ServerConfig.Port)
			err = srv.ListenAndServeTLS("", "")
		} else {
			// Start server without user certificates
			log.Printf("Starting HTTP server :%d", ServerConfig.Port)
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

// helper function to return http.Server object for different configurations
func getServer() *http.Server {
	srvCrt := ServerConfig.ServerCrt
	srvKey := ServerConfig.ServerKey
	port := ServerConfig.Port
	verbose := ServerConfig.Verbose
	rootCAs := ServerConfig.RootCA
	var srv *http.Server
	if srvCrt != "" && srvKey != "" {
		srv = TlsServer(srvCrt, srvKey, rootCAs, port, verbose)
	} else if ServerConfig.LetsEncrypt {
		srv = LetsEncryptServer(ServerConfig.DomainNames...)
	} else {
		addr := fmt.Sprintf(":%d", port)
		srv = &http.Server{
			Addr: addr,
		}
	}
	return srv
}

// LetsEncryptServer provides HTTPs server with Let's encrypt for
// given domain names (hosts)
func LetsEncryptServer(hosts ...string) *http.Server {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(hosts...),
		Cache:      autocert.DirCache("certs"),
	}

	server := &http.Server{
		Addr: ":https",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}
	// start cert Manager goroutine
	go http.ListenAndServe(":http", certManager.HTTPHandler(nil))
	return server
}

// TlsServer returns TLS enabled HTTP server
func TlsServer(serverCrt, serverKey, rootCAs string, port, verbose int) *http.Server {
	var certPool *x509.CertPool
	if rootCAs != "" {
		certPool := x509.NewCertPool()
		files, err := ioutil.ReadDir(rootCAs)
		if err != nil {
			log.Fatal(err)
			log.Fatalf("Unable to list files in '%s', error: %v\n", rootCAs, err)
		}
		for _, finfo := range files {
			fname := fmt.Sprintf("%s/%s", rootCAs, finfo.Name())
			caCert, err := os.ReadFile(filepath.Clean(fname))
			if err != nil {
				if verbose > 1 {
					log.Printf("Unable to read %s\n", fname)
				}
			}
			if ok := certPool.AppendCertsFromPEM(caCert); !ok {
				if verbose > 1 {
					log.Printf("invalid PEM format while importing trust-chain: %q", fname)
				}
			}
		}
	}
	// if we do not require custom verification we'll load server crt/key and present to client
	cert, err := tls.LoadX509KeyPair(serverCrt, serverKey)
	if err != nil {
		log.Fatalf("server loadkeys: %s", err)
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	if certPool != nil {
		tlsConfig.RootCAs = certPool
	}
	addr := fmt.Sprintf(":%d", port)
	server := &http.Server{
		Addr:      addr,
		TLSConfig: tlsConfig,
	}
	return server
}
