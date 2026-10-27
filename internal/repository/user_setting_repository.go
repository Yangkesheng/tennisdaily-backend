package repository

import (
	"context"
	"errors"

	"tennisdaily-backend/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserSettingRepository struct{ db *gorm.DB }

func NewUserSettingRepository(db *gorm.DB) *UserSettingRepository {
	return &UserSettingRepository{db: db}
}
func (r *UserSettingRepository) ListByUser(ctx context.Context, userID int64) ([]model.UserSetting, error) {
	var rows []model.UserSetting
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&rows).Error
	return rows, err
}

// Apply serializes all setting writers, including first-time insertion of absent keys.
func (r *UserSettingRepository) Apply(ctx context.Context, userID int64, updates []model.UserSetting, validate func([]model.UserSetting) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user model.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", userID).First(&user).Error; err != nil {
			return err
		}
		var rows []model.UserSetting
		if err := tx.Where("user_id = ?", userID).Find(&rows).Error; err != nil {
			return err
		}
		byKey := make(map[int16]int)
		for i, row := range rows {
			if _, ok := byKey[row.SettingKey]; ok {
				return errors.New("duplicate active user setting")
			}
			byKey[row.SettingKey] = i
		}
		for _, update := range updates {
			if i, ok := byKey[update.SettingKey]; ok {
				rows[i].SettingValue = update.SettingValue
				if err := tx.Save(&rows[i]).Error; err != nil {
					return err
				}
			} else {
				update.UserID = userID
				if err := tx.Create(&update).Error; err != nil {
					return err
				}
				byKey[update.SettingKey] = len(rows)
				rows = append(rows, update)
			}
		}
		return validate(rows)
	})
}
