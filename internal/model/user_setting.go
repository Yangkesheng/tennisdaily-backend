package model

import (
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

const (
	UserSettingSignature              int16 = 1
	UserSettingDefaultCourtName       int16 = 2
	UserSettingDefaultDurationMinutes int16 = 3
	UserSettingDefaultCourtLocation   int16 = 4
)

type UserSetting struct {
	ID           int64          `json:"id" gorm:"primaryKey"`
	UserID       int64          `json:"userId" gorm:"not null;index"`
	SettingKey   int16          `json:"settingKey" gorm:"not null"`
	SettingValue string         `json:"settingValue" gorm:"type:text;not null"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (UserSetting) TableName() string {
	return "user_settings"
}

type UserSettingsResponse struct {
	DefaultCourtLocation   *CourtLocation `json:"defaultCourtLocation"`
	Signature              string         `json:"signature"`
	DefaultCourtName       string         `json:"defaultCourtName"`
	DefaultDurationMinutes int            `json:"defaultDurationMinutes"`
}

type UpdateUserSettingsRequest struct {
	DefaultCourtLocation   LocationField `json:"defaultCourtLocation"`
	Signature              *string       `json:"signature"`
	DefaultCourtName       *string       `json:"defaultCourtName"`
	DefaultDurationMinutes *int          `json:"defaultDurationMinutes"`
}

func (r *UpdateUserSettingsRequest) UnmarshalJSON(data []byte) error {
	type plain UpdateUserSettingsRequest
	var value plain
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if _, ok := fields["defaultCourtName"]; ok && value.DefaultCourtName == nil {
		return errors.New("invalid court name")
	}
	*r = UpdateUserSettingsRequest(value)
	return nil
}
