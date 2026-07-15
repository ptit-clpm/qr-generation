package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"qr-generator/backend/internal/config"
	"qr-generator/backend/internal/models"
	"qr-generator/backend/internal/shared"
	"qr-generator/backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupHandlerTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.User{}, &models.Role{}, &models.Subscription{}, &models.Plan{}))
	db.Create(&models.Role{Name: shared.RoleNameUser})
	db.Create(&models.Role{Name: shared.RoleNameAdmin})
	db.Create(&models.Plan{
		Name:         shared.PlanNameFree,
		Price:        0,
		DurationDays: 3650,
		MaxQRCodes:   10,
		Status:       shared.PlanStatusActive,
	})
	return db
}

func seedUserForHandler(db *gorm.DB, fullName, email, password string, status shared.UserStatus, roleName shared.RoleName) models.User {
	hash, _ := utils.HashPassword(password)
	user := models.User{FullName: fullName, Email: email, PasswordHash: hash, Status: status}
	db.Create(&user)
	var role models.Role
	db.Where("name = ?", roleName).First(&role)
	db.Model(&user).Association("Roles").Append(&role)
	return user
}

type responseEnvelope struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
	Errors  json.RawMessage `json:"errors,omitempty"`
}

func makeHandler(t *testing.T, db *gorm.DB) *Handler {
	t.Helper()
	cfg := config.Config{
		JWTAccessSecret:    "test-access-secret",
		JWTRefreshSecret:   "test-refresh-secret",
		JWTAccessTTLMin:    60,
		JWTRefreshTTLHours: 720,
	}
	return NewHandler(db, cfg)
}

func setupRouter(handler *Handler) *gin.Engine {
	r := gin.New()
	r.POST("/api/v1/auth/register", handler.Register)
	r.POST("/api/v1/auth/login", handler.Login)
	r.POST("/api/v1/auth/refresh", handler.Refresh)
	return r
}

func doRequest(r *gin.Engine, method, path string, body any, token ...string) *httptest.ResponseRecorder {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}
	req := httptest.NewRequest(method, path, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if len(token) > 0 && token[0] != "" {
		req.Header.Set("Authorization", "Bearer "+token[0])
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func parseResponse(t *testing.T, w *httptest.ResponseRecorder) responseEnvelope {
	t.Helper()
	var resp responseEnvelope
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	return resp
}

// ========================================================================
// REGISTER TESTS
// ========================================================================

func TestRegister_Success(t *testing.T) {
	db := setupHandlerTestDB(t)
	handler := makeHandler(t, db)
	r := setupRouter(handler)

	w := doRequest(r, http.MethodPost, "/api/v1/auth/register", RegisterRequest{
		FullName: "Nguyen Van A", Email: "nguyenvana@example.com",
		Password: "password123", ConfirmPassword: "password123",
	})

	assert.Equal(t, 201, w.Code)
	resp := parseResponse(t, w)
	assert.Equal(t, "Register success", resp.Message)
	assert.True(t, resp.Success)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	db := setupHandlerTestDB(t)
	seedUserForHandler(db, "Existing", "existing@example.com", "password123", shared.UserStatusActive, shared.RoleNameUser)
	handler := makeHandler(t, db)
	r := setupRouter(handler)

	w := doRequest(r, http.MethodPost, "/api/v1/auth/register", RegisterRequest{
		FullName: "New User", Email: "existing@example.com",
		Password: "password123", ConfirmPassword: "password123",
	})

	assert.Equal(t, 409, w.Code)
	resp := parseResponse(t, w)
	assert.Equal(t, "Email already exists", resp.Message)
	assert.False(t, resp.Success)
}

func TestRegister_PasswordMismatch(t *testing.T) {
	db := setupHandlerTestDB(t)
	handler := makeHandler(t, db)
	r := setupRouter(handler)

	w := doRequest(r, http.MethodPost, "/api/v1/auth/register", RegisterRequest{
		FullName: "Mismatch User", Email: "mismatch@example.com",
		Password: "password123", ConfirmPassword: "differentPass",
	})

	assert.Equal(t, 400, w.Code)
	resp := parseResponse(t, w)
	assert.Equal(t, "Password confirmation does not match", resp.Message)
}

func TestRegister_ValidationError_EmptyFullName(t *testing.T) {
	db := setupHandlerTestDB(t)
	handler := makeHandler(t, db)
	r := setupRouter(handler)

	w := doRequest(r, http.MethodPost, "/api/v1/auth/register", RegisterRequest{
		FullName: "", Email: "valid@example.com",
		Password: "password123", ConfirmPassword: "password123",
	})

	assert.Equal(t, 400, w.Code)
}

func TestRegister_ValidationError_InvalidEmail(t *testing.T) {
	db := setupHandlerTestDB(t)
	handler := makeHandler(t, db)
	r := setupRouter(handler)

	w := doRequest(r, http.MethodPost, "/api/v1/auth/register", RegisterRequest{
		FullName: "User", Email: "not-an-email",
		Password: "password123", ConfirmPassword: "password123",
	})

	assert.Equal(t, 400, w.Code)
}

func TestRegister_ValidationError_ShortPassword(t *testing.T) {
	db := setupHandlerTestDB(t)
	handler := makeHandler(t, db)
	r := setupRouter(handler)

	w := doRequest(r, http.MethodPost, "/api/v1/auth/register", RegisterRequest{
		FullName: "User", Email: "valid@example.com",
		Password: "1234567", ConfirmPassword: "1234567",
	})

	assert.Equal(t, 400, w.Code)
}

func TestRegister_ValidationError_LongFullName(t *testing.T) {
	db := setupHandlerTestDB(t)
	handler := makeHandler(t, db)
	r := setupRouter(handler)

	longName := string(make([]byte, 151))
	w := doRequest(r, http.MethodPost, "/api/v1/auth/register", RegisterRequest{
		FullName: longName, Email: "valid@example.com",
		Password: "password123", ConfirmPassword: "password123",
	})

	assert.Equal(t, 400, w.Code)
}

func TestRegister_ValidationError_EmptyBody(t *testing.T) {
	db := setupHandlerTestDB(t)
	handler := makeHandler(t, db)
	r := setupRouter(handler)

	w := doRequest(r, http.MethodPost, "/api/v1/auth/register", nil)
	assert.Equal(t, 400, w.Code)
}

// ========================================================================
// LOGIN TESTS
// ========================================================================

func TestLogin_Success(t *testing.T) {
	db := setupHandlerTestDB(t)
	seedUserForHandler(db, "Test User", "test@example.com", "correctPassword", shared.UserStatusActive, shared.RoleNameUser)
	handler := makeHandler(t, db)
	r := setupRouter(handler)

	w := doRequest(r, http.MethodPost, "/api/v1/auth/login", LoginRequest{
		Email: "test@example.com", Password: "correctPassword",
	})

	assert.Equal(t, 200, w.Code)
	resp := parseResponse(t, w)
	assert.Equal(t, "Login success", resp.Message)
	assert.True(t, resp.Success)
}

func TestLogin_WrongPassword(t *testing.T) {
	db := setupHandlerTestDB(t)
	seedUserForHandler(db, "Test User", "test@example.com", "correctPassword", shared.UserStatusActive, shared.RoleNameUser)
	handler := makeHandler(t, db)
	r := setupRouter(handler)

	w := doRequest(r, http.MethodPost, "/api/v1/auth/login", LoginRequest{
		Email: "test@example.com", Password: "wrongPassword",
	})

	assert.Equal(t, 401, w.Code)
	assert.False(t, parseResponse(t, w).Success)
}

func TestLogin_UserNotFound(t *testing.T) {
	db := setupHandlerTestDB(t)
	handler := makeHandler(t, db)
	r := setupRouter(handler)

	w := doRequest(r, http.MethodPost, "/api/v1/auth/login", LoginRequest{
		Email: "nonexistent@example.com", Password: "password123",
	})

	assert.Equal(t, 401, w.Code)
}

func TestLogin_AccountLocked(t *testing.T) {
	db := setupHandlerTestDB(t)
	seedUserForHandler(db, "Locked User", "locked@example.com", "password123", shared.UserStatusLocked, shared.RoleNameUser)
	handler := makeHandler(t, db)
	r := setupRouter(handler)

	w := doRequest(r, http.MethodPost, "/api/v1/auth/login", LoginRequest{
		Email: "locked@example.com", Password: "password123",
	})

	assert.Equal(t, 403, w.Code)
	assert.False(t, parseResponse(t, w).Success)
}

func TestLogin_ValidationError_EmptyEmail(t *testing.T) {
	db := setupHandlerTestDB(t)
	handler := makeHandler(t, db)
	r := setupRouter(handler)

	w := doRequest(r, http.MethodPost, "/api/v1/auth/login", LoginRequest{
		Email: "", Password: "password123",
	})
	assert.Equal(t, 400, w.Code)
}

func TestLogin_ValidationError_InvalidEmailFormat(t *testing.T) {
	db := setupHandlerTestDB(t)
	handler := makeHandler(t, db)
	r := setupRouter(handler)

	w := doRequest(r, http.MethodPost, "/api/v1/auth/login", LoginRequest{
		Email: "invalid", Password: "password123",
	})
	assert.Equal(t, 400, w.Code)
}

func TestLogin_EmptyRequestBody(t *testing.T) {
	db := setupHandlerTestDB(t)
	handler := makeHandler(t, db)
	r := setupRouter(handler)

	w := doRequest(r, http.MethodPost, "/api/v1/auth/login", nil)
	assert.Equal(t, 400, w.Code)
}

// ========================================================================
// REFRESH TOKEN TESTS
// ========================================================================

func TestRefresh_Success(t *testing.T) {
	db := setupHandlerTestDB(t)
	user := seedUserForHandler(db, "Test User", "test@example.com", "password123", shared.UserStatusActive, shared.RoleNameUser)
	cfg := config.Config{JWTRefreshSecret: "test-refresh-secret", JWTAccessSecret: "test-access-secret", JWTAccessTTLMin: 60, JWTRefreshTTLHours: 720}
	handler := NewHandler(db, cfg)
	r := setupRouter(handler)

	refreshToken, _ := utils.GenerateToken(user.ID, user.Email, []string{"USER"}, cfg.JWTRefreshSecret, time.Duration(cfg.JWTRefreshTTLHours)*time.Hour)

	w := doRequest(r, http.MethodPost, "/api/v1/auth/refresh", RefreshRequest{RefreshToken: refreshToken})
	assert.Equal(t, 200, w.Code)
	resp := parseResponse(t, w)
	assert.Equal(t, "Refresh success", resp.Message)
}

func TestRefresh_InvalidToken(t *testing.T) {
	db := setupHandlerTestDB(t)
	handler := makeHandler(t, db)
	r := setupRouter(handler)

	w := doRequest(r, http.MethodPost, "/api/v1/auth/refresh", RefreshRequest{RefreshToken: "invalid-token"})
	assert.Equal(t, 401, w.Code)
}

func TestRefresh_ExpiredToken(t *testing.T) {
	db := setupHandlerTestDB(t)
	user := seedUserForHandler(db, "Test User", "test@example.com", "password123", shared.UserStatusActive, shared.RoleNameUser)
	cfg := config.Config{JWTRefreshSecret: "test-refresh-secret", JWTAccessSecret: "test-access-secret", JWTAccessTTLMin: 60, JWTRefreshTTLHours: 720}
	handler := NewHandler(db, cfg)
	r := setupRouter(handler)

	expiredRefresh, _ := utils.GenerateToken(user.ID, user.Email, []string{"USER"}, cfg.JWTRefreshSecret, -time.Hour)

	w := doRequest(r, http.MethodPost, "/api/v1/auth/refresh", RefreshRequest{RefreshToken: expiredRefresh})
	assert.Equal(t, 401, w.Code)
}

func TestRefresh_EmptyBody(t *testing.T) {
	db := setupHandlerTestDB(t)
	handler := makeHandler(t, db)
	r := setupRouter(handler)

	w := doRequest(r, http.MethodPost, "/api/v1/auth/refresh", nil)
	assert.Equal(t, 400, w.Code)
}

func TestRefresh_UserNotFound(t *testing.T) {
	db := setupHandlerTestDB(t)
	cfg := config.Config{JWTRefreshSecret: "test-refresh-secret", JWTAccessSecret: "test-access-secret", JWTAccessTTLMin: 60, JWTRefreshTTLHours: 720}
	handler := NewHandler(db, cfg)
	r := setupRouter(handler)

	token, _ := utils.GenerateToken(9999, "ghost@example.com", []string{"USER"}, cfg.JWTRefreshSecret, time.Hour)

	w := doRequest(r, http.MethodPost, "/api/v1/auth/refresh", RefreshRequest{RefreshToken: token})
	assert.Equal(t, 401, w.Code)
}

func TestRefresh_RefreshTokenUsesAccessSecret(t *testing.T) {
	db := setupHandlerTestDB(t)
	user := seedUserForHandler(db, "Test User", "test@example.com", "password123", shared.UserStatusActive, shared.RoleNameUser)
	cfg := config.Config{
		JWTAccessSecret:  "access-secret",
		JWTRefreshSecret: "refresh-secret",
	}
	handler := NewHandler(db, cfg)
	r := setupRouter(handler)

	accessToken, _ := utils.GenerateToken(user.ID, user.Email, []string{"USER"}, cfg.JWTAccessSecret, time.Hour)

	w := doRequest(r, http.MethodPost, "/api/v1/auth/refresh", RefreshRequest{RefreshToken: accessToken})
	assert.Equal(t, 401, w.Code, "access token should not work as refresh token")
}
