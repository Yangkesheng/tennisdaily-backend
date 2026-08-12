package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	UserSettingSignature              int16 = 1
	UserSettingDefaultCourtName       int16 = 2
	UserSettingDefaultDurationMinutes int16 = 3
)

type UserSetting struct {
	ID           int64          `json:"id" gorm:"primaryKey"`
	UserID       int64          `json:"userId" gorm:"not null;index"`
	SettingKey   int16          `json:"settingKey" gorm:"not null"`
	SettingValue string         `json:"settingValue" gorm:"size:255;not null;default:''"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (UserSetting) TableName() string {
	return "user_settings"
}

type UserSettingsResponse struct {
	Signature              string `json:"signature"`
	DefaultCourtName       string `json:"defaultCourtName"`
	DefaultDurationMinutes int    `json:"defaultDurationMinutes"`
}

type UpdateUserSettingsRequest struct {
	Signature              *string `json:"signature"`
	DefaultCourtName       *string `json:"defaultCourtName"`
	DefaultDurationMinutes *int    `json:"defaultDurationMinutes"`
}
