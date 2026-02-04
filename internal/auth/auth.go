package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const (
	// Argon2id parameters (OWASP recommendations)
	argonTime      = 1
	argonMemory    = 64 * 1024 // 64 MB
	argonThreads   = 4
	argonKeyLength = 32
	saltLength     = 16
)

// HashPassword hashes a password using argon2id.
func HashPassword(password string) (string, error) {
	// Generate random salt
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generating salt: %w", err)
	}

	// Hash password with argon2id
	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLength)

	// Encode salt + hash as base64
	// Format: salt (16 bytes) + hash (32 bytes) = 48 bytes
	encoded := make([]byte, saltLength+argonKeyLength)
	copy(encoded[:saltLength], salt)
	copy(encoded[saltLength:], hash)

	return base64.RawStdEncoding.EncodeToString(encoded), nil
}

// VerifyPassword verifies a password against an argon2id hash using constant-time comparison.
func VerifyPassword(password, encodedHash string) (bool, error) {
	// Decode the hash
	decoded, err := base64.RawStdEncoding.DecodeString(encodedHash)
	if err != nil {
		return false, fmt.Errorf("decoding hash: %w", err)
	}

	// Verify length
	if len(decoded) != saltLength+argonKeyLength {
		return false, fmt.Errorf("invalid hash format: expected %d bytes, got %d", saltLength+argonKeyLength, len(decoded))
	}

	// Extract salt and hash
	salt := decoded[:saltLength]
	expectedHash := decoded[saltLength:]

	// Hash the provided password with the same salt
	actualHash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLength)

	// Constant-time comparison
	match := subtle.ConstantTimeCompare(expectedHash, actualHash) == 1
	return match, nil
}

// GenerateSessionToken generates a random 32-byte session token.
func GenerateSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generating session token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// HashSessionToken hashes a session token using SHA-256 for storage.
func HashSessionToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.RawStdEncoding.EncodeToString(hash[:])
}
