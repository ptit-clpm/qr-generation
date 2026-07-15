package payments

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"regexp"
	"strings"
	"time"

	"qr-generator/backend/internal/config"
	"qr-generator/backend/internal/middleware"
	"qr-generator/backend/internal/models"
	"qr-generator/backend/internal/shared"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var errPaymentNotPending = errors.New("payment is not pending")

type CreatePaymentRequest struct {
	PlanID        uint                 `json:"plan_id" binding:"required"`
	PaymentMethod shared.PaymentMethod `json:"payment_method"`
}

type CreatePaymentResponse struct {
	Payment      models.Payment   `json:"payment"`
	Instructions SepayPaymentInfo `json:"instructions"`
}

type SepayPaymentInfo struct {
	Provider        string  `json:"provider"`
	BankAccount     string  `json:"bank_account"`
	BankName        string  `json:"bank_name"`
	AccountName     string  `json:"account_name"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	TransactionCode string  `json:"transaction_code"`
	TransferContent string  `json:"transfer_content"`
	ReturnURL       string  `json:"return_url"`
	Enabled         bool    `json:"enabled"`
}

type SepayWebhookRequest struct {
	Secret          string  `json:"secret"`
	TransactionCode string  `json:"transaction_code"`
	Code            string  `json:"code"`
	Content         string  `json:"content"`
	Description     string  `json:"description"`
	ReferenceCode   string  `json:"referenceCode"`
	TransferType    string  `json:"transferType"`
	TransferAmount  float64 `json:"transferAmount"`
	Amount          float64 `json:"amount"`
	Status          string  `json:"status"`
}

type MockSuccessRequest struct {
	TransactionCode string `json:"transaction_code" binding:"required"`
}

type Handler struct {
	db  *gorm.DB
	cfg config.Config
}

func RegisterPublicRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg config.Config) {
	h := &Handler{db: db, cfg: cfg}
	rg.POST("/payments/sepay/webhook", h.SepayWebhook)
}

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg config.Config) {
	h := &Handler{db: db, cfg: cfg}
	payments := rg.Group("/payments")
	payments.POST("/create", h.Create)
	if cfg.AppEnv == "development" && !cfg.SepayEnabled {
		payments.POST("/mock-success", h.MockSuccess)
	}
	payments.GET("/:transactionCode", h.Detail)
	payments.GET("", h.List)
	rg.POST("/subscriptions/upgrade", h.Create)
}

func (h *Handler) Create(c *gin.Context) {
	user, _ := middleware.CurrentUser(c)
	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, 400, "Validation error", err.Error())
		return
	}
	var plan models.Plan
	if err := h.db.First(&plan, req.PlanID).Error; err != nil || plan.Status != shared.PlanStatusActive || plan.Name != shared.PlanNamePro {
		shared.Error(c, 404, "Active Pro plan not found", nil)
		return
	}
	method := req.PaymentMethod
	if method == "" {
		method = shared.PaymentMethodSepay
	}
	payment := models.Payment{
		UserID:          user.ID,
		Amount:          plan.Price,
		Currency:        "VND",
		PaymentMethod:   method,
		TransactionCode: h.newTransactionCode(user.ID),
		Provider:        "SEPAY",
		Status:          shared.PaymentStatusPending,
	}
	if err := h.db.Create(&payment).Error; err != nil {
		shared.Error(c, 500, "Could not create payment", nil)
		return
	}
	shared.Created(c, "Payment created", CreatePaymentResponse{
		Payment:      payment,
		Instructions: h.sepayInfo(payment),
	})
}

func (h *Handler) SepayWebhook(c *gin.Context) {
	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		shared.Error(c, 400, "Could not read webhook payload", nil)
		return
	}
	var req SepayWebhookRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		shared.Error(c, 400, "Invalid webhook payload", err.Error())
		return
	}
	if !h.validWebhookSecret(c, req.Secret) {
		shared.Error(c, 401, "Invalid webhook secret", nil)
		return
	}
	transactionCode := h.extractTransactionCode(req)
	if transactionCode == "" {
		shared.Error(c, 400, "Transaction code not found", nil)
		return
	}
	if status, ok := webhookTerminalFailure(req.Status); ok {
		h.markPaymentFailed(c, transactionCode, status, raw, req.ReferenceCode)
		return
	}
	if strings.EqualFold(req.TransferType, "out") {
		shared.Error(c, 400, "Only incoming Sepay transfers are accepted", nil)
		return
	}

	var payment models.Payment
	if err := h.db.Where("transaction_code = ?", transactionCode).First(&payment).Error; err != nil {
		shared.Error(c, 404, "Payment not found", nil)
		return
	}
	if payment.Status == shared.PaymentStatusSuccess {
		shared.OK(c, "Payment already activated", payment)
		return
	}
	if payment.Status != shared.PaymentStatusPending {
		shared.Error(c, 409, "Payment is not pending", nil)
		return
	}
	amount := req.paymentAmount()
	if amount <= 0 || math.Abs(amount-payment.Amount) > 0.01 {
		shared.Error(c, 400, "Invalid payment amount", nil)
		return
	}

	now := time.Now()
	err = h.db.Transaction(func(tx *gorm.DB) error {
		var locked models.Payment
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("transaction_code = ?", transactionCode).First(&locked).Error; err != nil {
			return err
		}
		if locked.Status == shared.PaymentStatusSuccess {
			payment = locked
			return nil
		}
		if locked.Status != shared.PaymentStatusPending {
			return errPaymentNotPending
		}
		var pro models.Plan
		if err := tx.Where("name = ? AND status = ?", shared.PlanNamePro, shared.PlanStatusActive).First(&pro).Error; err != nil {
			return err
		}
		sub, err := h.createProSubscription(tx, locked.UserID, pro, now)
		if err != nil {
			return err
		}
		locked.Status = shared.PaymentStatusSuccess
		locked.PaidAt = &now
		locked.SubscriptionID = &sub.ID
		locked.ProviderRef = req.ReferenceCode
		locked.ProviderPayload = string(raw)
		if err := tx.Save(&locked).Error; err != nil {
			return err
		}
		payment = locked
		return nil
	})
	if err != nil {
		if errors.Is(err, errPaymentNotPending) {
			shared.Error(c, 409, "Payment is not pending", nil)
			return
		}
		shared.Error(c, 500, "Could not activate payment", nil)
		return
	}
	shared.OK(c, "Payment success and Pro activated", payment)
}

func (h *Handler) MockSuccess(c *gin.Context) {
	user, _ := middleware.CurrentUser(c)
	var req MockSuccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, 400, "Validation error", err.Error())
		return
	}
	var payment models.Payment
	if err := h.db.Where("transaction_code = ? AND user_id = ?", req.TransactionCode, user.ID).First(&payment).Error; err != nil {
		shared.Error(c, 404, "Payment not found", nil)
		return
	}
	if payment.Status == shared.PaymentStatusSuccess {
		shared.OK(c, "Payment already activated", payment)
		return
	}
	now := time.Now()
	err := h.db.Transaction(func(tx *gorm.DB) error {
		var pro models.Plan
		if err := tx.Where("name = ? AND status = ?", shared.PlanNamePro, shared.PlanStatusActive).First(&pro).Error; err != nil {
			return err
		}
		sub, err := h.createProSubscription(tx, payment.UserID, pro, now)
		if err != nil {
			return err
		}
		payment.Status = shared.PaymentStatusSuccess
		payment.PaidAt = &now
		payment.SubscriptionID = &sub.ID
		return tx.Save(&payment).Error
	})
	if err != nil {
		shared.Error(c, 500, "Could not activate payment", nil)
		return
	}
	shared.OK(c, "Payment success and Pro activated", payment)
}

func (h *Handler) List(c *gin.Context) {
	user, _ := middleware.CurrentUser(c)
	var payments []models.Payment
	h.db.Where("user_id = ?", user.ID).Order("created_at desc").Find(&payments)
	shared.OK(c, "Success", payments)
}

func (h *Handler) Detail(c *gin.Context) {
	user, _ := middleware.CurrentUser(c)
	var payment models.Payment
	if err := h.db.Where("transaction_code = ? AND user_id = ?", c.Param("transactionCode"), user.ID).First(&payment).Error; err != nil {
		shared.Error(c, 404, "Payment not found", nil)
		return
	}
	shared.OK(c, "Success", CreatePaymentResponse{
		Payment:      payment,
		Instructions: h.sepayInfo(payment),
	})
}

func (h *Handler) newTransactionCode(userID uint) string {
	prefix := strings.TrimSpace(h.cfg.SepayPaymentPrefix)
	if prefix == "" {
		prefix = "QRPRO"
	}
	return fmt.Sprintf("%s-%d-%s", prefix, userID, strings.ToUpper(uuid.NewString()[:8]))
}

func (h *Handler) sepayInfo(payment models.Payment) SepayPaymentInfo {
	return SepayPaymentInfo{
		Provider:        "SEPAY",
		BankAccount:     h.cfg.SepayBankAccount,
		BankName:        h.cfg.SepayBankName,
		AccountName:     h.cfg.SepayAccountName,
		Amount:          payment.Amount,
		Currency:        payment.Currency,
		TransactionCode: payment.TransactionCode,
		TransferContent: payment.TransactionCode,
		ReturnURL:       h.cfg.SepayReturnURL,
		Enabled:         h.cfg.SepayEnabled,
	}
}

func (h *Handler) createProSubscription(tx *gorm.DB, userID uint, pro models.Plan, now time.Time) (models.Subscription, error) {
	start := now
	var latest models.Subscription
	err := tx.Joins("JOIN plans ON plans.id = subscriptions.plan_id").
		Where("subscriptions.user_id = ? AND subscriptions.status = ? AND subscriptions.end_date > ? AND plans.name = ?",
			userID, shared.SubscriptionStatusActive, now, shared.PlanNamePro).
		Order("subscriptions.end_date desc").
		First(&latest).Error
	if err == nil && latest.EndDate.After(start) {
		start = latest.EndDate
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Subscription{}, err
	}
	sub := models.Subscription{
		UserID:    userID,
		PlanID:    pro.ID,
		StartDate: start,
		EndDate:   start.AddDate(0, 0, pro.DurationDays),
		Status:    shared.SubscriptionStatusActive,
	}
	if err := tx.Create(&sub).Error; err != nil {
		return models.Subscription{}, err
	}
	return sub, nil
}

func (h *Handler) validWebhookSecret(c *gin.Context, bodySecret string) bool {
	expected := strings.TrimSpace(h.cfg.SepayWebhookSecret)
	if expected == "" {
		return false
	}
	candidates := []string{
		c.GetHeader("X-Sepay-Signature"),
		c.GetHeader("X-Sepay-Webhook-Secret"),
		c.GetHeader("X-Webhook-Secret"),
		bodySecret,
	}
	if auth := c.GetHeader("Authorization"); strings.HasPrefix(strings.ToLower(auth), "bearer ") {
		candidates = append(candidates, strings.TrimSpace(auth[7:]))
	}
	for _, candidate := range candidates {
		candidate = strings.TrimSpace(candidate)
		if candidate != "" && subtle.ConstantTimeCompare([]byte(candidate), []byte(expected)) == 1 {
			return true
		}
	}
	return false
}

func (h *Handler) extractTransactionCode(req SepayWebhookRequest) string {
	prefix := h.cfg.SepayPaymentPrefix
	if prefix == "" {
		prefix = "QRPRO"
	}
	for _, candidate := range []string{req.TransactionCode, req.Code} {
		candidate = strings.TrimSpace(candidate)
		if strings.HasPrefix(strings.ToUpper(candidate), strings.ToUpper(prefix)) {
			return candidate
		}
	}
	pattern := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(prefix) + `[-_A-Z0-9]*`)
	if match := pattern.FindString(req.Content + " " + req.Description); match != "" {
		return match
	}
	return ""
}

func (req SepayWebhookRequest) paymentAmount() float64 {
	if req.TransferAmount > 0 {
		return req.TransferAmount
	}
	return req.Amount
}

func webhookTerminalFailure(status string) (shared.PaymentStatus, bool) {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "FAILED", "FAIL":
		return shared.PaymentStatusFailed, true
	case "CANCELLED", "CANCELED":
		return shared.PaymentStatusCancelled, true
	default:
		return "", false
	}
}

func (h *Handler) markPaymentFailed(c *gin.Context, transactionCode string, status shared.PaymentStatus, raw []byte, reference string) {
	var payment models.Payment
	if err := h.db.Where("transaction_code = ?", transactionCode).First(&payment).Error; err != nil {
		shared.Error(c, 404, "Payment not found", nil)
		return
	}
	if payment.Status == shared.PaymentStatusSuccess {
		shared.OK(c, "Payment already activated", payment)
		return
	}
	if payment.Status == shared.PaymentStatusPending {
		payment.Status = status
		payment.ProviderRef = reference
		payment.ProviderPayload = string(raw)
		if err := h.db.Save(&payment).Error; err != nil {
			shared.Error(c, 500, "Could not update payment", nil)
			return
		}
	}
	shared.OK(c, "Payment status updated", payment)
}
