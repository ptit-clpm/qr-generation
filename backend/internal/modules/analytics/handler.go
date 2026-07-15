package analytics

import (
	"time"

	"qr-generator/backend/internal/middleware"
	"qr-generator/backend/internal/models"
	"qr-generator/backend/internal/shared"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct{ db *gorm.DB }

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	h := &Handler{db: db}
	rg.GET("/qrcodes/:id/analytics/summary", h.Summary)
	rg.GET("/qrcodes/:id/analytics/by-date", h.ByDate)
	rg.GET("/qrcodes/:id/analytics/by-device", h.ByField("device_type"))
	rg.GET("/qrcodes/:id/analytics/by-browser", h.ByField("browser"))
	rg.GET("/qrcodes/:id/analytics/by-location", h.ByField("country"))
}

func (h *Handler) Summary(c *gin.Context) {
	qr, ok := h.authorize(c)
	if !ok {
		return
	}
	var lastScan models.QRScan
	h.db.Where("qr_code_id = ?", qr.ID).Order("scanned_at desc").First(&lastScan)
	shared.OK(c, "Success", gin.H{"qr_id": qr.ID, "scan_count": qr.ScanCount, "last_scan": lastScan.ScannedAt})
}

func (h *Handler) ByDate(c *gin.Context) {
	qr, ok := h.authorize(c)
	if !ok {
		return
	}
	type row struct {
		Label string `json:"label"`
		Count int64  `json:"count"`
	}
	var rows []row
	h.db.Model(&models.QRScan{}).
		Select("DATE(scanned_at) as label, COUNT(*) as count").
		Where("qr_code_id = ?", qr.ID).
		Group("DATE(scanned_at)").
		Order("label asc").
		Scan(&rows)
	shared.OK(c, "Success", rows)
}

func (h *Handler) ByField(field string) gin.HandlerFunc {
	return func(c *gin.Context) {
		qr, ok := h.authorize(c)
		if !ok {
			return
		}
		type row struct {
			Label string `json:"label"`
			Count int64  `json:"count"`
		}
		var rows []row
		h.db.Model(&models.QRScan{}).
			Select(field+" as label, COUNT(*) as count").
			Where("qr_code_id = ?", qr.ID).
			Group(field).
			Order("count desc").
			Scan(&rows)
		shared.OK(c, "Success", rows)
	}
}

func (h *Handler) authorize(c *gin.Context) (models.QRCode, bool) {
	user, _ := middleware.CurrentUser(c)
	if !h.hasAnalytics(user.ID) {
		shared.Error(c, 403, "Analytics is only available for Pro users", nil)
		return models.QRCode{}, false
	}
	var qr models.QRCode
	if err := h.db.Where("id = ? AND user_id = ? AND status <> ?", c.Param("id"), user.ID, shared.QRStatusDeleted).First(&qr).Error; err != nil {
		shared.Error(c, 404, "QR code not found", nil)
		return qr, false
	}
	return qr, true
}

func (h *Handler) hasAnalytics(userID uint) bool {
	var count int64
	h.db.Model(&models.Subscription{}).
		Joins("JOIN plans ON plans.id = subscriptions.plan_id").
		Where("subscriptions.user_id = ? AND subscriptions.status = ? AND subscriptions.end_date > ? AND plans.allow_analytics = ?",
			userID, shared.SubscriptionStatusActive, time.Now(), true).
		Count(&count)
	return count > 0
}
