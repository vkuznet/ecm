package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"embed"

	"golang.org/x/exp/errors"
)

//go:embed certificates
var certDir embed.FS

// helper function to read rootCA
// https://eli.thegreenplace.net/2021/go-https-servers-with-tls/
func httpClient() (*http.Client, error) {
	client := &http.Client{}
	certPool := x509.NewCertPool()

	// read rootCAs from certificates area
	if entries, err := certDir.ReadDir("certificates"); err == nil {
		for _, f := range entries {
			rootCA := f.Name()
			if _, err := os.Stat(rootCA); errors.Is(err, os.ErrNotExist) {
				return client, err
			}
			certPEM, err := os.ReadFile(rootCA)
			if err != nil {
				log.Println("unable to read", rootCA, "error", err)
				continue
			}
			// before adding new rootCA I need to check if it is still valid
			cert, err := x509.ParseCertificate(certPEM)
			today := time.Now()
			if today.Before(cert.NotAfter) && today.After(cert.NotBefore) {
				if ok := certPool.AppendCertsFromPEM(certPEM); !ok {
					msg := fmt.Sprintf("unable to parse cert from %s", rootCA)
					log.Println(msg)
				}
			}
		}
	} else {
		msg := fmt.Sprintf("unable to read rootCAs from certificates area %v", err)
		log.Println(msg)
		return client, errors.New(msg)
	}

	// setup HTTP client with content of rootCAs
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
	}
	return client, nil
}
