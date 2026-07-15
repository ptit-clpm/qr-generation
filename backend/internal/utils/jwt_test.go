package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateToken_Success(t *testing.T) {
	secret := "test-access-secret"
	token, err := GenerateToken(1, "user@example.com", []string{"USER"}, secret, time.Hour)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := ParseToken(token, secret)
	require.NoError(t, err)
	assert.Equal(t, uint(1), claims.UserID)
	assert.Equal(t, "user@example.com", claims.Email)
	assert.Equal(t, []string{"USER"}, claims.Roles)
	assert.True(t, claims.ExpiresAt.Time.After(time.Now()))
}

func TestGenerateToken_WithMultipleRoles(t *testing.T) {
	secret := "multi-role-secret"
	token, err := GenerateToken(2, "admin@example.com", []string{"ADMIN", "USER"}, secret, time.Hour)
	require.NoError(t, err)

	claims, err := ParseToken(token, secret)
	require.NoError(t, err)
	assert.Equal(t, uint(2), claims.UserID)
	assert.Equal(t, []string{"ADMIN", "USER"}, claims.Roles)
}

func TestParseToken_ExpiredToken(t *testing.T) {
	secret := "test-secret"
	token, err := GenerateToken(1, "user@example.com", nil, secret, -time.Hour)
	require.NoError(t, err)

	_, err = ParseToken(token, secret)
	assert.Error(t, err)
	assert.ErrorIs(t, err, jwt.ErrTokenExpired)
}

func TestParseToken_WrongSecret(t *testing.T) {
	token, err := GenerateToken(1, "user@example.com", nil, "correct-secret", time.Hour)
	require.NoError(t, err)

	_, err = ParseToken(token, "wrong-secret")
	assert.Error(t, err)
}

func TestParseToken_InvalidTokenString(t *testing.T) {
	_, err := ParseToken("invalid-token-string", "secret")
	assert.Error(t, err)
}

func TestParseToken_EmptyToken(t *testing.T) {
	_, err := ParseToken("", "secret")
	assert.Error(t, err)
}

func TestParseToken_TamperedToken(t *testing.T) {
	secret := "test-secret"
	token, err := GenerateToken(1, "user@example.com", nil, secret, time.Hour)
	require.NoError(t, err)

	tampered := token[:len(token)-5] + "XXXXX"
	_, err = ParseToken(tampered, secret)
	assert.Error(t, err)
}

func TestGenerateToken_ContextualData(t *testing.T) {
	secret := "context-secret"
	now := time.Now()
	token, err := GenerateToken(3, "context@example.com", []string{"USER"}, secret, 30*time.Minute)
	require.NoError(t, err)

	claims, err := ParseToken(token, secret)
	require.NoError(t, err)
	assert.Equal(t, uint(3), claims.UserID)
	assert.Equal(t, "context@example.com", claims.Email)
	assert.True(t, claims.IssuedAt.Time.After(now.Add(-time.Second)) || claims.IssuedAt.Time.Equal(now))
	assert.True(t, claims.IssuedAt.Time.Before(now.Add(time.Second)))
}
