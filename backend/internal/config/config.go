package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	mysqlcfg "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv                   string
	AppPort                  string
	AppURL                   string
	FrontendURL              string
	AdminFrontendURL         string
	DBHost                   string
	DBPort                   string
	DBUser                   string
	DBPassword               string
	DBName                   string
	DatabaseURL              string
	DBSSLCA                  string
	JWTAccessSecret          string
	JWTRefreshSecret         string
	JWTAccessTTLMin          int
	JWTRefreshTTLHours       int
	FreeMaxQRCodes           int
	AdminEmail               string
	AdminPassword            string
	SepayEnabled             bool
	SepayWebhookSecret       string
	SepayTransactionPrefix   string
	SepayReturnURL           string
	SepayAPIURL              string
	SepayAPIKey              string
	BankCode                 string
	AccountNo                string
	AccountName              string
	CloudinaryEnabled        bool
	CloudinaryCloudName      string
	CloudinaryAPIKey         string
	CloudinaryAPISecret      string
	CloudinaryFolder         string
	CloudinaryMaxUploadBytes int64
}

func Load() Config {
	loadEnv()

	return Config{
		AppEnv:                   getEnv("APP_ENV", "development"),
		AppPort:                  getEnv("APP_PORT", "8080"),
		AppURL:                   getEnv("APP_URL", "http://localhost:8080"),
		FrontendURL:              getEnv("FRONTEND_URL", "http://localhost:3000"),
		AdminFrontendURL:         getEnv("ADMIN_FRONTEND_URL", "http://localhost:5173"),
		DBHost:                   getEnv("DB_HOST", "localhost"),
		DBPort:                   getEnv("DB_PORT", "3306"),
		DBUser:                   getEnv("DB_USER", "qr_user"),
		DBPassword:               getEnv("DB_PASSWORD", "qr_password"),
		DBName:                   getEnv("DB_NAME", "qr_generator"),
		DatabaseURL:              getEnv("DATABASE_URL", ""),
		DBSSLCA:                  getEnv("DB_SSL_CA", ""),
		JWTAccessSecret:          getEnv("JWT_ACCESS_SECRET", "change-me-access-secret"),
		JWTRefreshSecret:         getEnv("JWT_REFRESH_SECRET", "change-me-refresh-secret"),
		JWTAccessTTLMin:          getEnvInt("JWT_ACCESS_TTL_MINUTES", 60),
		JWTRefreshTTLHours:       getEnvInt("JWT_REFRESH_TTL_HOURS", 720),
		FreeMaxQRCodes:           getEnvInt("FREE_MAX_QR_CODES", 10),
		AdminEmail:               getEnv("ADMIN_EMAIL", "admin@qr.local"),
		AdminPassword:            getEnv("ADMIN_PASSWORD", "Admin@123456"),
		SepayEnabled:             getEnvBool("SEPAY_ENABLED", true),
		SepayWebhookSecret:       getEnv("SEPAY_WEBHOOK_SECRET", ""),
		SepayTransactionPrefix:   getEnv("SEPAY_TRANSACTION_PREFIX", "QRPRO"),
		SepayReturnURL:           getEnv("SEPAY_RETURN_URL", getEnv("FRONTEND_URL", "http://localhost:3000")+"/pricing"),
		SepayAPIURL:              getEnv("SEPAY_API_URL", ""),
		SepayAPIKey:              getEnv("SEPAY_API_KEY", ""),
		BankCode:                 getEnv("BANK_CODE", ""),
		AccountNo:                getEnv("ACCOUNT_NO", ""),
		AccountName:              getEnv("ACCOUNT_NAME", ""),
		CloudinaryEnabled:        getEnvBool("CLOUDINARY_ENABLED", false),
		CloudinaryCloudName:      getEnv("CLOUDINARY_CLOUD_NAME", ""),
		CloudinaryAPIKey:         getEnv("CLOUDINARY_API_KEY", ""),
		CloudinaryAPISecret:      getEnv("CLOUDINARY_API_SECRET", ""),
		CloudinaryFolder:         getEnv("CLOUDINARY_FOLDER", "qr-generator/logos"),
		CloudinaryMaxUploadBytes: int64(getEnvInt("CLOUDINARY_MAX_UPLOAD_BYTES", 5*1024*1024)),
	}
}

func (c Config) ValidateForProduction() error {
	if c.AppEnv != "production" {
		return nil
	}
	if len(strings.TrimSpace(c.JWTAccessSecret)) < 32 || len(strings.TrimSpace(c.JWTRefreshSecret)) < 32 {
		return fmt.Errorf("JWT secrets must be at least 32 characters in production")
	}
	if c.SepayEnabled && strings.TrimSpace(c.SepayWebhookSecret) == "" {
		return fmt.Errorf("SEPAY_WEBHOOK_SECRET is required when Sepay is enabled in production")
	}
	if strings.TrimSpace(c.FrontendURL) == "" || strings.TrimSpace(c.AdminFrontendURL) == "" {
		return fmt.Errorf("FRONTEND_URL and ADMIN_FRONTEND_URL are required in production")
	}
	if c.CloudinaryEnabled && (strings.TrimSpace(c.CloudinaryCloudName) == "" || strings.TrimSpace(c.CloudinaryAPIKey) == "" || strings.TrimSpace(c.CloudinaryAPISecret) == "") {
		return fmt.Errorf("CLOUDINARY_CLOUD_NAME, CLOUDINARY_API_KEY, and CLOUDINARY_API_SECRET are required when Cloudinary is enabled")
	}
	return nil
}

func loadEnv() {
	if os.Getenv("APP_ENV") == "docker" && os.Getenv("DB_HOST") == "mysql" {
		_ = godotenv.Load(".env", "backend/.env")
		return
	}
	_ = godotenv.Overload(".env", "backend/.env")
}

func (c Config) DSN() (string, error) {
	if strings.TrimSpace(c.DatabaseURL) != "" {
		return mysqlDSNFromURL(c.DatabaseURL, c.DBSSLCA)
	}

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
	return cfg.FormatDSN(), nil
}

func mysqlDSNFromURL(raw string, caPath string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", fmt.Errorf("invalid DATABASE_URL: %w", err)
	}
	if u.Scheme != "mysql" {
		return "", fmt.Errorf("invalid DATABASE_URL scheme %q: expected mysql", u.Scheme)
	}
	if u.Host == "" || u.User == nil {
		return "", fmt.Errorf("invalid DATABASE_URL: host and user are required")
	}

	password, _ := u.User.Password()
	cfg := mysqlcfg.Config{
		User:                 u.User.Username(),
		Passwd:               password,
		Net:                  "tcp",
		Addr:                 u.Host,
		DBName:               strings.TrimPrefix(u.Path, "/"),
		AllowNativePasswords: true,
		ParseTime:            true,
		Loc:                  time.Local,
		Params: map[string]string{
			"charset": "utf8mb4",
		},
	}

	query := u.Query()
	if charset := query.Get("charset"); charset != "" {
		cfg.Params["charset"] = charset
	}
	sslMode := strings.ToLower(query.Get("ssl-mode"))
	if sslMode == "required" || sslMode == "verify-ca" || sslMode == "verify-identity" {
		cfg.TLSConfig = "true"
	}
	if tls := query.Get("tls"); tls != "" {
		cfg.TLSConfig = tls
	}
	if strings.TrimSpace(caPath) != "" {
		caPEM, err := os.ReadFile(os.ExpandEnv(strings.TrimSpace(caPath)))
		if err != nil {
			return "", fmt.Errorf("read DB_SSL_CA %q: %w", caPath, err)
		}
		rootCAs := x509.NewCertPool()
		if !rootCAs.AppendCertsFromPEM(caPEM) {
			return "", fmt.Errorf("DB_SSL_CA %q does not contain a valid PEM certificate", caPath)
		}
		if err := mysqlcfg.RegisterTLSConfig("aiven", &tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs:    rootCAs,
			ServerName: u.Hostname(),
		}); err != nil {
			return "", fmt.Errorf("register Aiven TLS config: %w", err)
		}
		cfg.TLSConfig = "aiven"
	}

	return cfg.FormatDSN(), nil
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
