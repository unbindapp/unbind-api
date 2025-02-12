package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/unbindapp/unbind-api/internal/generated"
	"github.com/unbindapp/unbind-api/internal/server"
)

func main() {
	// Implementation
	srvImpl := &server.Server{}

	// New chi router
	r := chi.NewRouter()

	// Register the routes from generated code
	generated.HandlerFromMux(srvImpl, r)

	// Start the server
	addr := ":8081"
	fmt.Printf("Starting server on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
