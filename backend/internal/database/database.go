package database

import (
	"fmt"

	"qr-generator/backend/internal/config"
	"qr-generator/backend/internal/models"
	"qr-generator/backend/internal/shared"
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
		zap.Bool("database_url_configured", strings.TrimSpace(cfg.DatabaseURL) != ""),
		zap.Int("db_password_length", len(cfg.DBPassword)),
		zap.Bool("db_password_empty", cfg.DBPassword == ""),
		zap.Bool("db_password_has_outer_whitespace", strings.TrimSpace(cfg.DBPassword) != cfg.DBPassword),
	}
	logger.Info("connecting to database", dbFields...)

	dsn, err := cfg.DSN()
	if err != nil {
		logger.Fatal("invalid database configuration", zap.Error(err))
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("database connection failed", append(dbFields, zap.Error(err))...)
	}

	// Migrate the referenced table first. Older deployments may already have
	// payments/subscriptions whose plan rows were removed; MySQL refuses to add
	// the foreign keys until those legacy rows are repaired.
	if err := db.AutoMigrate(&models.User{}, &models.Role{}, &models.Plan{}); err != nil {
		logger.Fatal("database migration failed", zap.Error(err))
	}
	freePlan, err := ensureBasePlans(db, cfg)
	if err != nil {
		logger.Fatal("database plan bootstrap failed", zap.Error(err))
	}
	if err := repairOrphanPlanReferences(db, freePlan.ID); err != nil {
		logger.Fatal("database legacy data repair failed", zap.Error(err))
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

func ensureBasePlans(db *gorm.DB, cfg config.Config) (models.Plan, error) {
	free := models.Plan{
		Name: shared.PlanNameFree, Price: 0, DurationDays: 3650,
		MaxQRCodes: cfg.FreeMaxQRCodes, AllowDynamicQR: false, AllowLogo: false,
		AllowAnalytics: false, AllowSVGPDFExport: false,
		Description: "Free plan for basic static QR codes", Status: shared.PlanStatusActive,
	}
	pro := models.Plan{
		Name: shared.PlanNamePro, Price: 99000, DurationDays: 30,
		MaxQRCodes: 1000, AllowDynamicQR: true, AllowLogo: true,
		AllowAnalytics: true, AllowSVGPDFExport: true,
		Description: "Pro plan with dynamic QR, logo, analytics and advanced exports", Status: shared.PlanStatusActive,
	}
	for _, plan := range []models.Plan{free, pro} {
		if err := db.Where("name = ?", plan.Name).FirstOrCreate(&plan).Error; err != nil {
			return models.Plan{}, err
		}
	}
	var persisted models.Plan
	if err := db.Where("name = ?", shared.PlanNameFree).First(&persisted).Error; err != nil {
		return models.Plan{}, err
	}
	return persisted, nil
}

func repairOrphanPlanReferences(db *gorm.DB, freePlanID uint) error {
	if freePlanID == 0 {
		return fmt.Errorf("free plan has no ID")
	}
	if db.Migrator().HasTable(&models.Subscription{}) {
		if err := db.Model(&models.Subscription{}).
			Where("NOT EXISTS (?)", db.Model(&models.Plan{}).Select("1").Where("plans.id = subscriptions.plan_id")).
			Update("plan_id", freePlanID).Error; err != nil {
			return fmt.Errorf("repair subscriptions with missing plans: %w", err)
		}
	}
	if db.Migrator().HasTable(&models.Payment{}) {
		if err := db.Model(&models.Payment{}).
			Where("NOT EXISTS (?)", db.Model(&models.Plan{}).Select("1").Where("plans.id = payments.plan_id")).
			Update("plan_id", freePlanID).Error; err != nil {
			return fmt.Errorf("repair payments with missing plans: %w", err)
		}
	}
	return nil
}
