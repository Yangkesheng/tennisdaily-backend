package repository

import (
	"errors"
	"time"

	"tennisdaily-backend/internal/model"

	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) List(userID int64) ([]model.TennisSession, error) {
	var sessions []model.TennisSession
	err := r.db.Where("user_id = ?", userID).Order("date DESC, created_at DESC").Find(&sessions).Error
	return sessions, err
}

func (r *SessionRepository) ListByDate(userID int64, date time.Time) ([]model.TennisSession, error) {
	var sessions []model.TennisSession
	err := r.db.Where("user_id = ? AND date = ?", userID, date).Order("created_at DESC").Find(&sessions).Error
	return sessions, err
}

func (r *SessionRepository) Create(session *model.TennisSession) error {
	return r.db.Create(session).Error
}

func (r *SessionRepository) FindByID(userID, id int64) (*model.TennisSession, error) {
	var session model.TennisSession
	err := r.db.Where("user_id = ? AND id = ?", userID, id).First(&session).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) Update(session *model.TennisSession) error {
	return r.db.Save(session).Error
}

func (r *SessionRepository) SoftDelete(userID, id int64) (bool, error) {
	result := r.db.Model(&model.TennisSession{}).
		Where("user_id = ? AND id = ? AND deleted_at IS NULL", userID, id).
		Updates(map[string]interface{}{"deleted_at": time.Now(), "updated_at": time.Now()})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *SessionRepository) Latest(userID int64) (*model.TennisSession, error) {
	var session model.TennisSession
	err := r.db.Where("user_id = ?", userID).Order("date DESC, created_at DESC").First(&session).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) CalendarDays(userID int64, start, end time.Time) ([]model.SessionCalendarDay, error) {
	var days []model.SessionCalendarDay
	err := r.db.Model(&model.TennisSession{}).
		Select("DATE_FORMAT(date, '%Y-%m-%d') AS date, COUNT(*) AS count").
		Where("user_id = ? AND date >= ? AND date < ?", userID, start, end).
		Group("date").
		Order("date ASC").
		Scan(&days).Error
	return days, err
}

func (r *SessionRepository) ListLatest(userID int64, limit int) ([]model.TennisSession, error) {
	var sessions []model.TennisSession
	err := r.db.Where("user_id = ?", userID).Order("date DESC, created_at DESC").Limit(limit).Find(&sessions).Error
	return sessions, err
}

func (r *SessionRepository) StatsByMonth(userID int64, start, end time.Time) (model.SessionStats, error) {
	var stats model.SessionStats

	if err := r.db.Model(&model.TennisSession{}).
		Select("COUNT(*) AS month_count, COALESCE(SUM(duration_minutes), 0) AS month_minutes, COALESCE(SUM(cost), 0) AS month_cost").
		Where("user_id = ? AND date >= ? AND date < ?", userID, start, end).
		Scan(&stats).Error; err != nil {
		return stats, err
	}

	var total int64
	if err := r.db.Model(&model.TennisSession{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return stats, err
	}
	stats.TotalCount = int(total)
	return stats, nil
}
