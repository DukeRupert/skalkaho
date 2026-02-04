package auth

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHashPassword verifies that password hashing produces valid argon2id hashes.
func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "simple_password",
			password: "password123",
		},
		{
			name:     "complex_password",
			password: "C0mpl3x!P@ssw0rd#2024",
		},
		{
			name:     "long_password",
			password: strings.Repeat("a", 100),
		},
		{
			name:     "unicode_password",
			password: "パスワード123",
		},
		{
			name:     "empty_password",
			password: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			require.NoError(t, err, "HashPassword should not return an error")
			assert.NotEmpty(t, hash, "hash should not be empty")

			// Verify hash format (argon2id should produce a base64-encoded string)
			assert.True(t, len(hash) > 50, "hash should be sufficiently long")

			// Verify that hashing the same password twice produces different hashes (due to salt)
			hash2, err := HashPassword(tt.password)
			require.NoError(t, err)
			assert.NotEqual(t, hash, hash2, "hashing same password twice should produce different hashes")
		})
	}
}

// TestVerifyPassword verifies password verification logic.
func TestVerifyPassword(t *testing.T) {
	t.Run("correct_password_matches", func(t *testing.T) {
		password := "mySecurePassword123"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		match, err := VerifyPassword(password, hash)
		require.NoError(t, err)
		assert.True(t, match, "correct password should match")
	})

	t.Run("incorrect_password_fails", func(t *testing.T) {
		password := "mySecurePassword123"
		wrongPassword := "wrongPassword456"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		match, err := VerifyPassword(wrongPassword, hash)
		require.NoError(t, err)
		assert.False(t, match, "incorrect password should not match")
	})

	t.Run("empty_password_fails", func(t *testing.T) {
		password := "mySecurePassword123"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		match, err := VerifyPassword("", hash)
		require.NoError(t, err)
		assert.False(t, match, "empty password should not match non-empty hash")
	})

	t.Run("case_sensitive", func(t *testing.T) {
		password := "Password123"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		match, err := VerifyPassword("password123", hash)
		require.NoError(t, err)
		assert.False(t, match, "password verification should be case-sensitive")
	})

	t.Run("invalid_hash_format_returns_error", func(t *testing.T) {
		_, err := VerifyPassword("password123", "invalid-hash-format")
		assert.Error(t, err, "invalid hash format should return error")
	})

	t.Run("constant_time_comparison", func(t *testing.T) {
		// This test doesn't directly verify constant-time behavior (that would require timing analysis),
		// but it ensures the function completes successfully for both matching and non-matching cases
		password := "password123"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		// Multiple verifications should all complete without error
		for i := 0; i < 10; i++ {
			match, err := VerifyPassword(password, hash)
			require.NoError(t, err)
			assert.True(t, match)

			match, err = VerifyPassword("wrong"+password, hash)
			require.NoError(t, err)
			assert.False(t, match)
		}
	})
}

// TestGenerateSessionToken verifies session token generation.
func TestGenerateSessionToken(t *testing.T) {
	t.Run("generates_token_of_correct_length", func(t *testing.T) {
		token, err := GenerateSessionToken()
		require.NoError(t, err)
		assert.NotEmpty(t, token, "token should not be empty")

		// 32 bytes encoded as hex should be 64 characters, or as base64 should be ~44 characters
		// The actual length depends on encoding choice - verify it's reasonable
		assert.True(t, len(token) >= 32, "token should be at least 32 characters")
	})

	t.Run("generates_unique_tokens", func(t *testing.T) {
		tokens := make(map[string]bool)
		iterations := 100

		for i := 0; i < iterations; i++ {
			token, err := GenerateSessionToken()
			require.NoError(t, err)

			// Check for collisions
			assert.False(t, tokens[token], "token should be unique (collision detected)")
			tokens[token] = true
		}

		assert.Equal(t, iterations, len(tokens), "all tokens should be unique")
	})

	t.Run("tokens_have_sufficient_entropy", func(t *testing.T) {
		token, err := GenerateSessionToken()
		require.NoError(t, err)

		// Check that token doesn't contain only a small set of characters (poor entropy)
		uniqueChars := make(map[rune]bool)
		for _, ch := range token {
			uniqueChars[ch] = true
		}

		// A good token should have at least 10 unique characters
		assert.True(t, len(uniqueChars) >= 10, "token should have sufficient character diversity")
	})
}

// TestHashSessionToken verifies session token hashing.
func TestHashSessionToken(t *testing.T) {
	t.Run("produces_consistent_hash", func(t *testing.T) {
		token := "test-session-token-12345"

		hash1 := HashSessionToken(token)
		hash2 := HashSessionToken(token)

		assert.Equal(t, hash1, hash2, "hashing same token should produce same hash")
		assert.NotEmpty(t, hash1, "hash should not be empty")
	})

	t.Run("different_tokens_produce_different_hashes", func(t *testing.T) {
		token1 := "test-token-1"
		token2 := "test-token-2"

		hash1 := HashSessionToken(token1)
		hash2 := HashSessionToken(token2)

		assert.NotEqual(t, hash1, hash2, "different tokens should produce different hashes")
	})

	t.Run("hash_is_not_reversible", func(t *testing.T) {
		token := "secret-session-token"
		hash := HashSessionToken(token)

		// The hash should not contain the original token
		assert.NotContains(t, hash, token, "hash should not contain original token")
		assert.NotEqual(t, hash, token, "hash should not equal token")
	})

	t.Run("produces_fixed_length_hash", func(t *testing.T) {
		// SHA-256 produces 32-byte hashes. Encoded as hex = 64 chars, base64 = 44 chars
		token1 := "short"
		token2 := strings.Repeat("very-long-token", 100)

		hash1 := HashSessionToken(token1)
		hash2 := HashSessionToken(token2)

		// Both hashes should be the same length
		assert.Equal(t, len(hash1), len(hash2), "hashes should be fixed length regardless of input")
	})

	t.Run("empty_token_produces_hash", func(t *testing.T) {
		hash := HashSessionToken("")
		assert.NotEmpty(t, hash, "empty token should still produce a hash")
	})
}

// TestPasswordHashingRoundTrip verifies the complete hash and verify cycle.
func TestPasswordHashingRoundTrip(t *testing.T) {
	passwords := []string{
		"simple",
		"Complex!P@ssw0rd#123",
		"パスワード",
		strings.Repeat("long", 50),
		"with spaces and symbols !@#$%^&*()",
	}

	for _, password := range passwords {
		t.Run("password_"+password[:min(10, len(password))], func(t *testing.T) {
			// Hash the password
			hash, err := HashPassword(password)
			require.NoError(t, err)

			// Verify correct password
			match, err := VerifyPassword(password, hash)
			require.NoError(t, err)
			assert.True(t, match, "correct password should verify")

			// Verify incorrect password
			match, err = VerifyPassword(password+"wrong", hash)
			require.NoError(t, err)
			assert.False(t, match, "incorrect password should not verify")
		})
	}
}

// Helper function for Go < 1.21
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
