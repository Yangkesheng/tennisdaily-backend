package service

import (
	"time"

	"tennisdaily-backend/internal/repository"
)

type AnalyticsInstallationService struct {
	repo *repository.AnalyticsInstallationRepository
}

func NewAnalyticsInstallationService(repo *repository.AnalyticsInstallationRepository) *AnalyticsInstallationService {
	return &AnalyticsInstallationService{repo: repo}
}

func (s *AnalyticsInstallationService) Touch(installationID string) error {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return err
	}
	return s.repo.Touch(installationID, time.Now().In(loc))
}

func (s *AnalyticsInstallationService) Identify(installationID string, userID int64) error {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return err
	}
	return s.repo.Identify(installationID, userID, time.Now().In(loc))
}
