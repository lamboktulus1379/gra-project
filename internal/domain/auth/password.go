package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// PasswordService provides methods for password hashing and verification
type PasswordService interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) (bool, error)
}

// ArgonParams defines the parameters used by the Argon2id algorithm
type ArgonParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// DefaultPasswordService is the default implementation of PasswordService
type DefaultPasswordService struct {
	params ArgonParams
}

// NewPasswordService creates a new password service with the given params
func NewPasswordService(params ArgonParams) PasswordService {
	return &DefaultPasswordService{
		params: params,
	}
}

// HashPassword hashes a password using Argon2id
func (s *DefaultPasswordService) HashPassword(password string) (string, error) {
	salt := make([]byte, s.params.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, s.params.Iterations, s.params.Memory, s.params.Parallelism, s.params.KeyLength)

	// Base64 encode the salt and hash
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Format: $argon2id$v=19$m=65536,t=3,p=2$<salt>$<hash>
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, s.params.Memory, s.params.Iterations, s.params.Parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

// VerifyPassword verifies the password against the hash
func (s *DefaultPasswordService) VerifyPassword(hashedPassword, password string) (bool, error) {
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid hash format")
	}

	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return false, errors.New("invalid hash format")
	}

	if version != argon2.Version {
		return false, errors.New("incompatible argon2 version")
	}

	var memory, iterations uint32
	var parallelism uint8
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return false, errors.New("invalid hash format")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, errors.New("invalid hash format")
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, errors.New("invalid hash format")
	}

	keyLength := uint32(len(decodedHash))

	// Compute the hash of the provided password using the same parameters
	comparisonHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLength)

	// Constant-time comparison to prevent timing attacks
	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1, nil
}
