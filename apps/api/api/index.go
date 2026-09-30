// Package handler is the Vercel serverless entrypoint. Every request is
// rewritten to this function (see vercel.json) and routed by pkg/server.
package handler

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/curet-dev/auth/apps/api/pkg/server"
)

var (
	app     *server.Server
	initErr error
	once    sync.Once
)

// Handler lazily initialises the server (DB pool + migrations) once per
// function instance and reuses it for subsequent invocations.
func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(func() {
		app, initErr = server.FromEnv(context.Background())
		if initErr != nil {
			log.Printf("init: %v", initErr)
		}
	})
	if initErr != nil {
		http.Error(w, `{"error":"server misconfigured"}`, http.StatusInternalServerError)
		return
	}
	app.ServeHTTP(w, r)
}
