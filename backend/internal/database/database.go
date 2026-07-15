package database

import (
	"qr-generator/backend/internal/config"
	"qr-generator/backend/internal/models"
	"strings"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect(cfg config.Config, logger *zap.Logger) *gorm.DB {
	dbFields := []zap.Field{
		zap.String("db_host", cfg.DBHost),
		zap.String("db_port", cfg.DBPort),
		zap.String("db_user", cfg.DBUser),
		zap.String("db_name", cfg.DBName),
		zap.Int("db_password_length", len(cfg.DBPassword)),
		zap.Bool("db_password_empty", cfg.DBPassword == ""),
		zap.Bool("db_password_has_outer_whitespace", strings.TrimSpace(cfg.DBPassword) != cfg.DBPassword),
	}
	logger.Info("connecting to database", dbFields...)

	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		logger.Fatal("database connection failed", append(dbFields, zap.Error(err))...)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Plan{},
		&models.Subscription{},
		&models.Payment{},
		&models.Folder{},
		&models.QRCode{},
		&models.QRDesign{},
		&models.QRTemplate{},
		&models.QRScan{},
		&models.SystemLog{},
	); err != nil {
		logger.Fatal("database migration failed", zap.Error(err))
	}

	return db
}
