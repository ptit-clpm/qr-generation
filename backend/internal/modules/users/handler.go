package users

import (
	"qr-generator/backend/internal/middleware"
	"qr-generator/backend/internal/models"
	"qr-generator/backend/internal/shared"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UpdateProfileRequest struct {
	FullName    string `json:"full_name" binding:"required,min=2,max=150"`
	PhoneNumber string `json:"phone_number"`
	AvatarURL   string `json:"avatar_url"`
}

type Handler struct {
	db *gorm.DB
}

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	h := &Handler{db: db}
	users := rg.Group("/users")
	users.GET("/profile", h.Profile)
	users.PUT("/profile", h.UpdateProfile)
	users.GET("/subscription", h.Subscription)
	users.GET("/payments", h.Payments)
}

func (h *Handler) Profile(c *gin.Context) {
	user, _ := middleware.CurrentUser(c)
	shared.OK(c, "Success", user)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	user, _ := middleware.CurrentUser(c)
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, 400, "Validation error", err.Error())
		return
	}
	h.db.Model(&models.User{}).Where("id = ?", user.ID).Updates(map[string]any{
		"full_name": req.FullName, "phone_number": req.PhoneNumber, "avatar_url": req.AvatarURL,
	})
	var updated models.User
	h.db.Preload("Roles").First(&updated, user.ID)
	shared.OK(c, "Profile updated", updated)
}

func (h *Handler) Subscription(c *gin.Context) {
	user, _ := middleware.CurrentUser(c)
	var sub models.Subscription
	if err := h.db.Preload("Plan").Where("user_id = ? AND status = ?", user.ID, shared.SubscriptionStatusActive).Order("end_date desc").First(&sub).Error; err != nil {
		shared.Error(c, 404, "Subscription not found", nil)
		return
	}
	shared.OK(c, "Success", sub)
}

func (h *Handler) Payments(c *gin.Context) {
	user, _ := middleware.CurrentUser(c)
	var payments []models.Payment
	h.db.Where("user_id = ?", user.ID).Order("created_at desc").Find(&payments)
	shared.OK(c, "Success", payments)
}
