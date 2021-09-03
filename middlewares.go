package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	limiter "github.com/ulule/limiter/v3"
	stdlib "github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
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
			rec := make(Record)
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
