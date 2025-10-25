package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dleonardoq/meli/orquestator"
)

func main() {

	httpPort := flag.String("http-port", "8080", "HTTP port")
	flag.Parse()

	logger := log.New(os.Stdout, "[MELI-API] ", log.LstdFlags|log.Lshortfile)

	logger.Println("Starting meli API")

	router := orquestator.Orquestator()

	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + *httpPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Println("Server started on", *httpPort)

	if err := srv.ListenAndServe(); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}
