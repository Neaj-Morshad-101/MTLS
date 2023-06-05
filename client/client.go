// client/client.go

package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	cert, err := ioutil.ReadFile("certs/ca.crt")
	if err != nil {
		log.Fatalf("Could not open certificate file: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)

	certificate, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")

	if err != nil {
		log.Fatalf("Could not load certificate: %v", err)
	}

	client := http.Client{
		Timeout: time.Minute * 3,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{certificate},
			},
		},
	}

	// change the address to match the common name of the certificate
	resp, err := client.Get("https://localhost:9090")
	if err != nil {
		log.Fatalf("error making get request: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading response: %v", err)
	}
	fmt.Println(string(body))
}
