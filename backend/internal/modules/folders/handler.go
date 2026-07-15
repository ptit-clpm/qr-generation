package folders

import (
	"qr-generator/backend/internal/middleware"
	"qr-generator/backend/internal/models"
	"qr-generator/backend/internal/shared"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FolderRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description"`
}

type Handler struct{ db *gorm.DB }

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	h := &Handler{db: db}
	folders := rg.Group("/folders")
	folders.POST("", h.Create)
	folders.GET("", h.List)
	folders.PUT("/:id", h.Update)
	folders.DELETE("/:id", h.Delete)
}

func (h *Handler) Create(c *gin.Context) {
	user, _ := middleware.CurrentUser(c)
	var req FolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, 400, "Validation error", err.Error())
		return
	}
	folder := models.Folder{UserID: user.ID, Name: req.Name, Description: req.Description}
	h.db.Create(&folder)
	shared.Created(c, "Folder created", folder)
}

func (h *Handler) List(c *gin.Context) {
	user, _ := middleware.CurrentUser(c)
	var folders []models.Folder
	h.db.Where("user_id = ?", user.ID).Order("created_at desc").Find(&folders)
	shared.OK(c, "Success", folders)
}

func (h *Handler) Update(c *gin.Context) {
	user, _ := middleware.CurrentUser(c)
	var req FolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, 400, "Validation error", err.Error())
		return
	}
	res := h.db.Model(&models.Folder{}).Where("id = ? AND user_id = ?", c.Param("id"), user.ID).Updates(req)
	if res.RowsAffected == 0 {
		shared.Error(c, 404, "Folder not found", nil)
		return
	}
	shared.OK(c, "Folder updated", nil)
}

func (h *Handler) Delete(c *gin.Context) {
	user, _ := middleware.CurrentUser(c)
	res := h.db.Where("id = ? AND user_id = ?", c.Param("id"), user.ID).Delete(&models.Folder{})
	if res.RowsAffected == 0 {
		shared.Error(c, 404, "Folder not found", nil)
		return
	}
	shared.OK(c, "Folder deleted", nil)
}
