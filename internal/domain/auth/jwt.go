package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lamboktulussimamora/gra-project/internal/domain/user"
)

// Common errors
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// JWTService provides methods to generate and validate JWT tokens
type JWTService interface {
	GenerateToken(user *user.User) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

// Claims represents the JWT claims
type Claims struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	jwt.RegisteredClaims
}

// JWTConfig holds JWT configuration parameters
type JWTConfig struct {
	SecretKey     string
	TokenDuration time.Duration
}

// DefaultJWTService is a default implementation of JWTService
type DefaultJWTService struct {
	config JWTConfig
}

// NewJWTService creates a new JWT service
func NewJWTService(config JWTConfig) JWTService {
	return &DefaultJWTService{
		config: config,
	}
}

// GenerateToken generates a new JWT token for the given user
func (s *DefaultJWTService) GenerateToken(user *user.User) (string, error) {
	claims := Claims{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.TokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.SecretKey))
}

// ValidateToken validates the provided token and returns the claims
func (s *DefaultJWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.config.SecretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
