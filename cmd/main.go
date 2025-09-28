package main

import (
	"fmt"
	"log"
	"net/http"

	"surf_challenge/internal/api/router"
	"surf_challenge/internal/container"
)

func main() {
	log.Println("Starting server...")

	dependencies := container.NewAppContainer()

	mux := router.New(dependencies)

	addr := fmt.Sprintf(":%d", 3000)
	log.Printf("Listening on %s", addr)

	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

	log.Println("Server stopped.")
}
