package usecase

import (
	"errors"
	"time"

	"github.com/lamboktulussimamora/gra-project/internal/domain/user"
)

// UserResponse represents the user data that is safe to return in API responses
type UserResponse struct {
	FirstName string
	LastName  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserUseCase defines the application use cases for user management
type UserUseCase struct {
	userRepo user.Repository
}

// NewUserUseCase creates a new user use case instance
func NewUserUseCase(repo user.Repository) *UserUseCase {
	return &UserUseCase{
		userRepo: repo,
	}
}

// Register registers a new user
func (uc *UserUseCase) Register(firstName, lastName, email, password string) (*UserResponse, error) {
	// Create a new user entity
	newUser := user.NewUser(firstName, lastName, email, password)

	// Validate user data
	if !newUser.Validate() {
		return nil, errors.New("missing required fields")
	}

	// Check if user already exists
	existingUser, _ := uc.userRepo.FindByEmail(email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// In a real application, you would:
	// 1. Hash the password
	// newUser.Password = hash(password)

	// Save user to repository
	if err := uc.userRepo.Save(newUser); err != nil {
		return nil, err
	}

	// Create response
	response := &UserResponse{
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}

	return response, nil
}
