package admin

import (
	"time"

	"qr-generator/backend/internal/models"
	"qr-generator/backend/internal/shared"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct{ db *gorm.DB }

type StatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type PlanRequest struct {
	Name              shared.PlanName   `json:"name" binding:"required"`
	Price             float64           `json:"price"`
	DurationDays      int               `json:"duration_days"`
	MaxQRCodes        int               `json:"max_qr_codes"`
	AllowDynamicQR    bool              `json:"allow_dynamic_qr"`
	AllowLogo         bool              `json:"allow_logo"`
	AllowAnalytics    bool              `json:"allow_analytics"`
	AllowSVGPDFExport bool              `json:"allow_svg_pdf_export"`
	Description       string            `json:"description"`
	Status            shared.PlanStatus `json:"status"`
}

type TemplateRequest struct {
	Name            string                `json:"name" binding:"required"`
	PreviewImageURL string                `json:"preview_image_url"`
	ConfigJSON      string                `json:"config_json"`
	IsPro           bool                  `json:"is_pro"`
	Status          shared.TemplateStatus `json:"status"`
}

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	h := &Handler{db: db}
	rg.GET("/dashboard", h.Dashboard)
	rg.GET("/users", h.Users)
	rg.GET("/users/:id", h.UserDetail)
	rg.PUT("/users/:id/status", h.UpdateUserStatus)
	rg.GET("/qrcodes", h.QRCodes)
	rg.GET("/qrcodes/:id", h.QRDetail)
	rg.PUT("/qrcodes/:id/status", h.UpdateQRStatus)
	rg.GET("/plans", h.Plans)
	rg.POST("/plans", h.CreatePlan)
	rg.PUT("/plans/:id", h.UpdatePlan)
	rg.GET("/payments", h.Payments)
	rg.GET("/templates", h.Templates)
	rg.POST("/templates", h.CreateTemplate)
	rg.PUT("/templates/:id", h.UpdateTemplate)
	rg.DELETE("/templates/:id", h.DeleteTemplate)
	rg.GET("/logs", h.Logs)
}

func (h *Handler) Dashboard(c *gin.Context) {
	var users, qrcodes, scans, payments int64
	var revenue float64
	h.db.Model(&models.User{}).Count(&users)
	h.db.Model(&models.QRCode{}).Where("status <> ?", shared.QRStatusDeleted).Count(&qrcodes)
	h.db.Model(&models.QRScan{}).Count(&scans)
	h.db.Model(&models.Payment{}).Where("status = ?", shared.PaymentStatusSuccess).Count(&payments)
	h.db.Model(&models.Payment{}).Where("status = ?", shared.PaymentStatusSuccess).Select("COALESCE(SUM(amount),0)").Scan(&revenue)
	shared.OK(c, "Success", gin.H{"users": users, "qrcodes": qrcodes, "scans": scans, "successful_payments": payments, "revenue": revenue})
}

func (h *Handler) Users(c *gin.Context) {
	var users []models.User
	h.db.Preload("Roles").Order("created_at desc").Find(&users)
	shared.OK(c, "Success", users)
}

func (h *Handler) UserDetail(c *gin.Context) {
	var user models.User
	if err := h.db.Preload("Roles").First(&user, c.Param("id")).Error; err != nil {
		shared.Error(c, 404, "User not found", nil)
		return
	}
	shared.OK(c, "Success", user)
}

func (h *Handler) UpdateUserStatus(c *gin.Context) {
	var req StatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, 400, "Validation error", err.Error())
		return
	}
	res := h.db.Model(&models.User{}).Where("id = ?", c.Param("id")).Update("status", req.Status)
	if res.RowsAffected == 0 {
		shared.Error(c, 404, "User not found", nil)
		return
	}
	shared.OK(c, "User status updated", nil)
}

func (h *Handler) QRCodes(c *gin.Context) {
	var qrs []models.QRCode
	h.db.Preload("Design").Order("created_at desc").Find(&qrs)
	shared.OK(c, "Success", qrs)
}

func (h *Handler) QRDetail(c *gin.Context) {
	var qr models.QRCode
	if err := h.db.Preload("Design").First(&qr, c.Param("id")).Error; err != nil {
		shared.Error(c, 404, "QR code not found", nil)
		return
	}
	shared.OK(c, "Success", qr)
}

func (h *Handler) UpdateQRStatus(c *gin.Context) {
	var req StatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, 400, "Validation error", err.Error())
		return
	}
	res := h.db.Model(&models.QRCode{}).Where("id = ?", c.Param("id")).Update("status", req.Status)
	if res.RowsAffected == 0 {
		shared.Error(c, 404, "QR code not found", nil)
		return
	}
	shared.OK(c, "QR status updated", nil)
}

func (h *Handler) Plans(c *gin.Context) {
	var plans []models.Plan
	h.db.Order("price asc").Find(&plans)
	shared.OK(c, "Success", plans)
}

func (h *Handler) CreatePlan(c *gin.Context) {
	var req PlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, 400, "Validation error", err.Error())
		return
	}
	plan := models.Plan{
		Name: req.Name, Price: req.Price, DurationDays: req.DurationDays, MaxQRCodes: req.MaxQRCodes,
		AllowDynamicQR: req.AllowDynamicQR, AllowLogo: req.AllowLogo, AllowAnalytics: req.AllowAnalytics,
		AllowSVGPDFExport: req.AllowSVGPDFExport, Description: req.Description, Status: req.Status,
	}
	if plan.Status == "" {
		plan.Status = shared.PlanStatusActive
	}
	if err := h.db.Create(&plan).Error; err != nil {
		shared.Error(c, 500, "Could not create plan", nil)
		return
	}
	shared.Created(c, "Plan created", plan)
}

func (h *Handler) UpdatePlan(c *gin.Context) {
	var req PlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, 400, "Validation error", err.Error())
		return
	}
	updates := map[string]any{
		"name": req.Name, "price": req.Price, "duration_days": req.DurationDays, "max_qr_codes": req.MaxQRCodes,
		"allow_dynamic_qr": req.AllowDynamicQR, "allow_logo": req.AllowLogo, "allow_analytics": req.AllowAnalytics,
		"allow_svg_pdf_export": req.AllowSVGPDFExport, "description": req.Description, "status": req.Status,
	}
	res := h.db.Model(&models.Plan{}).Where("id = ?", c.Param("id")).Updates(updates)
	if res.RowsAffected == 0 {
		shared.Error(c, 404, "Plan not found", nil)
		return
	}
	shared.OK(c, "Plan updated", nil)
}

func (h *Handler) Payments(c *gin.Context) {
	var payments []models.Payment
	h.db.Order("created_at desc").Find(&payments)
	shared.OK(c, "Success", payments)
}

func (h *Handler) Templates(c *gin.Context) {
	var templates []models.QRTemplate
	h.db.Order("created_at desc").Find(&templates)
	shared.OK(c, "Success", templates)
}

func (h *Handler) CreateTemplate(c *gin.Context) {
	var req TemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, 400, "Validation error", err.Error())
		return
	}
	template := models.QRTemplate{
		Name: req.Name, PreviewImageURL: req.PreviewImageURL, ConfigJSON: req.ConfigJSON,
		IsPro: req.IsPro, Status: req.Status,
	}
	if template.Status == "" {
		template.Status = shared.TemplateStatusActive
	}
	h.db.Create(&template)
	shared.Created(c, "Template created", template)
}

func (h *Handler) UpdateTemplate(c *gin.Context) {
	var req TemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, 400, "Validation error", err.Error())
		return
	}
	updates := map[string]any{
		"name": req.Name, "preview_image_url": req.PreviewImageURL, "config_json": req.ConfigJSON,
		"is_pro": req.IsPro, "status": req.Status,
	}
	res := h.db.Model(&models.QRTemplate{}).Where("id = ?", c.Param("id")).Updates(updates)
	if res.RowsAffected == 0 {
		shared.Error(c, 404, "Template not found", nil)
		return
	}
	shared.OK(c, "Template updated", nil)
}

func (h *Handler) DeleteTemplate(c *gin.Context) {
	res := h.db.Model(&models.QRTemplate{}).Where("id = ?", c.Param("id")).Updates(map[string]any{"status": shared.TemplateStatusDeleted, "updated_at": time.Now()})
	if res.RowsAffected == 0 {
		shared.Error(c, 404, "Template not found", nil)
		return
	}
	shared.OK(c, "Template deleted", nil)
}

func (h *Handler) Logs(c *gin.Context) {
	var logs []models.SystemLog
	h.db.Order("created_at desc").Limit(200).Find(&logs)
	shared.OK(c, "Success", logs)
}
