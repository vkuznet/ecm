package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "embed"

	"golang.org/x/exp/errors"
)

// helper function to read rootCA
// https://eli.thegreenplace.net/2021/go-https-servers-with-tls/
func httpClient(rootCA string) (*http.Client, error) {
	client := &http.Client{}
	if rootCA == "" {
		return client, nil
	}
	var cert []byte
	var err error
	//     if RootCACert != nil {
	//         cert = RootCACert
	//     } else {
	if _, err := os.Stat(rootCA); errors.Is(err, os.ErrNotExist) {
		return client, err
	}
	cert, err = os.ReadFile(rootCA)
	if err != nil {
		log.Println("unable to read", rootCA, "error", err)
		return client, err
	}
	//     }
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		msg := fmt.Sprintf("unable to parse cert from %s", rootCA)
		log.Println(msg)
		return client, errors.New(msg)
	}
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
	}
	return client, nil
}
