package model

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int64          `json:"id" gorm:"primaryKey"`
	OpenID    string         `json:"openid" gorm:"column:openid;size:128;uniqueIndex"`
	Phone     string         `json:"phone" gorm:"size:32;index"`
	Nickname  string         `json:"nickname" gorm:"size:128;not null;default:''"`
	AvatarURL string         `json:"avatarUrl" gorm:"not null;default:''"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (User) TableName() string {
	return "users"
}

type UserWechatIdentity struct {
	ID        int64          `json:"id" gorm:"primaryKey"`
	UserID    int64          `json:"userId" gorm:"not null;index"`
	AppID     string         `json:"appid" gorm:"size:64;not null;uniqueIndex:idx_user_wechat_appid_openid"`
	OpenID    string         `json:"openid" gorm:"size:128;not null;uniqueIndex:idx_user_wechat_appid_openid"`
	UnionID   string         `json:"unionid" gorm:"size:128;index"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (UserWechatIdentity) TableName() string {
	return "user_wechat_identities"
}

type UserResponse struct {
	ID          int64     `json:"id"`
	Phone       string    `json:"phone"`
	MaskedPhone string    `json:"maskedPhone"`
	Nickname    string    `json:"nickname"`
	AvatarURL   string    `json:"avatarUrl"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func NewUserResponse(user User) UserResponse {
	return UserResponse{
		ID:          user.ID,
		Phone:       user.Phone,
		MaskedPhone: MaskPhone(user.Phone),
		Nickname:    user.Nickname,
		AvatarURL:   user.AvatarURL,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

func MaskPhone(phone string) string {
	phone = strings.TrimSpace(phone)
	if len(phone) < 7 {
		return phone
	}
	return phone[:3] + "****" + phone[len(phone)-4:]
}
