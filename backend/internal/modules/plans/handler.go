package plans

import (
	"qr-generator/backend/internal/models"
	"qr-generator/backend/internal/shared"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct{ db *gorm.DB }

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	h := &Handler{db: db}
	plans := rg.Group("/plans")
	plans.GET("", h.List)
	plans.GET("/:id", h.Detail)
}

func (h *Handler) List(c *gin.Context) {
	var plans []models.Plan
	h.db.Where("status = ?", shared.PlanStatusActive).Order("price asc").Find(&plans)
	shared.OK(c, "Success", plans)
}

func (h *Handler) Detail(c *gin.Context) {
	var plan models.Plan
	if err := h.db.First(&plan, c.Param("id")).Error; err != nil {
		shared.Error(c, 404, "Plan not found", nil)
		return
	}
	shared.OK(c, "Success", plan)
}
