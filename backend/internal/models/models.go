package models

import (
	"time"

	"qr-generator/backend/internal/shared"
)

type User struct {
	ID           uint              `gorm:"primaryKey" json:"id"`
	FullName     string            `gorm:"size:150;not null" json:"full_name"`
	Email        string            `gorm:"size:150;uniqueIndex;not null" json:"email"`
	PasswordHash string            `gorm:"size:255;not null" json:"-"`
	PhoneNumber  string            `gorm:"size:30" json:"phone_number"`
	AvatarURL    string            `gorm:"size:255" json:"avatar_url"`
	Status       shared.UserStatus `gorm:"size:20;not null;default:ACTIVE" json:"status"`
	Roles        []Role            `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

type Role struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Name        shared.RoleName `gorm:"size:30;uniqueIndex;not null" json:"name"`
	Description string          `gorm:"size:255" json:"description"`
}

type Plan struct {
	ID                uint              `gorm:"primaryKey" json:"id"`
	Name              shared.PlanName   `gorm:"size:30;uniqueIndex;not null" json:"name"`
	Price             float64           `gorm:"not null;default:0" json:"price"`
	DurationDays      int               `gorm:"not null;default:30" json:"duration_days"`
	MaxQRCodes        int               `gorm:"not null;default:10" json:"max_qr_codes"`
	AllowDynamicQR    bool              `gorm:"not null;default:false" json:"allow_dynamic_qr"`
	AllowLogo         bool              `gorm:"not null;default:false" json:"allow_logo"`
	AllowAnalytics    bool              `gorm:"not null;default:false" json:"allow_analytics"`
	AllowSVGPDFExport bool              `gorm:"not null;default:false" json:"allow_svg_pdf_export"`
	Description       string            `gorm:"type:text" json:"description"`
	Status            shared.PlanStatus `gorm:"size:20;not null;default:ACTIVE" json:"status"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

type Subscription struct {
	ID        uint                      `gorm:"primaryKey" json:"id"`
	UserID    uint                      `gorm:"not null;index" json:"user_id"`
	PlanID    uint                      `gorm:"not null;index" json:"plan_id"`
	User      User                      `json:"user,omitempty"`
	Plan      Plan                      `json:"plan,omitempty"`
	StartDate time.Time                 `gorm:"not null" json:"start_date"`
	EndDate   time.Time                 `gorm:"not null" json:"end_date"`
	Status    shared.SubscriptionStatus `gorm:"size:20;not null;default:PENDING" json:"status"`
	AutoRenew bool                      `gorm:"not null;default:false" json:"auto_renew"`
	CreatedAt time.Time                 `json:"created_at"`
	UpdatedAt time.Time                 `json:"updated_at"`
}

type Payment struct {
	ID              uint                 `gorm:"primaryKey" json:"id"`
	UserID          uint                 `gorm:"not null;index" json:"user_id"`
	SubscriptionID  *uint                `gorm:"index" json:"subscription_id"`
	User            User                 `json:"user,omitempty"`
	Subscription    *Subscription        `json:"subscription,omitempty"`
	Amount          float64              `gorm:"not null" json:"amount"`
	Currency        string               `gorm:"size:10;not null;default:VND" json:"currency"`
	PaymentMethod   shared.PaymentMethod `gorm:"size:30;not null" json:"payment_method"`
	TransactionCode string               `gorm:"size:100;uniqueIndex" json:"transaction_code"`
	Provider        string               `gorm:"size:50" json:"provider"`
	ProviderRef     string               `gorm:"size:150" json:"provider_ref"`
	ProviderPayload string               `gorm:"type:text" json:"provider_payload"`
	Status          shared.PaymentStatus `gorm:"size:20;not null;default:PENDING" json:"status"`
	PaidAt          *time.Time           `json:"paid_at"`
	CreatedAt       time.Time            `json:"created_at"`
}

type Folder struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type QRCode struct {
	ID             uint            `gorm:"primaryKey" json:"id"`
	UserID         uint            `gorm:"not null;index" json:"user_id"`
	FolderID       *uint           `gorm:"index" json:"folder_id"`
	Title          string          `gorm:"size:150;not null" json:"title"`
	QRType         shared.QRType   `gorm:"size:30;not null" json:"qr_type"`
	Content        string          `gorm:"type:text;not null" json:"content"`
	ShortCode      *string         `gorm:"size:30;uniqueIndex" json:"short_code,omitempty"`
	IsDynamic      bool            `gorm:"not null;default:false" json:"is_dynamic"`
	DestinationURL string          `gorm:"size:2048" json:"destination_url"`
	QRImageURL     string          `gorm:"size:255" json:"qr_image_url"`
	ScanCount      int64           `gorm:"not null;default:0" json:"scan_count"`
	Status         shared.QRStatus `gorm:"size:20;not null;default:ACTIVE" json:"status"`
	Design         QRDesign        `json:"design,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type QRDesign struct {
	ID                   uint                        `gorm:"primaryKey" json:"id"`
	QRCodeID             uint                        `gorm:"uniqueIndex;not null" json:"qr_code_id"`
	TemplateID           *uint                       `gorm:"index" json:"template_id"`
	ForegroundColor      string                      `gorm:"size:20;not null;default:#111827" json:"foreground_color"`
	BackgroundColor      string                      `gorm:"size:20;not null;default:#FFFFFF" json:"background_color"`
	EyeStyle             string                      `gorm:"size:50" json:"eye_style"`
	DotStyle             string                      `gorm:"size:50" json:"dot_style"`
	FrameStyle           string                      `gorm:"size:50" json:"frame_style"`
	LogoURL              string                      `gorm:"size:255" json:"logo_url"`
	Size                 int                         `gorm:"not null;default:512" json:"size"`
	ErrorCorrectionLevel shared.ErrorCorrectionLevel `gorm:"size:5;not null;default:M" json:"error_correction_level"`
	CreatedAt            time.Time                   `json:"created_at"`
	UpdatedAt            time.Time                   `json:"updated_at"`
}

type QRTemplate struct {
	ID              uint                  `gorm:"primaryKey" json:"id"`
	Name            string                `gorm:"size:100;not null" json:"name"`
	PreviewImageURL string                `gorm:"size:255" json:"preview_image_url"`
	ConfigJSON      string                `gorm:"type:text" json:"config_json"`
	IsPro           bool                  `gorm:"not null;default:false" json:"is_pro"`
	Status          shared.TemplateStatus `gorm:"size:20;not null;default:ACTIVE" json:"status"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
}

type QRScan struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	QRCodeID        uint      `gorm:"not null;index" json:"qr_code_id"`
	ScannedAt       time.Time `gorm:"not null;index" json:"scanned_at"`
	IPAddress       string    `gorm:"size:100" json:"ip_address"`
	UserAgent       string    `gorm:"type:text" json:"user_agent"`
	DeviceType      string    `gorm:"size:50" json:"device_type"`
	Browser         string    `gorm:"size:50" json:"browser"`
	OperatingSystem string    `gorm:"size:50" json:"operating_system"`
	Country         string    `gorm:"size:100" json:"country"`
	City            string    `gorm:"size:100" json:"city"`
	Referer         string    `gorm:"size:255" json:"referer"`
}

type SystemLog struct {
	ID         uint            `gorm:"primaryKey" json:"id"`
	UserID     *uint           `gorm:"index" json:"user_id"`
	Action     string          `gorm:"size:100;not null" json:"action"`
	EntityType string          `gorm:"size:100" json:"entity_type"`
	EntityID   *uint           `json:"entity_id"`
	Level      shared.LogLevel `gorm:"size:20;not null;default:INFO" json:"level"`
	Message    string          `gorm:"type:text" json:"message"`
	IPAddress  string          `gorm:"size:100" json:"ip_address"`
	CreatedAt  time.Time       `json:"created_at"`
}
