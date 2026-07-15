package middleware

import (
	"strings"

	"qr-generator/backend/internal/config"
	"qr-generator/backend/internal/models"
	"qr-generator/backend/internal/shared"
	"qr-generator/backend/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const UserContextKey = "current_user"

func AuthRequired(db *gorm.DB, cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			shared.Error(c, 401, "Missing bearer token", nil)
			c.Abort()
			return
		}
		claims, err := utils.ParseToken(strings.TrimPrefix(header, "Bearer "), cfg.JWTAccessSecret)
		if err != nil {
			shared.Error(c, 401, "Invalid or expired token", nil)
			c.Abort()
			return
		}
		var user models.User
		if err := db.Preload("Roles").First(&user, claims.UserID).Error; err != nil || user.Status != shared.UserStatusActive {
			shared.Error(c, 401, "User is not active", nil)
			c.Abort()
			return
		}
		c.Set(UserContextKey, user)
		c.Next()
	}
}

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := CurrentUser(c)
		if !ok || !HasRole(user, shared.RoleNameAdmin) {
			shared.Error(c, 403, "Admin permission required", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

func CurrentUser(c *gin.Context) (models.User, bool) {
	value, ok := c.Get(UserContextKey)
	if !ok {
		return models.User{}, false
	}
	user, ok := value.(models.User)
	return user, ok
}

func HasRole(user models.User, role shared.RoleName) bool {
	for _, item := range user.Roles {
		if item.Name == role {
			return true
		}
	}
	return false
}
