package router

import (
	"time"

	"qr-generator/backend/internal/config"
	"qr-generator/backend/internal/middleware"
	"qr-generator/backend/internal/modules/admin"
	"qr-generator/backend/internal/modules/analytics"
	"qr-generator/backend/internal/modules/auth"
	"qr-generator/backend/internal/modules/folders"
	"qr-generator/backend/internal/modules/payments"
	"qr-generator/backend/internal/modules/plans"
	"qr-generator/backend/internal/modules/qrcodes"
	"qr-generator/backend/internal/modules/users"
	"qr-generator/backend/internal/shared"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, cfg config.Config, logger *zap.Logger) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Info("request", zap.String("method", c.Request.Method), zap.String("path", c.Request.URL.Path), zap.Int("status", c.Writer.Status()), zap.Duration("latency", time.Since(start)))
	})
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontendURL, cfg.AdminFrontendURL, "http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		shared.OK(c, "OK", gin.H{"status": "healthy"})
	})
	qrcodes.RegisterPublicRoutes(r, db, cfg)

	api := r.Group("/api/v1")
	auth.RegisterRoutes(api, db, cfg)
	plans.RegisterRoutes(api, db)
	payments.RegisterPublicRoutes(api, db, cfg)

	protected := api.Group("")
	protected.Use(middleware.AuthRequired(db, cfg))
	users.RegisterRoutes(protected, db)
	qrcodes.RegisterRoutes(protected, db, cfg)
	folders.RegisterRoutes(protected, db)
	payments.RegisterRoutes(protected, db, cfg)
	analytics.RegisterRoutes(protected, db)

	adminGroup := protected.Group("/admin")
	adminGroup.Use(middleware.AdminRequired())
	admin.RegisterRoutes(adminGroup, db)

	return r
}
