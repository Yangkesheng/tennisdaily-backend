package service

import (
	"strings"
	"tennisdaily-backend/internal/model"
	"tennisdaily-backend/internal/repository"
)

type TennisStringService struct {
	repo *repository.TennisStringRepository
}

func NewTennisStringService(repo *repository.TennisStringRepository) *TennisStringService {
	return &TennisStringService{repo: repo}
}
func normalizeStringQuery(q model.StringLibraryQuery) model.StringLibraryQuery {
	q.Keyword = strings.Join(strings.Fields(q.Keyword), " ")
	q.Brand = strings.TrimSpace(q.Brand)
	if q.Page == 0 {
		q.Page = 1
	}
	if q.PageSize == 0 {
		q.PageSize = 20
	}
	return q
}
func (s *TennisStringService) List(q model.StringLibraryQuery) (model.StringLibraryPage, error) {
	return s.repo.List(normalizeStringQuery(q))
}
func (s *TennisStringService) Brands() ([]string, error) { return s.repo.Brands() }
