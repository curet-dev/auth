package auth

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"createdAt"`
	passwordHash []byte
}

// UserStore persists users. PostgresStore is used whenever DATABASE_URL is
// set; MemoryStore is a fallback for tests and quick local runs.
type UserStore interface {
	Create(ctx context.Context, email, name, password string) (*User, error)
	Authenticate(ctx context.Context, email, password string) (*User, error)
	Ping(ctx context.Context) error
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func checkPassword(u *User, password string) (*User, error) {
	if bcrypt.CompareHashAndPassword(u.passwordHash, []byte(password)) != nil {
		return nil, ErrInvalidCredentials
	}
	return u, nil
}

type MemoryStore struct {
	mu      sync.RWMutex
	byEmail map[string]*User
	nextID  int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{byEmail: map[string]*User{}}
}

func (s *MemoryStore) Create(_ context.Context, email, name, password string) (*User, error) {
	email = normalizeEmail(email)
	hash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.byEmail[email]; ok {
		return nil, ErrEmailTaken
	}
	s.nextID++
	u := &User{
		ID:           "usr_" + strconv.Itoa(s.nextID),
		Email:        email,
		Name:         strings.TrimSpace(name),
		CreatedAt:    time.Now().UTC(),
		passwordHash: hash,
	}
	s.byEmail[email] = u
	return u, nil
}

func (s *MemoryStore) Authenticate(_ context.Context, email, password string) (*User, error) {
	s.mu.RLock()
	u, ok := s.byEmail[normalizeEmail(email)]
	s.mu.RUnlock()
	if !ok {
		return nil, ErrInvalidCredentials
	}
	return checkPassword(u, password)
}

func (s *MemoryStore) Ping(context.Context) error { return nil }
