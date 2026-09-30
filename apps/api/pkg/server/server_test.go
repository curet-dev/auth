package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/curet-dev/auth/apps/api/pkg/auth"
)

func do(t *testing.T, s *Server, method, path, body string, cookies ...*http.Cookie) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	for _, c := range cookies {
		req.AddCookie(c)
	}
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)
	return rec
}

func TestAuthFlow(t *testing.T) {
	s := New(Config{JWTSecret: "test", SessionTTL: time.Hour}, auth.NewMemoryStore())

	if rec := do(t, s, "GET", "/api/auth/me", ""); rec.Code != http.StatusUnauthorized {
		t.Fatalf("me without session: got %d", rec.Code)
	}

	body := `{"email":"Ada@example.com","password":"supersecret","name":"Ada"}`
	rec := do(t, s, "POST", "/api/auth/register", body)
	if rec.Code != http.StatusCreated {
		t.Fatalf("register: got %d %s", rec.Code, rec.Body)
	}
	if rec := do(t, s, "POST", "/api/auth/register", body); rec.Code != http.StatusConflict {
		t.Fatalf("duplicate register: got %d", rec.Code)
	}

	if rec := do(t, s, "POST", "/api/auth/login", `{"email":"ada@example.com","password":"wrong-pass"}`); rec.Code != http.StatusUnauthorized {
		t.Fatalf("login with wrong password: got %d", rec.Code)
	}
	rec = do(t, s, "POST", "/api/auth/login", `{"email":"ada@example.com","password":"supersecret"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("login: got %d %s", rec.Code, rec.Body)
	}
	cookies := rec.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Name != sessionCookie || !cookies[0].HttpOnly {
		t.Fatalf("login: unexpected cookies %v", cookies)
	}

	rec = do(t, s, "GET", "/api/auth/me", "", cookies[0])
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"email":"ada@example.com"`) {
		t.Fatalf("me: got %d %s", rec.Code, rec.Body)
	}
}

func TestRegisterValidation(t *testing.T) {
	s := New(Config{JWTSecret: "test", SessionTTL: time.Hour}, auth.NewMemoryStore())
	for _, body := range []string{
		`{"email":"nope","password":"supersecret","name":"Ada"}`,
		`{"email":"a@example.com","password":"short","name":"Ada"}`,
		`{"email":"a@example.com","password":"supersecret","name":" "}`,
		`not json`,
	} {
		if rec := do(t, s, "POST", "/api/auth/register", body); rec.Code != http.StatusBadRequest {
			t.Errorf("%s: got %d", body, rec.Code)
		}
	}
}
