package shared

type UserStatus string
type RoleName string
type PlanName string
type PlanStatus string
type SubscriptionStatus string
type PaymentStatus string
type PaymentMethod string
type QRType string
type QRStatus string
type TemplateStatus string
type ErrorCorrectionLevel string
type LogLevel string

const (
	UserStatusActive  UserStatus = "ACTIVE"
	UserStatusLocked  UserStatus = "LOCKED"
	UserStatusDeleted UserStatus = "DELETED"

	RoleNameUser  RoleName = "USER"
	RoleNameAdmin RoleName = "ADMIN"

	PlanNameFree PlanName = "FREE"
	PlanNamePro  PlanName = "PRO"

	PlanStatusActive   PlanStatus = "ACTIVE"
	PlanStatusInactive PlanStatus = "INACTIVE"
	PlanStatusDeleted  PlanStatus = "DELETED"

	SubscriptionStatusActive    SubscriptionStatus = "ACTIVE"
	SubscriptionStatusExpired   SubscriptionStatus = "EXPIRED"
	SubscriptionStatusCancelled SubscriptionStatus = "CANCELLED"
	SubscriptionStatusPending   SubscriptionStatus = "PENDING"

	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusSuccess   PaymentStatus = "SUCCESS"
	PaymentStatusFailed    PaymentStatus = "FAILED"
	PaymentStatusCancelled PaymentStatus = "CANCELLED"
	PaymentStatusRefunded  PaymentStatus = "REFUNDED"

	PaymentMethodVNPay        PaymentMethod = "VNPAY"
	PaymentMethodMomo         PaymentMethod = "MOMO"
	PaymentMethodZaloPay      PaymentMethod = "ZALOPAY"
	PaymentMethodPayPal       PaymentMethod = "PAYPAL"
	PaymentMethodStripe       PaymentMethod = "STRIPE"
	PaymentMethodBankTransfer PaymentMethod = "BANK_TRANSFER"
	PaymentMethodSepay        PaymentMethod = "SEPAY"

	QRTypeURL      QRType = "URL"
	QRTypeText     QRType = "TEXT"
	QRTypeWiFi     QRType = "WIFI"
	QRTypeVCard    QRType = "VCARD"
	QRTypeEmail    QRType = "EMAIL"
	QRTypeSMS      QRType = "SMS"
	QRTypeLocation QRType = "LOCATION"
	QRTypeSocial   QRType = "SOCIAL"
	QRTypePDF      QRType = "PDF"
	QRTypeMenu     QRType = "MENU"

	QRStatusActive   QRStatus = "ACTIVE"
	QRStatusDisabled QRStatus = "DISABLED"
	QRStatusDeleted  QRStatus = "DELETED"

	TemplateStatusActive  TemplateStatus = "ACTIVE"
	TemplateStatusHidden  TemplateStatus = "HIDDEN"
	TemplateStatusDeleted TemplateStatus = "DELETED"

	ErrorCorrectionLevelL ErrorCorrectionLevel = "L"
	ErrorCorrectionLevelM ErrorCorrectionLevel = "M"
	ErrorCorrectionLevelQ ErrorCorrectionLevel = "Q"
	ErrorCorrectionLevelH ErrorCorrectionLevel = "H"

	LogLevelInfo     LogLevel = "INFO"
	LogLevelWarning  LogLevel = "WARNING"
	LogLevelError    LogLevel = "ERROR"
	LogLevelSecurity LogLevel = "SECURITY"
)
