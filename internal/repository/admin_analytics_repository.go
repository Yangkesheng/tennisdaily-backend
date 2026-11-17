package repository

import (
	"time"

	"tennisdaily-backend/internal/model"

	"gorm.io/gorm"
)

type AdminAnalyticsRepository struct{ db *gorm.DB }

func NewAdminAnalyticsRepository(db *gorm.DB) *AdminAnalyticsRepository {
	return &AdminAnalyticsRepository{db: db}
}

type dailyCount struct {
	Date  string `gorm:"column:date"`
	Count int64  `gorm:"column:count"`
}

func (r *AdminAnalyticsRepository) dailyCounts(table string, start, end time.Time) ([]dailyCount, error) {
	var rows []dailyCount
	err := r.db.Table(table).
		Select("DATE_FORMAT(created_at, '%Y-%m-%d') AS date, COUNT(*) AS count").
		Where("created_at >= ? AND created_at < ? AND deleted_at IS NULL", start, end).
		Group("DATE_FORMAT(created_at, '%Y-%m-%d')").
		Scan(&rows).Error
	return rows, err
}

func (r *AdminAnalyticsRepository) TotalUsers() (int64, error) {
	var count int64
	err := r.db.Table("users").Where("deleted_at IS NULL").Count(&count).Error
	return count, err
}

func (r *AdminAnalyticsRepository) Daily(start, end time.Time) ([]model.AdminAnalyticsDay, int64, error) {
	days := make([]model.AdminAnalyticsDay, 0, int(end.Sub(start).Hours()/24))
	index := make(map[string]int, cap(days))
	for day := start; day.Before(end); day = day.AddDate(0, 0, 1) {
		key := day.Format("2006-01-02")
		index[key] = len(days)
		days = append(days, model.AdminAnalyticsDay{Date: key})
	}
	for _, target := range []struct {
		table string
		set   func(*model.AdminAnalyticsDay, int64)
	}{
		{"users", func(d *model.AdminAnalyticsDay, n int64) { d.NewUsers = n }},
		{"tennis_sessions", func(d *model.AdminAnalyticsDay, n int64) { d.NewSessions = n }},
		{"racket", func(d *model.AdminAnalyticsDay, n int64) { d.NewRackets = n }},
		{"racket_stringing_record", func(d *model.AdminAnalyticsDay, n int64) { d.NewStringingRecords = n }},
		{"my_shoes", func(d *model.AdminAnalyticsDay, n int64) { d.NewShoes = n }},
	} {
		rows, err := r.dailyCounts(target.table, start, end)
		if err != nil {
			return nil, 0, err
		}
		for _, row := range rows {
			if i, ok := index[row.Date]; ok {
				target.set(&days[i], row.Count)
			}
		}
	}
	var activeRows []dailyCount
	err := r.db.Raw(`
		SELECT DATE_FORMAT(activity.created_at, '%Y-%m-%d') AS date, COUNT(DISTINCT activity.user_id) AS count
		FROM (
			SELECT user_id, created_at FROM tennis_sessions WHERE created_at >= ? AND created_at < ? AND deleted_at IS NULL
			UNION ALL
			SELECT user_id, created_at FROM racket WHERE created_at >= ? AND created_at < ? AND deleted_at IS NULL
			UNION ALL
			SELECT user_id, created_at FROM racket_stringing_record WHERE created_at >= ? AND created_at < ? AND deleted_at IS NULL
			UNION ALL
			SELECT user_id, created_at FROM my_shoes WHERE created_at >= ? AND created_at < ? AND deleted_at IS NULL
		) AS activity
		GROUP BY DATE_FORMAT(activity.created_at, '%Y-%m-%d')`, start, end, start, end, start, end, start, end).Scan(&activeRows).Error
	if err != nil {
		return nil, 0, err
	}
	for _, row := range activeRows {
		if i, ok := index[row.Date]; ok {
			days[i].ActiveUsers = row.Count
		}
	}
	var anonymousRows []dailyCount
	err = r.db.Table("analytics_installation_days AS activity_day").
		Select("DATE_FORMAT(activity_day.activity_date, '%Y-%m-%d') AS date, COUNT(DISTINCT activity_day.installation_id) AS count").
		Joins("JOIN analytics_installations AS installation ON installation.installation_id = activity_day.installation_id").
		Where("activity_day.activity_date >= ? AND activity_day.activity_date < ? AND installation.user_id IS NULL", start, end).
		Group("DATE_FORMAT(activity_day.activity_date, '%Y-%m-%d')").Scan(&anonymousRows).Error
	if err != nil {
		return nil, 0, err
	}
	for _, row := range anonymousRows {
		if i, ok := index[row.Date]; ok {
			days[i].UnregisteredOpeners = row.Count
		}
	}
	var anonymousTotal int64
	err = r.db.Table("analytics_installation_days AS activity_day").
		Joins("JOIN analytics_installations AS installation ON installation.installation_id = activity_day.installation_id").
		Where("activity_day.activity_date >= ? AND activity_day.activity_date < ? AND installation.user_id IS NULL", start, end).
		Distinct("activity_day.installation_id").Count(&anonymousTotal).Error
	if err != nil {
		return nil, 0, err
	}
	return days, anonymousTotal, nil
}
