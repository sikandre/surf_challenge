package main

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
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
	readTimeout       = 15 * time.Second
	readHeaderTimeout = 5 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
	headerBytes       = 1 << 20
	duration          = 20 * time.Second
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	sugar := logger.Sugar()
	defer logger.Sync()

	sugar.Info("Logger initialized")

	dependencies := container.NewAppContainer(sugar)

	mux := router.New(sugar, dependencies)

	addr := fmt.Sprintf(":%d", port)

	sugar.Infof("Starting server on %s", addr)

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
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

	sugar.Info("Received interrupt signal, shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exiting")
}
