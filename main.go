package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	version   = "1.0.1"
	startTime = time.Now()
)

func main() {
	appVersion := getEnv("APP_VERSION", version)
	appEnv := getEnv("APP_ENV", "development")
	port := getEnv("PORT", "8080")

	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	// Routes
	e.GET("/", func(c echo.Context) error {
		hostname, _ := os.Hostname()
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":     "Hello from Go Sample API",
			"version":     appVersion,
			"hostname":    hostname,
			"environment": appEnv,
			"timestamp":   time.Now().UTC().Format(time.RFC3339),
			"uptime":      time.Since(startTime).String(),
		})
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
		})
	})

	e.GET("/ready", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ready",
		})
	})

	// Graceful shutdown
	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatalf("shutting down the server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	e.Logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	e.Logger.Info("Server stopped gracefully")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
