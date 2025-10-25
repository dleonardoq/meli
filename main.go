package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dleonardoq/meli/handlers"
	"github.com/dleonardoq/meli/middleware"
	"github.com/dleonardoq/meli/orquestator"
)

func main() {

	httpPort := flag.String("http-port", "8080", "HTTP port")
	flag.Parse()

	logger := log.New(os.Stdout, "[MELI-API] ", log.LstdFlags|log.Lshortfile)

	logger.Println("Starting meli API")
	logger.Println("Data will be saved to JSON on shutdown")

	router := orquestator.Orquestator()

	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.CorsMiddleware)

	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + *httpPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Start server in a goroutine
	go func() {
		logger.Println("Server started on port", *httpPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop

	logger.Println("\nServer shutting down...")

	// Save data to disk before shutting down
	if err := handlers.SaveDataToDisk(); err != nil {
		logger.Printf("Error saving data to disk: %v", err)
	} else {
		logger.Println("Data successfully saved to disk")
	}

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		logger.Printf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server stopped gracefully")
}
