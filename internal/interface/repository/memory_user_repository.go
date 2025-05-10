package repository

import (
	"errors"
	"sync"

	"github.com/lamboktulussimamora/gra-project/internal/domain/user"
)

// InMemoryUserRepository is an in-memory implementation of the user repository
type InMemoryUserRepository struct {
	users map[string]*user.User
	mu    sync.RWMutex
}

// NewInMemoryUserRepository creates a new in-memory user repository
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*user.User),
	}
}

// Save saves a user in the in-memory store
func (r *InMemoryUserRepository) Save(user *user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if user already exists
	if _, exists := r.users[user.Email]; exists {
		return errors.New("user already exists")
	}

	// Store the user
	r.users[user.Email] = user
	return nil
}

// FindByEmail finds a user by email
func (r *InMemoryUserRepository) FindByEmail(email string) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[email]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}
