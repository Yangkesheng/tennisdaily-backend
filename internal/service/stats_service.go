package service

import (
	"time"

	"tennisdaily-backend/internal/model"
	"tennisdaily-backend/internal/repository"
)

type StatsService struct {
	repo *repository.SessionRepository
	loc  *time.Location
}

func NewStatsService(repo *repository.SessionRepository) *StatsService {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.Local
	}
	return &StatsService{repo: repo, loc: loc}
}

func (s *StatsService) Month(userID int64) (model.SessionStats, error) {
	now := time.Now().In(s.loc)
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, s.loc)
	end := start.AddDate(0, 1, 0)
	return s.repo.StatsByMonth(userID, start, end)
}
