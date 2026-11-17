package service

import (
	"time"

	"tennisdaily-backend/internal/model"
	"tennisdaily-backend/internal/repository"
)

type AdminAnalyticsService struct {
	repo *repository.AdminAnalyticsRepository
}

func NewAdminAnalyticsService(repo *repository.AdminAnalyticsRepository) *AdminAnalyticsService {
	return &AdminAnalyticsService{repo: repo}
}

func (s *AdminAnalyticsService) Daily(startDate, endDate string) (*model.AdminAnalyticsResponse, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return nil, err
	}
	start, err := time.ParseInLocation("2006-01-02", startDate, loc)
	if err != nil {
		return nil, ErrInvalidRequest
	}
	endInclusive, err := time.ParseInLocation("2006-01-02", endDate, loc)
	if err != nil || endInclusive.Before(start) {
		return nil, ErrInvalidRequest
	}
	if endInclusive.Sub(start) > 365*24*time.Hour {
		return nil, ErrInvalidRequest
	}
	end := endInclusive.AddDate(0, 0, 1)
	days, unregisteredOpeners, err := s.repo.Daily(start, end)
	if err != nil {
		return nil, err
	}
	return &model.AdminAnalyticsResponse{
		StartDate: startDate, EndDate: endDate, TimeZone: "Asia/Shanghai",
		GeneratedAt: time.Now().In(loc).Format("2006-01-02 15:04:05"), UnregisteredOpeners: unregisteredOpeners, Days: days,
	}, nil
}
