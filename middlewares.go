package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/context"
	limiter "github.com/ulule/limiter/v3"
	stdlib "github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
	vt "github.com/vkuznet/ecm/vault"
)

// limiter middleware pointer
var limiterMiddleware *stdlib.Middleware

// initialize Limiter middleware pointer
func initLimiter(period string) {
	log.Printf("limiter rate='%s'", period)
	// create rate limiter with 5 req/second
	rate, err := limiter.NewRateFromFormatted(period)
	if err != nil {
		panic(err)
	}
	store := memory.NewStore()
	instance := limiter.New(store, rate)
	limiterMiddleware = stdlib.NewMiddleware(instance)
}

// helper function to check auth/authz headers
func checkAuthnAuthz(header http.Header) bool {
	return true
}

// helper to auth/authz incoming requests to the server
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// perform authentication
		status := checkAuthnAuthz(r.Header)
		if !status {
			log.Printf("ERROR: fail to authenticate, HTTP headers %+v\n", r.Header)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if ServerConfig.Verbose > 2 {
			log.Printf("Auth layer status: %v headers: %+v\n", status, r.Header)
		}
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func validate(r *http.Request) error {
	return nil
}

// helper to validate incoming requests' parameters
func validateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			next.ServeHTTP(w, r)
			return
		}
		// perform validation of input parameters
		err := validate(r)
		if err != nil {
			uri, e := url.QueryUnescape(r.RequestURI)
			if e == nil {
				log.Printf("HTTP %s %s validation error %v\n", r.Method, uri, err)
			} else {
				log.Printf("HTTP %s %v validation error %v\n", r.Method, r.RequestURI, err)
			}
			w.WriteHeader(http.StatusBadRequest)
			rec := make(vt.Record)
			rec["error"] = fmt.Sprintf("Validation error %v", err)
			if r, e := json.Marshal(rec); e == nil {
				w.Write(r)
			}
			return
		}
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// limit middleware limits incoming requests
func limitMiddleware(next http.Handler) http.Handler {
	return limiterMiddleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	}))
}

// cors middleware provide CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			//             origin := fmt.Sprintf("http://localhost:%d/*", ServerConfig.Port)
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers:", "Origin, Content-Type, X-Auth-Token, Authorization")
			w.Header().Set("Content-Type", "application/json")
		}
		log.Println("call next ServeHTTP")
		next.ServeHTTP(w, r)
	})
}

// ValidateMiddleware provides authentication of user credentials
func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//         log.Printf("#### ValidateMiddleware request header %+v", r.Header)
		var user string
		if r.Method == "GET" {
			query := r.URL.Query()
			user = query.Get("user") //user="bla"
			if user == "" {
				w.Write([]byte("Please provide user name"))
			}
		} else if r.Method == "POST" {
			err := r.ParseForm()
			if err == nil {
				user = r.FormValue("user")
			} else {
				var otpToken OtpToken
				err := json.NewDecoder(r.Body).Decode(&otpToken)
				if err != nil {
					log.Println("error", err)
					json.NewEncoder(w).Encode(err)
					return
				}
				user = otpToken.User
			}
		}
		secret := findUserSecret(user)
		if secret == "" {
			err := errors.New("Unable to find user credentials, please obtain proper QR code")
			log.Println(err.Error())
			json.NewEncoder(w).Encode(err)
			return
		}
		bearerToken, err := getBearerToken(r.Header.Get("authorization"))
		if err != nil {
			log.Println("unable to get authorization token", err)
			json.NewEncoder(w).Encode(err)
			return
		}
		// for verification we can use either user's secret
		// or server secret
		// in latter case it should be global and available to all APIs
		decodedToken, err := VerifyJwt(bearerToken, secret)
		if err != nil {
			log.Println("unable to verify token", err)
			json.NewEncoder(w).Encode(err)
			return
		}
		if decodedToken["authorized"] == true {
			context.Set(r, "decoded", decodedToken)
			next(w, r)
		} else {
			json.NewEncoder(w).Encode("2FA is required")
		}
	})
}
