package qrcodes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"qr-generator/backend/internal/config"
	"qr-generator/backend/internal/middleware"
	"qr-generator/backend/internal/models"
	"qr-generator/backend/internal/shared"
	"qr-generator/backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(
		&models.User{}, &models.Role{}, &models.Plan{}, &models.Subscription{},
		&models.QRCode{}, &models.QRDesign{}, &models.QRScan{},
	))

	db.Create(&models.Role{Name: shared.RoleNameUser})
	db.Create(&models.Role{Name: shared.RoleNameAdmin})

	freePlan := models.Plan{
		Name:           shared.PlanNameFree,
		Price:          0,
		MaxQRCodes:     10,
		AllowDynamicQR: false,
		AllowAnalytics: false,
		Status:         shared.PlanStatusActive,
	}
	proPlan := models.Plan{
		Name:           shared.PlanNamePro,
		Price:          99000,
		MaxQRCodes:     1000,
		AllowDynamicQR: true,
		AllowAnalytics: true,
		Status:         shared.PlanStatusActive,
	}
	db.Create(&freePlan)
	db.Create(&proPlan)

	return db
}

func createUser(db *gorm.DB, email string, isPro bool) (models.User, string) {
	hash, _ := utils.HashPassword("Password123!")
	user := models.User{
		FullName:     "Test User",
		Email:        email,
		PasswordHash: hash,
		Status:       shared.UserStatusActive,
	}
	db.Create(&user)

	var userRole models.Role
	db.Where("name = ?", shared.RoleNameUser).First(&userRole)
	db.Model(&user).Association("Roles").Append(&userRole)

	var plan models.Plan
	planName := shared.PlanNameFree
	if isPro {
		planName = shared.PlanNamePro
	}
	db.Where("name = ?", planName).First(&plan)

	sub := models.Subscription{
		UserID:    user.ID,
		PlanID:    plan.ID,
		StartDate: time.Now().Add(-24 * time.Hour),
		EndDate:   time.Now().Add(30 * 24 * time.Hour),
		Status:    shared.SubscriptionStatusActive,
	}
	db.Create(&sub)

	token, _ := utils.GenerateToken(user.ID, user.Email, []string{string(shared.RoleNameUser)}, "test-secret", 24*time.Hour)
	return user, token
}

func setupRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	cfg := config.Config{
		JWTAccessSecret: "test-secret",
		AppURL:          "http://localhost:8080",
	}

	api := r.Group("/api/v1")
	api.Use(middleware.AuthRequired(db, cfg))
	RegisterRoutes(api, db, cfg)
	RegisterPublicRoutes(r, db, cfg)

	return r
}

func TestCreate_StaticURL_Success(t *testing.T) {
	db := setupTestDB(t)
	_, token := createUser(db, "free@example.com", false)
	r := setupRouter(db)

	reqBody := map[string]any{
		"title":      "Static URL QR",
		"qr_type":    "URL",
		"content":    "https://example.com",
		"is_dynamic": false,
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/qrcodes", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreate_DynamicURL_Success(t *testing.T) {
	db := setupTestDB(t)
	_, token := createUser(db, "pro@example.com", true)
	r := setupRouter(db)

	reqBody := map[string]any{
		"title":           "Dynamic URL QR",
		"qr_type":         "URL",
		"content":         "https://example.com",
		"is_dynamic":      true,
		"destination_url": "https://example.com/target",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/qrcodes", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var qr models.QRCode
	db.Where("title = ?", "Dynamic URL QR").First(&qr)
	assert.True(t, qr.IsDynamic)
	assert.NotNil(t, qr.ShortCode)
	assert.Equal(t, "https://example.com/target", qr.DestinationURL)
}

func TestCreate_DynamicNonURL_Rejected(t *testing.T) {
	db := setupTestDB(t)
	_, token := createUser(db, "pro@example.com", true)
	r := setupRouter(db)

	nonUrlTypes := []string{"SOCIAL", "PDF", "MENU", "WIFI", "VCARD", "EMAIL", "SMS", "LOCATION"}

	for _, qrType := range nonUrlTypes {
		reqBody := map[string]any{
			"title":           "Invalid Dynamic QR",
			"qr_type":         qrType,
			"content":         "https://example.com",
			"is_dynamic":      true,
			"destination_url": "https://example.com/target",
		}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/qrcodes", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code, "QR type %s should be rejected when is_dynamic=true", qrType)
		assert.Contains(t, w.Body.String(), "Dynamic QR chỉ hỗ trợ QR type URL")
	}
}

func TestCreate_DynamicURL_MissingDestination(t *testing.T) {
	db := setupTestDB(t)
	_, token := createUser(db, "pro@example.com", true)
	r := setupRouter(db)

	reqBody := map[string]any{
		"title":      "Missing Dest QR",
		"qr_type":    "URL",
		"content":    "https://example.com",
		"is_dynamic": true,
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/qrcodes", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdate_StaticQR_ContentImmutable(t *testing.T) {
	db := setupTestDB(t)
	user, token := createUser(db, "pro@example.com", true)
	r := setupRouter(db)

	qr := models.QRCode{
		UserID:    user.ID,
		Title:     "Static QR",
		QRType:    shared.QRTypeURL,
		Content:   "https://example.com",
		IsDynamic: false,
		Status:    shared.QRStatusActive,
	}
	db.Create(&qr)

	reqBody := map[string]any{
		"title":   "Updated Static QR Title",
		"content": "https://changed.com",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/qrcodes/1", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Static QR content cannot be changed")
}

func TestUpdate_DynamicURL_DestinationUrlSuccess(t *testing.T) {
	db := setupTestDB(t)
	user, token := createUser(db, "pro@example.com", true)
	r := setupRouter(db)

	shortCode := "abcd123"
	qr := models.QRCode{
		UserID:         user.ID,
		Title:          "Dynamic QR",
		QRType:         shared.QRTypeURL,
		Content:        "http://localhost:8080/q/abcd123",
		ShortCode:      &shortCode,
		IsDynamic:      true,
		DestinationURL: "https://original.com",
		Status:         shared.QRStatusActive,
	}
	db.Create(&qr)

	reqBody := map[string]any{
		"title":           "Dynamic QR Updated",
		"destination_url": "https://updated-target.com",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/qrcodes/1", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var updatedQr models.QRCode
	db.First(&updatedQr, 1)
	assert.Equal(t, "https://updated-target.com", updatedQr.DestinationURL)
}

func TestRedirect_DynamicURL_Success(t *testing.T) {
	db := setupTestDB(t)
	user, _ := createUser(db, "pro@example.com", true)
	r := setupRouter(db)

	shortCode := "testcode123"
	qr := models.QRCode{
		UserID:         user.ID,
		Title:          "Dynamic QR Redirect Test",
		QRType:         shared.QRTypeURL,
		Content:        "http://localhost:8080/q/testcode123",
		ShortCode:      &shortCode,
		IsDynamic:      true,
		DestinationURL: "https://destination.com",
		ScanCount:      0,
		Status:         shared.QRStatusActive,
	}
	db.Create(&qr)

	req := httptest.NewRequest(http.MethodGet, "/q/testcode123", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X)")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "https://destination.com", w.Header().Get("Location"))

	var updatedQr models.QRCode
	db.First(&updatedQr, qr.ID)
	assert.Equal(t, int64(1), updatedQr.ScanCount)

	var scanCount int64
	db.Model(&models.QRScan{}).Where("qr_code_id = ?", qr.ID).Count(&scanCount)
	assert.Equal(t, int64(1), scanCount)
}

func TestRedirect_StaticQR_NotFound(t *testing.T) {
	db := setupTestDB(t)
	user, _ := createUser(db, "user@example.com", false)
	r := setupRouter(db)

	qr := models.QRCode{
		UserID:    user.ID,
		Title:     "Static QR",
		QRType:    shared.QRTypeURL,
		Content:   "https://static-url.com",
		IsDynamic: false,
		Status:    shared.QRStatusActive,
	}
	db.Create(&qr)

	req := httptest.NewRequest(http.MethodGet, "/q/nonexistent", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
