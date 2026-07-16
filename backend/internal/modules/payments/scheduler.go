package payments

import (
	"time"

	"qr-generator/backend/internal/models"
	"qr-generator/backend/internal/shared"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// StartCleanupScheduler runs a periodic job every 10 minutes to mark expired payments as CANCELLED.
func StartCleanupScheduler(db *gorm.DB, logger *zap.Logger) {
	// Run immediately on startup to clean up any payments that expired while server was offline
	cleanupExpiredPayments(db, logger)

	ticker := time.NewTicker(10 * time.Minute)
	go func() {
		for range ticker.C {
			cleanupExpiredPayments(db, logger)
		}
	}()
}

func cleanupExpiredPayments(db *gorm.DB, logger *zap.Logger) {
	threshold := time.Now().Add(-10 * time.Minute)
	result := db.Model(&models.Payment{}).
		Where("status = ? AND created_at < ?", shared.PaymentStatusPending, threshold).
		Update("status", shared.PaymentStatusCancelled)
	if result.Error != nil {
		logger.Error("failed to cleanup expired payments", zap.Error(result.Error))
	} else if result.RowsAffected > 0 {
		logger.Info("cleaned up expired payments", zap.Int64("count", result.RowsAffected))
	}
}
