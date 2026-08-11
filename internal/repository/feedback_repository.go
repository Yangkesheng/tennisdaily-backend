package repository

import (
	"tennisdaily-backend/internal/model"

	"gorm.io/gorm"
)

type FeedbackRepository struct {
	db *gorm.DB
}

func NewFeedbackRepository(db *gorm.DB) *FeedbackRepository {
	return &FeedbackRepository{db: db}
}

func (r *FeedbackRepository) Create(feedback *model.Feedback) error {
	return r.db.Create(feedback).Error
}
