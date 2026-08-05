package analytics

import (
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
	}

	api := r.Group("/api/v1")
	api.Use(middleware.AuthRequired(db, cfg))
	RegisterRoutes(api, db)

	return r
}

func TestAnalytics_Forbidden_FreeUser(t *testing.T) {
	db := setupTestDB(t)
	user, token := createUser(db, "free@example.com", false)
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

	req := httptest.NewRequest(http.MethodGet, "/api/v1/qrcodes/1/analytics/summary", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Analytics is only available for active Pro subscription")
}

func TestAnalytics_NotFound_UnownedQR(t *testing.T) {
	db := setupTestDB(t)
	user1, _ := createUser(db, "user1@example.com", true)
	_, token2 := createUser(db, "user2@example.com", true)
	r := setupRouter(db)

	qr := models.QRCode{
		UserID:    user1.ID,
		Title:     "User1 QR",
		QRType:    shared.QRTypeURL,
		Content:   "https://example.com",
		IsDynamic: true,
		Status:    shared.QRStatusActive,
	}
	db.Create(&qr)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/qrcodes/1/analytics/summary", nil)
	req.Header.Set("Authorization", "Bearer "+token2)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAnalytics_Success_ProUser(t *testing.T) {
	db := setupTestDB(t)
	user, token := createUser(db, "pro@example.com", true)
	r := setupRouter(db)

	shortCode := "abc1234"
	qr := models.QRCode{
		UserID:         user.ID,
		Title:          "Pro Dynamic QR",
		QRType:         shared.QRTypeURL,
		Content:        "http://localhost:8080/q/abc1234",
		ShortCode:      &shortCode,
		IsDynamic:      true,
		DestinationURL: "https://example.com",
		ScanCount:      2,
		Status:         shared.QRStatusActive,
	}
	db.Create(&qr)

	scan1 := models.QRScan{
		QRCodeID:        qr.ID,
		ScannedAt:       time.Now().Add(-2 * time.Hour),
		DeviceType:      "Mobile",
		Browser:         "Chrome",
		OperatingSystem: "Android",
	}
	scan2 := models.QRScan{
		QRCodeID:        qr.ID,
		ScannedAt:       time.Now().Add(-1 * time.Hour),
		DeviceType:      "Desktop",
		Browser:         "Edge",
		OperatingSystem: "Windows",
	}
	db.Create(&scan1)
	db.Create(&scan2)

	// Summary test
	req := httptest.NewRequest(http.MethodGet, "/api/v1/qrcodes/1/analytics/summary", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// By-Date test
	reqDate := httptest.NewRequest(http.MethodGet, "/api/v1/qrcodes/1/analytics/by-date", nil)
	reqDate.Header.Set("Authorization", "Bearer "+token)
	wDate := httptest.NewRecorder()

	r.ServeHTTP(wDate, reqDate)
	assert.Equal(t, http.StatusOK, wDate.Code)

	// By-Location test when country is empty
	reqLoc := httptest.NewRequest(http.MethodGet, "/api/v1/qrcodes/1/analytics/by-location", nil)
	reqLoc.Header.Set("Authorization", "Bearer "+token)
	wLoc := httptest.NewRecorder()

	r.ServeHTTP(wLoc, reqLoc)
	assert.Equal(t, http.StatusOK, wLoc.Code)
	
	var env struct {
		Data []StatRow `json:"data"`
	}
	json.Unmarshal(wLoc.Body.Bytes(), &env)
	assert.Equal(t, 0, len(env.Data), "Location array should be empty [] when no country data exists")
}
