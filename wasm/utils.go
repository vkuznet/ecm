package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"

	"golang.org/x/exp/errors"
)

// helper function to read rootCA
// https://eli.thegreenplace.net/2021/go-https-servers-with-tls/
func httpClient(rootCA string) *http.Client {
	client := &http.Client{}
	if _, err := os.Stat(rootCA); errors.Is(err, os.ErrNotExist) {
		return client
	}
	cert, err := os.ReadFile(rootCA)
	if err != nil {
		log.Println("unable to read", rootCA, "error", err)
		return client
	}
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		log.Printf("unable to parse cert from %s", rootCA)
		return client
	}
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
	}
	return client
}
