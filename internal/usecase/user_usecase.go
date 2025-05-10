package usecase

import (
	"errors"
	"time"

	"github.com/lamboktulussimamora/gra-project/internal/domain/auth"
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

// AuthResponse represents the authentication response with user data and token
type AuthResponse struct {
	User  UserResponse
	Token string
}

// UserUseCase defines the application use cases for user management
type UserUseCase struct {
	userRepo        user.Repository
	passwordService auth.PasswordService
	jwtService      auth.JWTService
}

// NewUserUseCase creates a new user use case instance
func NewUserUseCase(
	repo user.Repository,
	passwordService auth.PasswordService,
	jwtService auth.JWTService,
) *UserUseCase {
	return &UserUseCase{
		userRepo:        repo,
		passwordService: passwordService,
		jwtService:      jwtService,
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

	// Hash the password
	hashedPassword, err := uc.passwordService.HashPassword(password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	newUser.Password = hashedPassword

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

// Login authenticates a user and returns a token
func (uc *UserUseCase) Login(email, password string) (*AuthResponse, error) {
	// Find user by email
	user, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	valid, err := uc.passwordService.VerifyPassword(user.Password, password)
	if err != nil || !valid {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := uc.jwtService.GenerateToken(user)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Create user response
	userResp := UserResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// Return auth response with token
	return &AuthResponse{
		User:  userResp,
		Token: token,
	}, nil
}
