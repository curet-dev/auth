// Package handler is the Vercel serverless entrypoint. Every request is
// rewritten to this function (see vercel.json) and routed by pkg/server.
package handler

import (
	"net/http"

	"github.com/curet-dev/auth/apps/api/pkg/server"
)

var app = server.New(server.ConfigFromEnv())

func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
