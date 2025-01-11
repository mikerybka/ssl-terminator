package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"golang.org/x/crypto/acme/autocert"
)

func main() {
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

	// Create a reverse proxy to forward requests to the backend
	reverseProxy := httputil.NewSingleHostReverseProxy(parsedURL)

	// Define the handler that will use the reverse proxy
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Modify the request if needed, e.g., change headers
		r.Host = parsedURL.Host
		reverseProxy.ServeHTTP(w, r)
	})

	// Create a certificate manager
	certManager := autocert.Manager{
		Cache:      autocert.DirCache("certs"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(), // Set allowed hostnames here
	}

	// Create an HTTPS server using autocert
	httpsServer := &http.Server{
		Addr:      ":443",
		Handler:   handler,
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
