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
	rg.GET("/qrcodes/:id/analytics/by-location", h.ByLocation)
}

type StatRow struct {
	Label string `json:"label"`
	Count int64  `json:"count"`
}

func (h *Handler) Summary(c *gin.Context) {
	qr, ok := h.authorize(c)
	if !ok {
		return
	}

	var firstScan, lastScan models.QRScan
	var firstTime, lastTime *time.Time

	if err := h.db.Where("qr_code_id = ?", qr.ID).Order("scanned_at asc").First(&firstScan).Error; err == nil {
		firstTime = &firstScan.ScannedAt
	}
	if err := h.db.Where("qr_code_id = ?", qr.ID).Order("scanned_at desc").First(&lastScan).Error; err == nil {
		lastTime = &lastScan.ScannedAt
	}

	topDevice := ""
	var deviceRows []StatRow
	if err := h.db.Model(&models.QRScan{}).
		Select("device_type as label, COUNT(*) as count").
		Where("qr_code_id = ? AND device_type IS NOT NULL AND device_type <> ''", qr.ID).
		Group("device_type").
		Order("count desc").
		Limit(1).
		Scan(&deviceRows).Error; err == nil && len(deviceRows) > 0 {
		topDevice = deviceRows[0].Label
	}

	topBrowser := ""
	var browserRows []StatRow
	if err := h.db.Model(&models.QRScan{}).
		Select("browser as label, COUNT(*) as count").
		Where("qr_code_id = ? AND browser IS NOT NULL AND browser <> ''", qr.ID).
		Group("browser").
		Order("count desc").
		Limit(1).
		Scan(&browserRows).Error; err == nil && len(browserRows) > 0 {
		topBrowser = browserRows[0].Label
	}

	shared.OK(c, "Success", gin.H{
		"qr_id":       qr.ID,
		"scan_count":  qr.ScanCount,
		"first_scan":  firstTime,
		"last_scan":   lastTime,
		"top_device":  topDevice,
		"top_browser": topBrowser,
	})
}

func (h *Handler) ByDate(c *gin.Context) {
	qr, ok := h.authorize(c)
	if !ok {
		return
	}

	fromDate := c.Query("from")
	toDate := c.Query("to")

	if fromDate == "" {
		fromDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	if toDate == "" {
		toDate = time.Now().Format("2006-01-02")
	}

	dateExpr := "DATE_FORMAT(scanned_at, '%Y-%m-%d')"
	if h.db.Dialector.Name() == "sqlite" {
		dateExpr = "strftime('%Y-%m-%d', scanned_at)"
	}

	rows := make([]StatRow, 0)
	err := h.db.Model(&models.QRScan{}).
		Select(dateExpr+" as label, COUNT(*) as count").
		Where("qr_code_id = ? AND "+dateExpr+" >= ? AND "+dateExpr+" <= ?", qr.ID, fromDate, toDate).
		Group(dateExpr).
		Order("label asc").
		Scan(&rows).Error

	if err != nil {
		shared.Error(c, 500, "Database query failed", err.Error())
		return
	}

	shared.OK(c, "Success", rows)
}

func (h *Handler) ByField(field string) gin.HandlerFunc {
	return func(c *gin.Context) {
		qr, ok := h.authorize(c)
		if !ok {
			return
		}
		rows := make([]StatRow, 0)
		err := h.db.Model(&models.QRScan{}).
			Select(field+" as label, COUNT(*) as count").
			Where("qr_code_id = ? AND "+field+" IS NOT NULL AND "+field+" <> ''", qr.ID).
			Group(field).
			Order("count desc").
			Scan(&rows).Error

		if err != nil {
			shared.Error(c, 500, "Database query failed", err.Error())
			return
		}

		shared.OK(c, "Success", rows)
	}
}

func (h *Handler) ByLocation(c *gin.Context) {
	qr, ok := h.authorize(c)
	if !ok {
		return
	}
	rows := make([]StatRow, 0)
	err := h.db.Model(&models.QRScan{}).
		Select("country as label, COUNT(*) as count").
		Where("qr_code_id = ? AND country IS NOT NULL AND country <> ''", qr.ID).
		Group("country").
		Order("count desc").
		Scan(&rows).Error

	if err != nil {
		shared.Error(c, 500, "Database query failed", err.Error())
		return
	}

	shared.OK(c, "Success", rows)
}

func (h *Handler) authorize(c *gin.Context) (models.QRCode, bool) {
	user, _ := middleware.CurrentUser(c)
	var qr models.QRCode
	if err := h.db.Where("id = ? AND user_id = ? AND status <> ?", c.Param("id"), user.ID, shared.QRStatusDeleted).First(&qr).Error; err != nil {
		shared.Error(c, 404, "QR code not found", nil)
		return qr, false
	}
	if !h.hasAnalytics(user.ID) {
		shared.Error(c, 403, "Analytics is only available for active Pro subscription", nil)
		return models.QRCode{}, false
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

