package auth

import (
	"time"

	"qr-generator/backend/internal/config"
	"qr-generator/backend/internal/middleware"
	"qr-generator/backend/internal/models"
	"qr-generator/backend/internal/shared"
	"qr-generator/backend/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db  *gorm.DB
	cfg config.Config
}

func NewHandler(db *gorm.DB, cfg config.Config) *Handler {
	return &Handler{db: db, cfg: cfg}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, 400, "Validation error", err.Error())
		return
	}
	if req.Password != req.ConfirmPassword {
		shared.Error(c, 400, "Password confirmation does not match", nil)
		return
	}

	var count int64
	h.db.Model(&models.User{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		shared.Error(c, 409, "Email already exists", nil)
		return
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		shared.Error(c, 500, "Could not hash password", nil)
		return
	}
	var userRole models.Role
	var freePlan models.Plan
	if err := h.db.Where("name = ?", shared.RoleNameUser).First(&userRole).Error; err != nil {
		shared.Error(c, 500, "Default role is missing", nil)
		return
	}
	if err := h.db.Where("name = ?", shared.PlanNameFree).First(&freePlan).Error; err != nil {
		shared.Error(c, 500, "Default plan is missing", nil)
		return
	}

	user := models.User{
		FullName: req.FullName, Email: req.Email, PhoneNumber: req.PhoneNumber,
		PasswordHash: hash, Status: shared.UserStatusActive, Roles: []models.Role{userRole},
	}
	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		return tx.Create(&models.Subscription{
			UserID: user.ID, PlanID: freePlan.ID, StartDate: time.Now(),
			EndDate: time.Now().AddDate(10, 0, 0), Status: shared.SubscriptionStatusActive,
		}).Error
	})
	if err != nil {
		shared.Error(c, 500, "Could not create account", nil)
		return
	}
	h.db.Preload("Roles").First(&user, user.ID)
	h.respondWithTokens(c, user, "Register success", 201)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, 400, "Validation error", err.Error())
		return
	}
	var user models.User
	if err := h.db.Preload("Roles").Where("email = ?", req.Email).First(&user).Error; err != nil {
		shared.Error(c, 401, "Invalid email or password", nil)
		return
	}
	if user.Status == shared.UserStatusLocked {
		shared.Error(c, 403, "Account is locked", nil)
		return
	}
	if user.Status != shared.UserStatusActive || !utils.CheckPassword(user.PasswordHash, req.Password) {
		shared.Error(c, 401, "Invalid email or password", nil)
		return
	}
	h.respondWithTokens(c, user, "Login success", 200)
}

func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, 400, "Validation error", err.Error())
		return
	}
	claims, err := utils.ParseToken(req.RefreshToken, h.cfg.JWTRefreshSecret)
	if err != nil {
		shared.Error(c, 401, "Invalid refresh token", nil)
		return
	}
	var user models.User
	if err := h.db.Preload("Roles").First(&user, claims.UserID).Error; err != nil || user.Status != shared.UserStatusActive {
		shared.Error(c, 401, "User is not active", nil)
		return
	}
	h.respondWithTokens(c, user, "Refresh success", 200)
}

func (h *Handler) Logout(c *gin.Context) {
	shared.OK(c, "Logout success", nil)
}

func (h *Handler) Me(c *gin.Context) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		shared.Error(c, 401, "Unauthorized", nil)
		return
	}
	shared.OK(c, "Success", user)
}

func (h *Handler) ChangePassword(c *gin.Context) {
	user, _ := middleware.CurrentUser(c)
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, 400, "Validation error", err.Error())
		return
	}
	var fresh models.User
	if err := h.db.First(&fresh, user.ID).Error; err != nil || !utils.CheckPassword(fresh.PasswordHash, req.OldPassword) {
		shared.Error(c, 400, "Old password is incorrect", nil)
		return
	}
	hash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		shared.Error(c, 500, "Could not hash password", nil)
		return
	}
	h.db.Model(&fresh).Update("password_hash", hash)
	shared.OK(c, "Password changed", nil)
}

func (h *Handler) respondWithTokens(c *gin.Context, user models.User, message string, status int) {
	roleNames := make([]string, 0, len(user.Roles))
	for _, role := range user.Roles {
		roleNames = append(roleNames, string(role.Name))
	}
	access, _ := utils.GenerateToken(user.ID, user.Email, roleNames, h.cfg.JWTAccessSecret, time.Duration(h.cfg.JWTAccessTTLMin)*time.Minute)
	refresh, _ := utils.GenerateToken(user.ID, user.Email, roleNames, h.cfg.JWTRefreshSecret, time.Duration(h.cfg.JWTRefreshTTLHours)*time.Hour)
	data := AuthResponse{AccessToken: access, RefreshToken: refresh, User: user}
	if status == 201 {
		shared.Created(c, message, data)
		return
	}
	shared.OK(c, message, data)
}
