package model

import "time"

type AnalyticsInstallation struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	InstallationID string    `json:"-" gorm:"column:installation_id;size:64;not null;uniqueIndex"`
	UserID         *int64    `json:"-" gorm:"column:user_id;index"`
	FirstSeenAt    time.Time `json:"-" gorm:"column:first_seen_at;not null"`
	LastSeenAt     time.Time `json:"-" gorm:"column:last_seen_at;not null;index"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

func (AnalyticsInstallation) TableName() string { return "analytics_installations" }

type AnalyticsInstallationDay struct {
	ID             int64     `json:"-" gorm:"primaryKey"`
	InstallationID string    `json:"-" gorm:"column:installation_id;size:64;not null;uniqueIndex:uk_analytics_installation_day"`
	ActivityDate   time.Time `json:"-" gorm:"column:activity_date;type:date;not null;uniqueIndex:uk_analytics_installation_day"`
	CreatedAt      time.Time `json:"-"`
}

func (AnalyticsInstallationDay) TableName() string { return "analytics_installation_days" }

type AppOpenRequest struct {
	InstallationID string `json:"installationId" binding:"required,min=32,max=64,alphanum"`
}
