package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	FeedbackStatusPending   = 0
	FeedbackStatusProcessed = 1
)

// Feedback 用户意见反馈。
type Feedback struct {
	ID        int64          `json:"id" gorm:"primaryKey"`
	UserID    int64          `json:"userId" gorm:"not null;index"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	Contact   string         `json:"contact" gorm:"size:100;not null;default:''"`
	Status    int            `json:"status" gorm:"not null;default:0"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Feedback) TableName() string {
	return "feedback"
}

type CreateFeedbackRequest struct {
	Content string `json:"content" binding:"required"`
	Contact string `json:"contact"`
}

type FeedbackResponse struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Contact   string    `json:"contact"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewFeedbackResponse(feedback Feedback) FeedbackResponse {
	return FeedbackResponse{
		ID:        feedback.ID,
		Content:   feedback.Content,
		Contact:   feedback.Contact,
		Status:    feedback.Status,
		CreatedAt: feedback.CreatedAt,
	}
}
