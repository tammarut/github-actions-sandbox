// package main is the entry point for the Go application
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Create a new Echo instance
	e := echo.New()

	// Configure middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 10 * time.Second,
	}))
	e.Use(middleware.CORS())
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
	}))

	// Define routes
	e.GET("/", func(c echo.Context) error {
		// log request details
		slog.Default().Info("Received request", "method", c.Request().Method, "path", c.Request().URL.Path)
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello, GitHub Actions!"})
	})

	// Add enhanced health check endpoint
	e.GET("/health", func(c echo.Context) error {
		slog.Default().Info("Health check requested", "remote_addr", c.RealIP())

		// Get environment information
		hostname, _ := os.Hostname()

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
			"version":   "0.0.3",
			"build":     "development",
			"uptime":    time.Since(time.Now().Add(-24 * time.Hour)).String(),
			"system": map[string]interface{}{
				"hostname":    hostname,
				"environment": "development",
				"go_version":  runtime.Version(),
				"goroutines":  runtime.NumGoroutine(),
			},
			"dependencies": map[string]string{
				"database": "connected", // Mock status, replace with actual DB check
				"cache":    "connected", // Mock status, replace with actual cache check
			},
		})
	})

	// Configure server
	server := &http.Server{
		Addr:              ":8080",
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           e,
	}

	// Start server in a goroutine
	go func() {
		slog.Default().Info("Server running on :8080")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			slog.Default().Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		slog.Default().Error("Failed to shutdown server gracefully", "error", err)
		fmt.Printf("Shutdown error: %v\n", err)
	} else {
		slog.Default().Info("Server shutdown gracefully")
	}
}
