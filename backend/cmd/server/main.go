package main

import (
	"fmt"

	"qr-generator/backend/internal/config"
	"qr-generator/backend/internal/database"
	"qr-generator/backend/internal/modules/payments"
	"qr-generator/backend/internal/router"

	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()
	logger, _ := zap.NewProduction()
	if cfg.AppEnv == "development" {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()

	db := database.Connect(cfg, logger)
	if err := database.Seed(db, cfg); err != nil {
		logger.Fatal("seed failed", zap.Error(err))
	}

	payments.StartCleanupScheduler(db, logger)

	app := router.Setup(db, cfg, logger)
	if err := app.Run(fmt.Sprintf(":%s", cfg.AppPort)); err != nil {
		logger.Fatal("server failed", zap.Error(err))
	}
}
 