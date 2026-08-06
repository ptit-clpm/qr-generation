package qrcodes

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"qr-generator/backend/internal/config"
	"qr-generator/backend/internal/middleware"
	"qr-generator/backend/internal/models"
	"qr-generator/backend/internal/services"
	"qr-generator/backend/internal/shared"
	"qr-generator/backend/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db         *gorm.DB
	cfg        config.Config
	cloudinary *services.CloudinaryService
}

func NewHandler(db *gorm.DB, cfg config.Config) *Handler {
	return &Handler{db: db, cfg: cfg, cloudinary: services.NewCloudinaryService(cfg)}
}

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg config.Config) {
	h := NewHandler(db, cfg)
	qr := rg.Group("/qrcodes")
	qr.POST("", h.Create)
	qr.GET("", h.List)
	qr.GET("/:id", h.Detail)
	qr.PUT("/:id", h.Update)
	qr.DELETE("/:id", h.Delete)
	qr.POST("/:id/duplicate", h.Duplicate)
	qr.GET("/:id/download", h.Download)
	qr.GET("/:id/design", h.GetDesign)
	qr.PUT("/:id/design", h.UpdateDesign)
	qr.POST("/:id/logo", h.UploadLogo)
	qr.DELETE("/:id/logo", h.DeleteLogo)
}

func RegisterPublicRoutes(r *gin.Engine, db *gorm.DB, cfg config.Config) {
	h := NewHandler(db, cfg)
	r.GET("/q/:shortCode", h.Redirect)
}

func (h *Handler) Create(c *gin.Context) {
	user, _ := middleware.CurrentUser(c)
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, 400, "Validation error", err.Error())
		return
	}
	if req.IsDynamic && req.QRType != shared.QRTypeURL {
		shared.Error(c, 400, "Dynamic QR chỉ hỗ trợ QR type URL", nil)
		return
	}
	if err := h.validateQRInput(req.QRType, req.Content); err != nil {
		shared.Error(c, 400, err.Error(), nil)
		return
	}
	plan, isPro := h.currentPlan(user.ID)
	if req.IsDynamic && !plan.AllowDynamicQR {
		shared.Error(c, 403, "Dynamic QR is only available for Pro users", nil)
		return
	}
	if (req.QRType == shared.QRTypePDF || req.QRType == shared.QRTypeMenu || req.QRType == shared.QRTypeSocial) && !isPro {
		shared.Error(c, 403, "This QR type is only available for Pro users", nil)
		return
	}
	if req.Design != nil && req.Design.LogoURL != "" {
		if !plan.AllowLogo {
			shared.Error(c, 403, "Logo is only available for Pro users", nil)
			return
		}
		shared.Error(c, 400, "Logo must be uploaded through the logo upload endpoint", nil)
		return
	}
	if !isPro && h.countUserQRCodes(user.ID) >= int64(plan.MaxQRCodes) {
		shared.Error(c, 403, "Free QR limit reached. Please upgrade to Pro", nil)
		return
	}

	qr := models.QRCode{
		UserID: user.ID, FolderID: req.FolderID, Title: req.Title, QRType: req.QRType,
		Content: req.Content, IsDynamic: req.IsDynamic, DestinationURL: "",
		Status: shared.QRStatusActive,
	}
	if req.IsDynamic {
		if req.DestinationURL == "" {
			shared.Error(c, 400, "destination_url is required for Dynamic QR", nil)
			return
		}
		if !isValidURL(req.DestinationURL) {
			shared.Error(c, 400, "destination_url must be a valid URL", nil)
			return
		}
		qr.DestinationURL = req.DestinationURL
		shortCode := utils.GenerateShortCode()
		qr.ShortCode = &shortCode
		qr.Content = utils.DynamicURL(h.cfg.AppURL, shortCode)
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&qr).Error; err != nil {
			return err
		}
		design := designFromRequest(req.Design)
		design.QRCodeID = qr.ID
		return tx.Create(&design).Error
	})
	if err != nil {
		shared.Error(c, 500, "Could not create QR code", nil)
		return
	}
	h.db.Preload("Design").First(&qr, qr.ID)
	shared.Created(c, "QR code created", qr)
}

func (h *Handler) List(c *gin.Context) {
	user, _ := middleware.CurrentUser(c)
	page := shared.Pagination(c)
	query := h.db.Preload("Design").Where("user_id = ? AND status <> ?", user.ID, shared.QRStatusDeleted)
	if q := c.Query("q"); q != "" {
		query = query.Where("title LIKE ?", "%"+q+"%")
	}
	if qrType := c.Query("qr_type"); qrType != "" {
		query = query.Where("qr_type = ?", qrType)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if folderID := c.Query("folder_id"); folderID != "" {
		if folderID == "null" || folderID == "0" || folderID == "uncategorized" {
			query = query.Where("folder_id IS NULL")
		} else {
			query = query.Where("folder_id = ?", folderID)
		}
	}
	var total int64
	query.Model(&models.QRCode{}).Count(&total)
	var items []models.QRCode
	query.Order("created_at desc").Offset(page.Offset()).Limit(page.Limit).Find(&items)
	shared.OK(c, "Success", gin.H{"items": items, "total": total, "page": page.Page, "limit": page.Limit})
}

func (h *Handler) Detail(c *gin.Context) {
	user, _ := middleware.CurrentUser(c)
	qr, ok := h.findOwnedQR(c, user.ID)
	if !ok {
		return
	}
	shared.OK(c, "Success", qr)
}

func (h *Handler) Update(c *gin.Context) {
	user, _ := middleware.CurrentUser(c)
	qr, ok := h.findOwnedQR(c, user.ID)
	if !ok {
		return
	}
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, 400, "Validation error", err.Error())
		return
	}
	updates := map[string]any{"title": req.Title}
	if req.FolderID != nil {
		updates["folder_id"] = req.FolderID
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if qr.IsDynamic && qr.QRType == shared.QRTypeURL {
		if req.DestinationURL == "" || !isValidURL(req.DestinationURL) {
			shared.Error(c, 400, "destination_url must be a valid non-empty URL for Dynamic QR", nil)
			return
		}
		updates["destination_url"] = req.DestinationURL
	} else {
		if req.DestinationURL != "" && req.DestinationURL != qr.DestinationURL {
			shared.Error(c, 400, "Dynamic QR chỉ hỗ trợ QR type URL", nil)
			return
		}
		if req.Content != "" && req.Content != qr.Content {
			shared.Error(c, 400, "Static QR content cannot be changed after creation", nil)
			return
		}
	}
	h.db.Model(&qr).Updates(updates)
	h.db.Preload("Design").First(&qr, qr.ID)
	shared.OK(c, "QR code updated", qr)
}

func (h *Handler) Delete(c *gin.Context) {
	user, _ := middleware.CurrentUser(c)
	qr, ok := h.findOwnedQR(c, user.ID)
	if !ok {
		return
	}
	h.db.Model(&qr).Update("status", shared.QRStatusDeleted)
	shared.OK(c, "QR code deleted", nil)
}

func (h *Handler) Duplicate(c *gin.Context) {
	user, _ := middleware.CurrentUser(c)
	_, isPro := h.currentPlan(user.ID)
	if !isPro {
		shared.Error(c, 403, "Duplicate QR is only available for Pro users", nil)
		return
	}
	qr, ok := h.findOwnedQR(c, user.ID)
	if !ok {
		return
	}
	copy := qr
	copy.ID = 0
	copy.Design = models.QRDesign{}
	copy.Title = qr.Title + " Copy"
	copy.ShortCode = nil
	copy.IsDynamic = false
	copy.ScanCount = 0
	copy.CreatedAt = time.Time{}
	copy.UpdatedAt = time.Time{}
	if err := h.db.Create(&copy).Error; err != nil {
		shared.Error(c, 500, "Could not duplicate QR", nil)
		return
	}
	shared.Created(c, "QR code duplicated", copy)
}

func (h *Handler) Download(c *gin.Context) {
	user, _ := middleware.CurrentUser(c)
	qr, ok := h.findOwnedQR(c, user.ID)
	if !ok {
		return
	}
	content := qr.Content
	if qr.IsDynamic && qr.ShortCode != nil {
		content = utils.DynamicURL(h.cfg.AppURL, *qr.ShortCode)
	}
	size := 512
	if qr.Design.Size > 0 {
		size = qr.Design.Size
	}
	png, err := utils.GeneratePNGWithLogo(content, size, utils.RenderOptions{LogoURL: qr.Design.LogoURL, Foreground: qr.Design.ForegroundColor, Background: qr.Design.BackgroundColor})
	if err != nil {
		shared.Error(c, 500, "Could not generate QR image", nil)
		return
	}
	c.Header("Content-Disposition", `attachment; filename="qr-code.png"`)
	c.Data(http.StatusOK, "image/png", png)
}

func (h *Handler) GetDesign(c *gin.Context) {
	user, _ := middleware.CurrentUser(c)
	qr, ok := h.findOwnedQR(c, user.ID)
	if !ok {
		return
	}
	shared.OK(c, "Success", qr.Design)
}

func (h *Handler) UpdateDesign(c *gin.Context) {
	user, _ := middleware.CurrentUser(c)
	qr, ok := h.findOwnedQR(c, user.ID)
	if !ok {
		return
	}
	plan, _ := h.currentPlan(user.ID)
	var req DesignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, 400, "Validation error", err.Error())
		return
	}
	if req.LogoURL != "" && !plan.AllowLogo {
		shared.Error(c, 403, "Logo is only available for Pro users", nil)
		return
	}
	if req.LogoURL != "" && req.LogoURL != qr.Design.LogoURL {
		shared.Error(c, 400, "Logo must be uploaded through the logo upload endpoint", nil)
		return
	}
	design := designFromRequest(&req)
	design.LogoURL = qr.Design.LogoURL
	design.LogoPublicID = qr.Design.LogoPublicID
	design.LogoWidth = qr.Design.LogoWidth
	design.LogoHeight = qr.Design.LogoHeight
	design.LogoFormat = qr.Design.LogoFormat
	design.LogoBytes = qr.Design.LogoBytes
	h.db.Model(&models.QRDesign{}).Where("qr_code_id = ?", qr.ID).Updates(design)
	h.db.Where("qr_code_id = ?", qr.ID).First(&design)
	shared.OK(c, "QR design updated", design)
}

func (h *Handler) UploadLogo(c *gin.Context) {
	user, _ := middleware.CurrentUser(c)
	qr, ok := h.findOwnedQR(c, user.ID)
	if !ok {
		return
	}
	plan, _ := h.currentPlan(user.ID)
	if !plan.AllowLogo {
		shared.Error(c, http.StatusForbidden, "Logo upload is only available for Pro users", nil)
		return
	}
	if !h.cfg.CloudinaryEnabled {
		shared.Error(c, http.StatusServiceUnavailable, "Logo upload is currently disabled because Cloudinary is not configured", nil)
		return
	}
	maxBytes := h.cfg.CloudinaryMaxUploadBytes
	if maxBytes <= 0 {
		maxBytes = 5 * 1024 * 1024
	}
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		shared.Error(c, http.StatusBadRequest, "Logo file is required", nil)
		return
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, maxBytes+1))
	if err != nil || int64(len(data)) > maxBytes {
		shared.Error(c, http.StatusBadRequest, fmt.Sprintf("Logo file must not exceed %d bytes", maxBytes), nil)
		return
	}
	contentType := http.DetectContentType(data)
	allowed := map[string]bool{"image/png": true, "image/jpeg": true, "image/webp": true}
	if !allowed[contentType] {
		shared.Error(c, http.StatusBadRequest, "Only PNG, JPEG, and WebP logo files are supported", nil)
		return
	}
	decoded, format, err := image.Decode(bytes.NewReader(data))
	if err != nil || decoded.Bounds().Dx() <= 0 || decoded.Bounds().Dy() <= 0 {
		shared.Error(c, http.StatusBadRequest, "The uploaded file is not a valid image", nil)
		return
	}
	if decoded.Bounds().Dx() > 4096 || decoded.Bounds().Dy() > 4096 {
		shared.Error(c, http.StatusBadRequest, "Logo dimensions must not exceed 4096x4096 pixels", nil)
		return
	}
	if format == "jpeg" {
		format = "jpg"
	}
	asset, err := h.cloudinary.Upload(bytes.NewReader(data), header.Filename, user.ID, qr.ID)
	if err != nil {
		shared.Error(c, http.StatusBadGateway, err.Error(), nil)
		return
	}
	oldPublicID := qr.Design.LogoPublicID
	updates := map[string]any{"logo_url": asset.SecureURL, "logo_public_id": asset.PublicID, "logo_width": decoded.Bounds().Dx(), "logo_height": decoded.Bounds().Dy(), "logo_format": format, "logo_bytes": int64(len(data)), "error_correction_level": shared.ErrorCorrectionLevelH}
	if err := h.db.Model(&models.QRDesign{}).Where("qr_code_id = ?", qr.ID).Updates(updates).Error; err != nil {
		_ = h.cloudinary.Delete(asset.PublicID)
		shared.Error(c, http.StatusInternalServerError, "Could not save logo metadata", nil)
		return
	}
	if oldPublicID != "" && oldPublicID != asset.PublicID {
		if err := h.cloudinary.Delete(oldPublicID); err != nil {
			shared.Error(c, http.StatusInternalServerError, "Logo saved, but the previous Cloudinary asset could not be deleted", nil)
			return
		}
	}
	h.db.Where("qr_code_id = ?", qr.ID).First(&qr.Design)
	shared.OK(c, "Logo uploaded", qr.Design)
}

func (h *Handler) DeleteLogo(c *gin.Context) {
	user, _ := middleware.CurrentUser(c)
	qr, ok := h.findOwnedQR(c, user.ID)
	if !ok {
		return
	}
	if qr.Design.LogoPublicID != "" {
		if err := h.cloudinary.Delete(qr.Design.LogoPublicID); err != nil {
			shared.Error(c, http.StatusBadGateway, "Could not delete logo from Cloudinary", nil)
			return
		}
	}
	updates := map[string]any{"logo_url": "", "logo_public_id": "", "logo_width": 0, "logo_height": 0, "logo_format": "", "logo_bytes": 0, "error_correction_level": shared.ErrorCorrectionLevelM}
	if err := h.db.Model(&models.QRDesign{}).Where("qr_code_id = ?", qr.ID).Updates(updates).Error; err != nil {
		shared.Error(c, http.StatusInternalServerError, "Logo was removed from Cloudinary but database cleanup failed", nil)
		return
	}
	h.db.Where("qr_code_id = ?", qr.ID).First(&qr.Design)
	shared.OK(c, "Logo deleted", qr.Design)
}

func (h *Handler) Redirect(c *gin.Context) {
	var qr models.QRCode
	if err := h.db.Where("short_code = ? AND is_dynamic = ? AND qr_type = ?", c.Param("shortCode"), true, shared.QRTypeURL).First(&qr).Error; err != nil {
		c.String(404, "QR code not found")
		return
	}
	if qr.Status != shared.QRStatusActive || qr.DestinationURL == "" {
		c.String(404, "QR code is not available")
		return
	}
	scan := models.QRScan{
		QRCodeID: qr.ID, ScannedAt: time.Now(), IPAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(), DeviceType: parseDevice(c.Request.UserAgent()),
		Browser: parseBrowser(c.Request.UserAgent()), OperatingSystem: parseOS(c.Request.UserAgent()),
		Referer: c.GetHeader("Referer"),
	}
	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&scan).Error; err != nil {
			return err
		}
		return tx.Model(&qr).UpdateColumn("scan_count", gorm.Expr("scan_count + ?", 1)).Error
	})
	if err != nil {
		c.String(500, "Could not record scan")
		return
	}
	c.Redirect(http.StatusFound, qr.DestinationURL)
}

func (h *Handler) findOwnedQR(c *gin.Context, userID uint) (models.QRCode, bool) {
	var qr models.QRCode
	if err := h.db.Preload("Design").Where("id = ? AND user_id = ? AND status <> ?", c.Param("id"), userID, shared.QRStatusDeleted).First(&qr).Error; err != nil {
		shared.Error(c, 404, "QR code not found", nil)
		return qr, false
	}
	return qr, true
}

func (h *Handler) currentPlan(userID uint) (models.Plan, bool) {
	var subs []models.Subscription
	if err := h.db.Preload("Plan").
		Where("user_id = ? AND status = ? AND end_date > ?", userID, shared.SubscriptionStatusActive, time.Now()).
		Order("end_date desc").
		Find(&subs).Error; err == nil {
		for _, sub := range subs {
			if sub.Plan.Name == shared.PlanNamePro {
				return sub.Plan, true
			}
		}
		if len(subs) > 0 {
			return subs[0].Plan, false
		}
	}
	var free models.Plan
	h.db.Where("name = ?", shared.PlanNameFree).First(&free)
	return free, false
}

func (h *Handler) countUserQRCodes(userID uint) int64 {
	var count int64
	h.db.Model(&models.QRCode{}).Where("user_id = ? AND status <> ?", userID, shared.QRStatusDeleted).Count(&count)
	return count
}

func (h *Handler) validateQRInput(qrType shared.QRType, content string) error {
	switch qrType {
	case shared.QRTypeURL, shared.QRTypeSocial, shared.QRTypePDF, shared.QRTypeMenu:
		if !isValidURL(content) {
			return errString("content must be a valid URL")
		}
	}
	return nil
}

func designFromRequest(req *DesignRequest) models.QRDesign {
	design := models.QRDesign{ForegroundColor: "#111827", BackgroundColor: "#FFFFFF", Size: 512, ErrorCorrectionLevel: shared.ErrorCorrectionLevelM}
	if req == nil {
		return design
	}
	if req.ForegroundColor != "" {
		design.ForegroundColor = req.ForegroundColor
	}
	if req.BackgroundColor != "" {
		design.BackgroundColor = req.BackgroundColor
	}
	if req.Size > 0 {
		design.Size = req.Size
	}
	if req.ErrorCorrectionLevel != "" {
		design.ErrorCorrectionLevel = req.ErrorCorrectionLevel
	}
	design.TemplateID = req.TemplateID
	design.EyeStyle = req.EyeStyle
	design.DotStyle = req.DotStyle
	design.FrameStyle = req.FrameStyle
	design.LogoURL = req.LogoURL
	return design
}

func isValidURL(value string) bool {
	parsed, err := url.ParseRequestURI(value)
	return err == nil && parsed.Scheme != "" && parsed.Host != ""
}

type errString string

func (e errString) Error() string { return string(e) }

func parseDevice(ua string) string {
	lower := strings.ToLower(ua)
	if strings.Contains(lower, "mobile") || strings.Contains(lower, "android") || strings.Contains(lower, "iphone") {
		return "Mobile"
	}
	if strings.Contains(lower, "ipad") || strings.Contains(lower, "tablet") {
		return "Tablet"
	}
	return "Desktop"
}

func parseBrowser(ua string) string {
	lower := strings.ToLower(ua)
	switch {
	case strings.Contains(lower, "edg"):
		return "Edge"
	case strings.Contains(lower, "chrome"):
		return "Chrome"
	case strings.Contains(lower, "firefox"):
		return "Firefox"
	case strings.Contains(lower, "safari"):
		return "Safari"
	default:
		return "Other"
	}
}

func parseOS(ua string) string {
	lower := strings.ToLower(ua)
	switch {
	case strings.Contains(lower, "windows"):
		return "Windows"
	case strings.Contains(lower, "android"):
		return "Android"
	case strings.Contains(lower, "iphone") || strings.Contains(lower, "ipad") || strings.Contains(lower, "mac os"):
		return "Apple"
	case strings.Contains(lower, "linux"):
		return "Linux"
	default:
		return "Other"
	}
}
