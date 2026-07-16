package config

import (
	"os"
	"strconv"
	"time"

	mysqlcfg "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv             string
	AppPort            string
	AppURL             string
	FrontendURL        string
	AdminFrontendURL   string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	JWTAccessSecret    string
	JWTRefreshSecret   string
	JWTAccessTTLMin    int
	JWTRefreshTTLHours int
	FreeMaxQRCodes     int
	AdminEmail         string
	AdminPassword      string
	SepayEnabled           bool
	SepayWebhookSecret     string
	SepayTransactionPrefix string
	SepayReturnURL         string
	SepayAPIURL            string
	SepayAPIKey            string
	BankCode               string
	AccountNo              string
	AccountName            string
}

func Load() Config {
	loadEnv()

	return Config{
		AppEnv:             getEnv("APP_ENV", "development"),
		AppPort:            getEnv("APP_PORT", "8080"),
		AppURL:             getEnv("APP_URL", "http://localhost:8080"),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:3000"),
		AdminFrontendURL:   getEnv("ADMIN_FRONTEND_URL", "http://localhost:5173"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "3306"),
		DBUser:             getEnv("DB_USER", "qr_user"),
		DBPassword:         getEnv("DB_PASSWORD", "qr_password"),
		DBName:             getEnv("DB_NAME", "qr_generator"),
		JWTAccessSecret:    getEnv("JWT_ACCESS_SECRET", "change-me-access-secret"),
		JWTRefreshSecret:   getEnv("JWT_REFRESH_SECRET", "change-me-refresh-secret"),
		JWTAccessTTLMin:    getEnvInt("JWT_ACCESS_TTL_MINUTES", 60),
		JWTRefreshTTLHours: getEnvInt("JWT_REFRESH_TTL_HOURS", 720),
		FreeMaxQRCodes:     getEnvInt("FREE_MAX_QR_CODES", 10),
		AdminEmail:         getEnv("ADMIN_EMAIL", "admin@qr.local"),
		AdminPassword:      getEnv("ADMIN_PASSWORD", "Admin@123456"),
		SepayEnabled:           getEnvBool("SEPAY_ENABLED", true),
		SepayWebhookSecret:     getEnv("SEPAY_WEBHOOK_SECRET", ""),
		SepayTransactionPrefix: getEnv("SEPAY_TRANSACTION_PREFIX", "QRPRO"),
		SepayReturnURL:         getEnv("SEPAY_RETURN_URL", getEnv("FRONTEND_URL", "http://localhost:3000")+"/pricing"),
		SepayAPIURL:            getEnv("SEPAY_API_URL", ""),
		SepayAPIKey:            getEnv("SEPAY_API_KEY", ""),
		BankCode:               getEnv("BANK_CODE", ""),
		AccountNo:              getEnv("ACCOUNT_NO", ""),
		AccountName:            getEnv("ACCOUNT_NAME", ""),
	}
}

func loadEnv() {
	if os.Getenv("APP_ENV") == "docker" && os.Getenv("DB_HOST") == "mysql" {
		_ = godotenv.Load(".env", "backend/.env")
		return
	}
	_ = godotenv.Overload(".env", "backend/.env")
}

func (c Config) DSN() string {
	cfg := mysqlcfg.Config{
		User:                 c.DBUser,
		Passwd:               c.DBPassword,
		Net:                  "tcp",
		Addr:                 c.DBHost + ":" + c.DBPort,
		DBName:               c.DBName,
		AllowNativePasswords: true,
		ParseTime:            true,
		Loc:                  time.Local,
		Params: map[string]string{
			"charset": "utf8mb4",
		},
	}
	return cfg.FormatDSN()
}

func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value, err := strconv.Atoi(getEnv(key, ""))
	if err != nil {
		return fallback
	}
	return value
}

func getEnvBool(key string, fallback bool) bool {
	value := getEnv(key, "")
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}
