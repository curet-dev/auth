package auth

import (
	"context"
	"errors"
	"os"
	"testing"
)

// Runs against a real database when TEST_DATABASE_URL is set, e.g.
// TEST_DATABASE_URL=postgres://postgres:postgres@localhost:5432/auth?sslmode=disable
func TestPostgresStore(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	ctx := context.Background()

	s, err := NewPostgresStore(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if err := s.Migrate(ctx); err != nil {
		t.Fatal(err)
	}
	// Migrations must be idempotent.
	if err := s.Migrate(ctx); err != nil {
		t.Fatalf("second migrate: %v", err)
	}
	if _, err := s.pool.Exec(ctx, "DELETE FROM users WHERE lower(email) = 'pg-test@example.com'"); err != nil {
		t.Fatal(err)
	}

	u, err := s.Create(ctx, " PG-Test@example.com ", "Ada", "supersecret")
	if err != nil {
		t.Fatal(err)
	}
	if u.ID == "" || u.Email != "pg-test@example.com" || u.CreatedAt.IsZero() {
		t.Fatalf("unexpected user %+v", u)
	}
	if _, err := s.Create(ctx, "pg-test@example.com", "Ada", "supersecret"); !errors.Is(err, ErrEmailTaken) {
		t.Fatalf("duplicate create: got %v", err)
	}
	if _, err := s.Authenticate(ctx, "pg-test@example.com", "wrong-pass"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("wrong password: got %v", err)
	}
	if _, err := s.Authenticate(ctx, "nobody@example.com", "supersecret"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("unknown user: got %v", err)
	}
	got, err := s.Authenticate(ctx, "PG-TEST@example.com", "supersecret")
	if err != nil || got.ID != u.ID {
		t.Fatalf("authenticate: got %+v, %v", got, err)
	}
}
