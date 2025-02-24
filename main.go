package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/mikerybka/frontend"
	"github.com/mikerybka/twilio"
	"github.com/mikerybka/util"
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
			fmt.Println(host)
			return nil
		},
	}

	// Create the http handler
	h := &frontend.Server{
		TwilioClient: &twilio.Client{
			AccountSID:  util.RequireEnvVar("TWILIO_ACCOUNT_SID"),
			AuthToken:   util.RequireEnvVar("TWILIO_AUTH_TOKEN"),
			PhoneNumber: util.RequireEnvVar("TWILIO_PHONE_NUMBER"),
		},
		AdminPhone: util.RequireEnvVar("ADMIN_PHONE"),
		BackendURL: parsedURL,
	}

	// Create an HTTPS server using autocert
	httpsServer := &http.Server{
		Addr:      ":443",
		Handler:   h,
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
