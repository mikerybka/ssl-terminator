package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"golang.org/x/crypto/acme/autocert"
)

func main() {
	// Fetch certDir from env var
	certDir := os.Getenv("CERT_DIR")
	if certDir == "" {
		log.Fatal("CERT_DIR environment variable is not set")
	}

	// Fetch BACKEND_URL from environment variable
	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		log.Fatal("BACKEND_URL environment variable is not set")
	}

	// Parse the backend URL
	parsedURL, err := url.Parse(backendURL)
	if err != nil {
		log.Fatalf("Invalid BACKEND_URL: %v", err)
	}

	// Create a certificate manager
	certManager := autocert.Manager{
		Cache:  autocert.DirCache(certDir),
		Prompt: autocert.AcceptTOS,
		HostPolicy: func(ctx context.Context, host string) error {
			return nil
		},
	}

	// Create an HTTPS server using autocert
	httpsServer := &http.Server{
		Addr:      ":443",
		Handler:   httputil.NewSingleHostReverseProxy(parsedURL),
		TLSConfig: certManager.TLSConfig(),
	}

	// Redirect HTTP to HTTPS
	go func() {
		http.ListenAndServe(":80", certManager.HTTPHandler(nil))
	}()

	// Start the HTTPS server
	log.Println("Starting server on :443")
	if err := httpsServer.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("Failed to start HTTPS server: %v", err)
	}
}
