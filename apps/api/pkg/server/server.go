// Package server contains the HTTP routes of the API. It is shared by the
// local dev server (cmd/server) and the Vercel function (api/index.go).
package server

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/curet-dev/auth/apps/api/pkg/auth"
)

const sessionCookie = "session"

type Config struct {
	JWTSecret   string
	SessionTTL  time.Duration
	DatabaseURL string
}

func ConfigFromEnv() Config {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("WARNING: JWT_SECRET is not set, using an insecure development secret")
		secret = "dev-secret-change-me"
	}
	return Config{
		JWTSecret:   secret,
		SessionTTL:  7 * 24 * time.Hour,
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}

// FromEnv builds a server from environment variables. With DATABASE_URL set
// users are stored in Postgres, otherwise in memory (not allowed on Vercel,
// where every cold start would lose all users).
func FromEnv(ctx context.Context) (*Server, error) {
	cfg := ConfigFromEnv()
	if cfg.DatabaseURL == "" {
		if os.Getenv("VERCEL") != "" {
			return nil, errors.New("DATABASE_URL is not set")
		}
		log.Println("WARNING: DATABASE_URL is not set, users are stored in memory")
		return New(cfg, auth.NewMemoryStore()), nil
	}

	store, err := auth.NewPostgresStore(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	if err := store.Migrate(ctx); err != nil {
		store.Close()
		return nil, err
	}
	return New(cfg, store), nil
}

type Server struct {
	users  auth.UserStore
	tokens *auth.Tokens
	mux    *http.ServeMux
}

func New(cfg Config, users auth.UserStore) *Server {
	s := &Server{
		users:  users,
		tokens: auth.NewTokens(cfg.JWTSecret, cfg.SessionTTL),
		mux:    http.NewServeMux(),
	}
	s.mux.HandleFunc("GET /api/health", s.health)
	s.mux.HandleFunc("POST /api/auth/register", s.register)
	s.mux.HandleFunc("POST /api/auth/login", s.login)
	s.mux.HandleFunc("POST /api/auth/logout", s.logout)
	s.mux.HandleFunc("GET /api/auth/me", s.me)
	s.mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		writeError(w, http.StatusNotFound, "not found")
	})
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	if err := s.users.Ping(r.Context()); err != nil {
		log.Printf("health: database ping failed: %v", err)
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "database unavailable"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	var in credentials
	if !decode(w, r, &in) {
		return
	}
	if _, err := mail.ParseAddress(in.Email); err != nil {
		writeError(w, http.StatusBadRequest, "please enter a valid email address")
		return
	}
	if len(in.Password) < 8 {
		writeError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}
	if strings.TrimSpace(in.Name) == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	u, err := s.users.Create(r.Context(), in.Email, in.Name, in.Password)
	if errors.Is(err, auth.ErrEmailTaken) {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	if err != nil {
		log.Printf("register: %v", err)
		writeError(w, http.StatusInternalServerError, "could not create user")
		return
	}
	s.startSession(w, r, u, http.StatusCreated)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var in credentials
	if !decode(w, r, &in) {
		return
	}
	u, err := s.users.Authenticate(r.Context(), in.Email, in.Password)
	if errors.Is(err, auth.ErrInvalidCredentials) {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if err != nil {
		log.Printf("login: %v", err)
		writeError(w, http.StatusInternalServerError, "could not sign in")
		return
	}
	s.startSession(w, r, u, http.StatusOK)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   isHTTPS(r),
		SameSite: http.SameSiteLaxMode,
	})
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) me(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(sessionCookie)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "not signed in")
		return
	}
	claims, err := s.tokens.Parse(c.Value)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "session expired")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"user": map[string]string{"id": claims.Subject, "email": claims.Email, "name": claims.Name},
	})
}

func (s *Server) startSession(w http.ResponseWriter, r *http.Request, u *auth.User, status int) {
	token, err := s.tokens.Issue(u)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create session")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    token,
		Path:     "/",
		MaxAge:   int(s.tokens.TTL().Seconds()),
		HttpOnly: true,
		Secure:   isHTTPS(r),
		SameSite: http.SameSiteLaxMode,
	})
	writeJSON(w, status, map[string]any{"user": u})
}

func isHTTPS(r *http.Request) bool {
	return r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
