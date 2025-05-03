// package main is the entry point for the Go application
package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// log request details
		slog.Default().Info("Received request", "method", r.Method, "path", r.URL.Path)
		if _, err := fmt.Fprintf(w, "Hello, GitHub Actions!"); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			slog.Default().Error("Error writing response", "error", err)
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

	slog.Default().Info("Server running on :8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Server error: %v\n", err)
		slog.Default().Error("Server error", "error", err)
	}
}
