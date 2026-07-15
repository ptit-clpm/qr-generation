package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword_Success(t *testing.T) {
	hash, err := HashPassword("password123")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, "password123", hash)
}

func TestCheckPassword_Correct(t *testing.T) {
	hash, err := HashPassword("my-secure-pass")
	require.NoError(t, err)
	assert.True(t, CheckPassword(hash, "my-secure-pass"))
}

func TestCheckPassword_Wrong(t *testing.T) {
	hash, err := HashPassword("correct-password")
	require.NoError(t, err)
	assert.False(t, CheckPassword(hash, "wrong-password"))
}

func TestHashPassword_EmptyString(t *testing.T) {
	hash, err := HashPassword("")
	require.NoError(t, err)
	assert.True(t, CheckPassword(hash, ""))
}

func TestHashPassword_Uniqueness(t *testing.T) {
	hash1, _ := HashPassword("samepassword")
	hash2, _ := HashPassword("samepassword")
	assert.NotEqual(t, hash1, hash2, "bcrypt should generate different salts")
}

func TestCheckPassword_InvalidHash(t *testing.T) {
	assert.False(t, CheckPassword("not-a-valid-hash", "password"))
}

func TestHashPassword_EmptyHash(t *testing.T) {
	assert.False(t, CheckPassword("", "password"))
}

func TestHashPassword_LongPassword(t *testing.T) {
	longPass := string(make([]byte, 72))
	hash, err := HashPassword(longPass)
	require.NoError(t, err)
	assert.True(t, CheckPassword(hash, longPass))
}
