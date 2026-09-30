package auth

import (
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

// UserStore persists users. MemoryStore is the only implementation for now;
// swap in a database-backed store for production.
type UserStore interface {
	Create(email, name, password string) (*User, error)
	Authenticate(email, password string) (*User, error)
}

type MemoryStore struct {
	mu      sync.RWMutex
	byEmail map[string]*User
	nextID  int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{byEmail: map[string]*User{}}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func (s *MemoryStore) Create(email, name, password string) (*User, error) {
	email = normalizeEmail(email)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
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

func (s *MemoryStore) Authenticate(email, password string) (*User, error) {
	s.mu.RLock()
	u, ok := s.byEmail[normalizeEmail(email)]
	s.mu.RUnlock()
	if !ok {
		return nil, ErrInvalidCredentials
	}
	if bcrypt.CompareHashAndPassword(u.passwordHash, []byte(password)) != nil {
		return nil, ErrInvalidCredentials
	}
	return u, nil
}
