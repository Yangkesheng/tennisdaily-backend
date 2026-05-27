package model

import "time"

type User struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	OpenID    string    `json:"openid" gorm:"column:openid;size:128;uniqueIndex;not null"`
	Nickname  string    `json:"nickname" gorm:"size:128;not null;default:''"`
	AvatarURL string    `json:"avatarUrl" gorm:"not null;default:''"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (User) TableName() string {
	return "users"
}
