package repository

import (
	"errors"

	"tennisdaily-backend/internal/model"

	"gorm.io/gorm"
)

type UserSettingRepository struct {
	db *gorm.DB
}

func NewUserSettingRepository(db *gorm.DB) *UserSettingRepository {
	return &UserSettingRepository{db: db}
}

func (r *UserSettingRepository) ListByUser(userID int64) ([]model.UserSetting, error) {
	var settings []model.UserSetting
	err := r.db.Where("user_id = ?", userID).Find(&settings).Error
	return settings, err
}

func (r *UserSettingRepository) Upsert(userID int64, key int16, value string) error {
	var setting model.UserSetting
	err := r.db.Where("user_id = ? AND setting_key = ?", userID, key).First(&setting).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.Create(&model.UserSetting{
			UserID:       userID,
			SettingKey:   key,
			SettingValue: value,
		}).Error
	}
	if err != nil {
		return err
	}

	setting.SettingValue = value
	return r.db.Save(&setting).Error
}
