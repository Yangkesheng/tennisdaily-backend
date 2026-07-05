package service

import (
	"time"

	"tennisdaily-backend/internal/model"
	"tennisdaily-backend/internal/repository"
)

const homeRatingTrendLimit = 5

type HomeService struct {
	sessionRepo             *repository.SessionRepository
	racketRepo              *repository.RacketRepository
	sessionCategoryResolver model.SessionCategoryTextResolver
	loc                     *time.Location
}

func NewHomeService(sessionRepo *repository.SessionRepository, racketRepo *repository.RacketRepository, sessionCategoryResolver model.SessionCategoryTextResolver) *HomeService {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.Local
	}
	return &HomeService{sessionRepo: sessionRepo, racketRepo: racketRepo, sessionCategoryResolver: sessionCategoryResolver, loc: loc}
}

func (s *HomeService) Summary(userID int64, query model.HomeSummaryQuery) (model.HomeSummaryResponse, error) {
	now := time.Now().In(s.loc)
	year := query.Year
	month := query.Month

	if year == 0 && month == 0 {
		year = now.Year()
		month = int(now.Month())
	} else if year == 0 || month == 0 {
		return model.HomeSummaryResponse{}, ErrInvalidRequest
	}

	if year < 2000 || year > 2100 || month < 1 || month > 12 {
		return model.HomeSummaryResponse{}, ErrInvalidRequest
	}

	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, s.loc)
	end := start.AddDate(0, 1, 0)
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, s.loc)
	yearEnd := yearStart.AddDate(1, 0, 0)

	sessionStats, err := s.sessionRepo.StatsByMonth(userID, start, end)
	if err != nil {
		return model.HomeSummaryResponse{}, err
	}

	racketCost, err := s.racketRepo.SumPurchaseCostByMonth(userID, start, end)
	if err != nil {
		return model.HomeSummaryResponse{}, err
	}

	stringingCost, err := s.racketRepo.SumStringingCostByMonth(userID, start, end)
	if err != nil {
		return model.HomeSummaryResponse{}, err
	}

	yearCount, err := s.sessionRepo.CountByDateRange(userID, yearStart, yearEnd)
	if err != nil {
		return model.HomeSummaryResponse{}, err
	}

	latestSession, err := s.sessionRepo.Latest(userID)
	if err != nil {
		return model.HomeSummaryResponse{}, err
	}

	trendSessions, err := s.sessionRepo.ListLatest(userID, homeRatingTrendLimit)
	if err != nil {
		return model.HomeSummaryResponse{}, err
	}

	var latestSessionResponse *model.SessionResponse
	if latestSession != nil {
		resp := model.NewSessionResponse(*latestSession, s.sessionCategoryResolver)
		latestSessionResponse = &resp
	}

	ratingTrend := make([]model.HomeRatingTrendItemResponse, 0, len(trendSessions))
	for _, session := range trendSessions {
		ratingTrend = append(ratingTrend, model.HomeRatingTrendItemResponse{
			ID:     session.ID,
			Date:   session.Date.Format("2006-01-02"),
			Rating: session.Rating,
		})
	}

	sessionCost := sessionStats.MonthCost

	return model.HomeSummaryResponse{
		Year:  year,
		Month: month,
		Session: model.HomeSessionSummaryResponse{
			MonthCount:   sessionStats.MonthCount,
			MonthMinutes: sessionStats.MonthMinutes,
			MonthCost:    sessionCost,
			YearCount:    yearCount,
			TotalCount:   sessionStats.TotalCount,
		},
		Expense: model.HomeExpenseSummaryResponse{
			SessionCost:   sessionCost,
			RacketCost:    racketCost,
			StringingCost: stringingCost,
			TotalCost:     sessionCost + racketCost + stringingCost,
		},
		LatestSession: latestSessionResponse,
		RatingTrend:   ratingTrend,
	}, nil
}
