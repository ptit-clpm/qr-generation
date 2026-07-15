package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

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

func setupAuthTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.User{}, &models.Role{}))
	db.Create(&models.Role{Name: shared.RoleNameUser})
	db.Create(&models.Role{Name: shared.RoleNameAdmin})
	return db
}

func seedUser(db *gorm.DB, fullName, email string, status shared.UserStatus, roleNames ...shared.RoleName) models.User {
	user := models.User{FullName: fullName, Email: email, PasswordHash: "hash", Status: status}
	db.Create(&user)
	for _, name := range roleNames {
		var role models.Role
		db.Where("name = ?", name).First(&role)
		db.Model(&user).Association("Roles").Append(&role)
	}
	return user
}

func TestAuthRequired_MissingBearerToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupAuthTestDB(t)
	cfg := config.Config{JWTAccessSecret: "test-secret"}

	r := gin.New()
	r.GET("/protected", AuthRequired(db, cfg), func(c *gin.Context) {
		c.JSON(200, gin.H{"data": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestAuthRequired_InvalidBearerFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupAuthTestDB(t)
	cfg := config.Config{JWTAccessSecret: "test-secret"}

	r := gin.New()
	r.GET("/protected", AuthRequired(db, cfg), func(c *gin.Context) {
		c.JSON(200, gin.H{"data": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "NotBearer token123")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestAuthRequired_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupAuthTestDB(t)
	cfg := config.Config{JWTAccessSecret: "test-secret"}

	r := gin.New()
	r.GET("/protected", AuthRequired(db, cfg), func(c *gin.Context) {
		c.JSON(200, gin.H{"data": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-string")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestAuthRequired_TokenFromWrongSecret(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupAuthTestDB(t)
	cfg := config.Config{JWTAccessSecret: "actual-secret"}

	wrongToken, _ := utils.GenerateToken(1, "user@test.com", []string{"USER"}, "wrong-secret", 3600)

	r := gin.New()
	r.GET("/protected", AuthRequired(db, cfg), func(c *gin.Context) {
		c.JSON(200, gin.H{"data": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+wrongToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestAuthRequired_ExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupAuthTestDB(t)
	cfg := config.Config{JWTAccessSecret: "test-secret"}

	expiredToken, _ := utils.GenerateToken(1, "user@test.com", []string{"USER"}, "test-secret", -1)

	r := gin.New()
	r.GET("/protected", AuthRequired(db, cfg), func(c *gin.Context) {
		c.JSON(200, gin.H{"data": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+expiredToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestAuthRequired_UserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupAuthTestDB(t)
	cfg := config.Config{JWTAccessSecret: "test-secret"}

	token, _ := utils.GenerateToken(9999, "nonexistent@test.com", []string{"USER"}, "test-secret", 3600)

	r := gin.New()
	r.GET("/protected", AuthRequired(db, cfg), func(c *gin.Context) {
		c.JSON(200, gin.H{"data": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestAuthRequired_UserIsLocked(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupAuthTestDB(t)
	cfg := config.Config{JWTAccessSecret: "test-secret"}

	user := seedUser(db, "Locked User", "locked@test.com", shared.UserStatusLocked, shared.RoleNameUser)
	token, _ := utils.GenerateToken(user.ID, user.Email, []string{"USER"}, "test-secret", 3600)

	r := gin.New()
	r.GET("/protected", AuthRequired(db, cfg), func(c *gin.Context) {
		c.JSON(200, gin.H{"data": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestAuthRequired_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupAuthTestDB(t)
	cfg := config.Config{JWTAccessSecret: "test-secret"}

	user := seedUser(db, "Active User", "active@test.com", shared.UserStatusActive, shared.RoleNameUser)
	token, _ := utils.GenerateToken(user.ID, user.Email, []string{"USER"}, "test-secret", 3600)

	r := gin.New()
	r.GET("/protected", AuthRequired(db, cfg), func(c *gin.Context) {
		currentUser, ok := CurrentUser(c)
		assert.True(t, ok)
		assert.Equal(t, user.ID, currentUser.ID)
		assert.Equal(t, user.Email, currentUser.Email)
		c.JSON(200, gin.H{"data": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestAdminRequired_NonAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupAuthTestDB(t)
	cfg := config.Config{JWTAccessSecret: "test-secret"}

	user := seedUser(db, "Regular User", "user@test.com", shared.UserStatusActive, shared.RoleNameUser)
	token, _ := utils.GenerateToken(user.ID, user.Email, []string{"USER"}, "test-secret", 3600)

	r := gin.New()
	r.GET("/admin", AuthRequired(db, cfg), AdminRequired(), func(c *gin.Context) {
		c.JSON(200, gin.H{"data": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code)
}

func TestAdminRequired_AdminSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupAuthTestDB(t)
	cfg := config.Config{JWTAccessSecret: "test-secret"}

	user := seedUser(db, "Admin User", "admin@test.com", shared.UserStatusActive, shared.RoleNameAdmin, shared.RoleNameUser)
	token, _ := utils.GenerateToken(user.ID, user.Email, []string{"ADMIN", "USER"}, "test-secret", 3600)

	r := gin.New()
	r.GET("/admin", AuthRequired(db, cfg), AdminRequired(), func(c *gin.Context) {
		currentUser, ok := CurrentUser(c)
		assert.True(t, ok)
		assert.True(t, HasRole(currentUser, shared.RoleNameAdmin))
		c.JSON(200, gin.H{"data": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestHasRole(t *testing.T) {
	db := setupAuthTestDB(t)
	user := seedUser(db, "Multi Role", "multi@test.com", shared.UserStatusActive, shared.RoleNameAdmin, shared.RoleNameUser)

	assert.True(t, HasRole(user, shared.RoleNameAdmin))
	assert.True(t, HasRole(user, shared.RoleNameUser))
	assert.False(t, HasRole(user, shared.RoleName("MANAGER")))
}

func TestCurrentUser_NoUserInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	_, ok := CurrentUser(c)
	assert.False(t, ok)
}

func TestCurrentUser_WrongTypeInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(UserContextKey, "not-a-user")

	_, ok := CurrentUser(c)
	assert.False(t, ok)
}
