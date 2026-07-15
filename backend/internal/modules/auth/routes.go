package auth

import (
	"qr-generator/backend/internal/config"
	"qr-generator/backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg config.Config) {
	h := NewHandler(db, cfg)
	auth := rg.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.POST("/refresh", h.Refresh)
	auth.POST("/logout", h.Logout)
	protected := auth.Group("")
	protected.Use(middleware.AuthRequired(db, cfg))
	protected.GET("/me", h.Me)
	protected.POST("/change-password", h.ChangePassword)
}
