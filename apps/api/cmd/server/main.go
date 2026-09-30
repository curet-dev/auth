// Command server runs the API as an HTTP server. Vercel's Go framework preset
// picks up cmd/server/main.go as the entrypoint and sets PORT.
package main

import (
	"bufio"
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/curet-dev/auth/apps/api/pkg/server"
)

func main() {
	loadDotEnv(".env")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app, err := server.FromEnv(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	addr := ":" + port
	log.Printf("api listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, app); err != nil {
		log.Fatal(err)
	}
}

// loadDotEnv sets KEY=VALUE pairs from path without overriding variables
// that are already set. Missing files are ignored.
func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		if _, set := os.LookupEnv(key); !set {
			os.Setenv(key, strings.Trim(strings.TrimSpace(value), `"'`))
		}
	}
}
