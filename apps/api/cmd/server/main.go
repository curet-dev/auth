// Command server runs the API as a regular HTTP server for local development.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/curet-dev/auth/apps/api/pkg/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("api listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, server.New(server.ConfigFromEnv())); err != nil {
		log.Fatal(err)
	}
}
