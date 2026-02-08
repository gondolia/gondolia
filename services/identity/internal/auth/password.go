package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 12
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("hashing password: %w", err)
	}
	return string(bytes), nil
}

// VerifyPassword compares a password with a hash
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateSecureToken generates a cryptographically secure random token
func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generating random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// HashToken creates a SHA256 hash of a token (for storage)
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// GenerateInvitationToken generates an invitation token
func GenerateInvitationToken() (string, error) {
	return GenerateSecureToken(32) // 64 hex characters
}

// GenerateRefreshToken generates a refresh token
func GenerateRefreshToken() (string, error) {
	return GenerateSecureToken(32) // 64 hex characters
}

// GeneratePasswordResetToken generates a password reset token
func GeneratePasswordResetToken() (string, error) {
	return GenerateSecureToken(32) // 64 hex characters
}

// ValidatePasswordStrength validates password meets requirements
func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	// Add more validation as needed (uppercase, lowercase, numbers, special chars)
	return nil
}
