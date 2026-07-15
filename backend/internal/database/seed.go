package database

import (
	"time"

	"qr-generator/backend/internal/config"
	"qr-generator/backend/internal/models"
	"qr-generator/backend/internal/shared"
	"qr-generator/backend/internal/utils"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB, cfg config.Config) error {
	roles := []models.Role{
		{Name: shared.RoleNameUser, Description: "Default application user"},
		{Name: shared.RoleNameAdmin, Description: "System administrator"},
	}
	for _, role := range roles {
		if err := db.FirstOrCreate(&role, models.Role{Name: role.Name}).Error; err != nil {
			return err
		}
	}

	free := models.Plan{
		Name: shared.PlanNameFree, Price: 0, DurationDays: 3650, MaxQRCodes: cfg.FreeMaxQRCodes,
		AllowDynamicQR: false, AllowLogo: false, AllowAnalytics: false, AllowSVGPDFExport: false,
		Description: "Free plan for basic static QR codes", Status: shared.PlanStatusActive,
	}
	pro := models.Plan{
		Name: shared.PlanNamePro, Price: 99000, DurationDays: 30, MaxQRCodes: 1000,
		AllowDynamicQR: true, AllowLogo: true, AllowAnalytics: true, AllowSVGPDFExport: true,
		Description: "Pro plan with dynamic QR, logo, analytics and advanced exports", Status: shared.PlanStatusActive,
	}
	for _, plan := range []models.Plan{free, pro} {
		if err := db.FirstOrCreate(&plan, models.Plan{Name: plan.Name}).Error; err != nil {
			return err
		}
	}

	templates := []models.QRTemplate{
		{Name: "Classic", ConfigJSON: `{"foreground":"#111827","background":"#FFFFFF"}`, IsPro: false, Status: shared.TemplateStatusActive},
		{Name: "Pro Dark", ConfigJSON: `{"foreground":"#F8FAFC","background":"#111827"}`, IsPro: true, Status: shared.TemplateStatusActive},
	}
	for _, template := range templates {
		if err := db.FirstOrCreate(&template, models.QRTemplate{Name: template.Name}).Error; err != nil {
			return err
		}
	}

	var count int64
	if err := db.Model(&models.User{}).Where("email = ?", cfg.AdminEmail).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		passwordHash, err := utils.HashPassword(cfg.AdminPassword)
		if err != nil {
			return err
		}
		var adminRole models.Role
		var userRole models.Role
		var freePlan models.Plan
		if err := db.Where("name = ?", shared.RoleNameAdmin).First(&adminRole).Error; err != nil {
			return err
		}
		if err := db.Where("name = ?", shared.RoleNameUser).First(&userRole).Error; err != nil {
			return err
		}
		if err := db.Where("name = ?", shared.PlanNameFree).First(&freePlan).Error; err != nil {
			return err
		}
		admin := models.User{
			FullName: "System Admin", Email: cfg.AdminEmail, PasswordHash: passwordHash,
			Status: shared.UserStatusActive, Roles: []models.Role{adminRole, userRole},
		}
		if err := db.Create(&admin).Error; err != nil {
			return err
		}
		sub := models.Subscription{
			UserID: admin.ID, PlanID: freePlan.ID, StartDate: time.Now(),
			EndDate: time.Now().AddDate(10, 0, 0), Status: shared.SubscriptionStatusActive,
		}
		if err := db.Create(&sub).Error; err != nil {
			return err
		}
	}
	return nil
}
