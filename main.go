// package main is the entry point for the Go application
package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// log request details
		fmt.Printf("Received request: %s %s\n", r.Method, r.URL.Path)
		if _, err := fmt.Fprintf(w, "Hello, GitHub Actions!"); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			fmt.Printf("Error writing response: %v\n", err)
		}
	})

	// Create a custom server with timeout settings
	server := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	fmt.Println("Server running on :8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
