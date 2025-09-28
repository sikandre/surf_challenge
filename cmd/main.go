package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"surf_challenge/internal/api/router"
	"surf_challenge/internal/container"
)

const (
	port              = 3000
	ReadTimeout       = 15 * time.Second
	ReadHeaderTimeout = 5 * time.Second
	WriteTimeout      = 15 * time.Second
	IdleTimeout       = 60 * time.Second
	headerBytes       = 1 << 20
	duration          = 20 * time.Second
)

func main() {
	log.Println("Starting server...")

	dependencies := container.NewAppContainer()

	mux := router.New(dependencies)

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Listening on %s", addr)

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadTimeout:       ReadTimeout,
		ReadHeaderTimeout: ReadHeaderTimeout,
		WriteTimeout:      WriteTimeout,
		IdleTimeout:       IdleTimeout,
		MaxHeaderBytes:    headerBytes,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
