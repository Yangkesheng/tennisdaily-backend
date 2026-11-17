package repository

import (
	"time"

	"tennisdaily-backend/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AnalyticsInstallationRepository struct{ db *gorm.DB }

func NewAnalyticsInstallationRepository(db *gorm.DB) *AnalyticsInstallationRepository {
	return &AnalyticsInstallationRepository{db: db}
}

func (r *AnalyticsInstallationRepository) Touch(installationID string, now time.Time) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		row := model.AnalyticsInstallation{InstallationID: installationID, FirstSeenAt: now, LastSeenAt: now}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "installation_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"last_seen_at": now}),
		}).Create(&row).Error; err != nil {
			return err
		}
		day := model.AnalyticsInstallationDay{InstallationID: installationID, ActivityDate: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())}
		return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&day).Error
	})
}

func (r *AnalyticsInstallationRepository) Identify(installationID string, userID int64, now time.Time) error {
	row := model.AnalyticsInstallation{InstallationID: installationID, UserID: &userID, FirstSeenAt: now, LastSeenAt: now}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "installation_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"user_id": userID, "last_seen_at": now}),
	}).Create(&row).Error
}
