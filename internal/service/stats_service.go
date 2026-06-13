package service

import (
	"fmt"
	"math"
	"time"

	"tennisdaily-backend/internal/model"
	"tennisdaily-backend/internal/repository"
)

const statsRatingTrendLimit = 8

type StatsService struct {
	sessionRepo *repository.SessionRepository
	racketRepo  *repository.RacketRepository
	loc         *time.Location
}

func NewStatsService(sessionRepo *repository.SessionRepository, racketRepo *repository.RacketRepository) *StatsService {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.Local
	}
	return &StatsService{sessionRepo: sessionRepo, racketRepo: racketRepo, loc: loc}
}

func (s *StatsService) Month(userID int64) (model.SessionStats, error) {
	now := time.Now().In(s.loc)
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, s.loc)
	end := start.AddDate(0, 1, 0)
	return s.sessionRepo.StatsByMonth(userID, start, end)
}

func (s *StatsService) Charts(userID int64, query model.StatsChartsQuery) (model.StatsChartsResultResponse, error) {
	if query.Period == model.StatsPeriodMonth {
		return s.monthCharts(userID, query.Year, query.Month)
	}
	if query.Period == model.StatsPeriodYear {
		return s.yearCharts(userID, query.Year)
	}
	return model.StatsChartsResultResponse{}, ErrInvalidRequest
}

func (s *StatsService) monthCharts(userID int64, year, month int) (model.StatsChartsResultResponse, error) {
	if year < 2000 || year > 2100 || month < 1 || month > 12 {
		return model.StatsChartsResultResponse{}, ErrInvalidRequest
	}

	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, s.loc)
	end := start.AddDate(0, 1, 0)

	summary, err := s.buildSummary(userID, start, end)
	if err != nil {
		return model.StatsChartsResultResponse{}, err
	}

	days, err := s.sessionRepo.CalendarDays(userID, start, end)
	if err != nil {
		return model.StatsChartsResultResponse{}, err
	}

	ratingTrend, err := s.monthRatingTrend(userID, start, end)
	if err != nil {
		return model.StatsChartsResultResponse{}, err
	}

	return model.StatsChartsResultResponse{
		Period:    model.StatsPeriodMonth,
		Year:      year,
		Month:     month,
		RangeText: fmt.Sprintf("%d年%d月", year, month),
		Summary:   summary,
		Charts: model.StatsChartsResponse{
			Frequency:            s.weeklyFrequency(start, end, days),
			RatingTrend:          ratingTrend,
			ExpenseBreakdown:     expenseBreakdown(summary.SessionCost, summary.RacketCost, summary.StringingCost),
			SessionTypeBreakdown: sessionTypeBreakdown(summary),
		},
	}, nil
}

func (s *StatsService) yearCharts(userID int64, year int) (model.StatsChartsResultResponse, error) {
	if year < 2000 || year > 2100 {
		return model.StatsChartsResultResponse{}, ErrInvalidRequest
	}

	start := time.Date(year, 1, 1, 0, 0, 0, 0, s.loc)
	end := start.AddDate(1, 0, 0)

	summary, err := s.buildSummary(userID, start, end)
	if err != nil {
		return model.StatsChartsResultResponse{}, err
	}

	frequency, err := s.yearFrequency(userID, start, end)
	if err != nil {
		return model.StatsChartsResultResponse{}, err
	}

	ratingTrend, err := s.yearRatingTrend(userID, start, end)
	if err != nil {
		return model.StatsChartsResultResponse{}, err
	}

	return model.StatsChartsResultResponse{
		Period:    model.StatsPeriodYear,
		Year:      year,
		Month:     0,
		RangeText: fmt.Sprintf("%d年", year),
		Summary:   summary,
		Charts: model.StatsChartsResponse{
			Frequency:            frequency,
			RatingTrend:          ratingTrend,
			ExpenseBreakdown:     expenseBreakdown(summary.SessionCost, summary.RacketCost, summary.StringingCost),
			SessionTypeBreakdown: sessionTypeBreakdown(summary),
		},
	}, nil
}

func (s *StatsService) buildSummary(userID int64, start, end time.Time) (model.StatsChartsSummaryResponse, error) {
	sessionStats, err := s.sessionRepo.StatsAggregate(userID, start, end)
	if err != nil {
		return model.StatsChartsSummaryResponse{}, err
	}

	racketCost, err := s.racketRepo.SumPurchaseCostByRange(userID, start, end)
	if err != nil {
		return model.StatsChartsSummaryResponse{}, err
	}

	stringingCost, err := s.racketRepo.SumStringingCostByRange(userID, start, end)
	if err != nil {
		return model.StatsChartsSummaryResponse{}, err
	}

	return model.StatsChartsSummaryResponse{
		SessionCount:   sessionStats.SessionCount,
		ActiveDayCount: sessionStats.ActiveDayCount,
		TotalMinutes:   sessionStats.TotalMinutes,
		AverageMinutes: round1(sessionStats.AverageMinutes),
		AverageRating:  round1(sessionStats.AverageRating),
		SessionCost:    sessionStats.SessionCost,
		RacketCost:     racketCost,
		StringingCost:  stringingCost,
		TotalCost:      sessionStats.SessionCost + racketCost + stringingCost,
		TrainingCount:  sessionStats.TrainingCount,
		SinglesCount:   sessionStats.SinglesCount,
		DoublesCount:   sessionStats.DoublesCount,
		MatchCount:     sessionStats.MatchCount,
	}, nil
}

func (s *StatsService) weeklyFrequency(start, end time.Time, days []model.SessionCalendarDay) []model.StatsFrequencyChartItemResponse {
	mondayStart := start.AddDate(0, 0, -weekdayOffset(start))
	weeks := 0
	for weekStart := mondayStart; weekStart.Before(end); weekStart = weekStart.AddDate(0, 0, 7) {
		weeks++
	}

	frequency := make([]model.StatsFrequencyChartItemResponse, weeks)
	for index := range frequency {
		frequency[index] = model.StatsFrequencyChartItemResponse{
			Label: fmt.Sprintf("第%d周", index+1),
			Value: 0,
		}
	}

	for _, day := range days {
		date, err := time.ParseInLocation("2006-01-02", day.Date, s.loc)
		if err != nil {
			continue
		}
		index := int(date.Sub(mondayStart).Hours() / 24 / 7)
		if index >= 0 && index < len(frequency) {
			frequency[index].Value += int64(day.Count)
		}
	}

	return frequency
}

func (s *StatsService) monthRatingTrend(userID int64, start, end time.Time) ([]model.StatsRatingTrendItemResponse, error) {
	sessions, err := s.sessionRepo.ListRatingTrendByRange(userID, start, end, statsRatingTrendLimit)
	if err != nil {
		return nil, err
	}

	trend := make([]model.StatsRatingTrendItemResponse, 0, len(sessions))
	for index := len(sessions) - 1; index >= 0; index-- {
		session := sessions[index]
		trend = append(trend, model.StatsRatingTrendItemResponse{
			Label:  fmt.Sprintf("%d/%d", int(session.Date.Month()), session.Date.Day()),
			Date:   session.Date.Format("2006-01-02"),
			Rating: float64(session.Rating),
		})
	}
	return trend, nil
}

func (s *StatsService) yearFrequency(userID int64, start, end time.Time) ([]model.StatsFrequencyChartItemResponse, error) {
	rows, err := s.sessionRepo.MonthlyFrequency(userID, start, end)
	if err != nil {
		return nil, err
	}

	frequency := make([]model.StatsFrequencyChartItemResponse, 12)
	for month := 1; month <= 12; month++ {
		frequency[month-1] = model.StatsFrequencyChartItemResponse{Label: fmt.Sprintf("%d月", month), Value: 0}
	}
	for _, row := range rows {
		if row.Month >= 1 && row.Month <= 12 {
			frequency[row.Month-1].Value = row.Count
		}
	}
	return frequency, nil
}

func (s *StatsService) yearRatingTrend(userID int64, start, end time.Time) ([]model.StatsRatingTrendItemResponse, error) {
	rows, err := s.sessionRepo.MonthlyRatingTrend(userID, start, end)
	if err != nil {
		return nil, err
	}

	trend := make([]model.StatsRatingTrendItemResponse, 0, len(rows))
	for _, row := range rows {
		if row.Month < 1 || row.Month > 12 {
			continue
		}
		trend = append(trend, model.StatsRatingTrendItemResponse{
			Label:  fmt.Sprintf("%d月", row.Month),
			Date:   fmt.Sprintf("%04d-%02d", start.Year(), row.Month),
			Rating: round1(row.Rating),
		})
	}
	return trend, nil
}

func expenseBreakdown(sessionCost, racketCost, stringingCost float64) []model.StatsBreakdownItemResponse {
	totalCost := sessionCost + racketCost + stringingCost
	return []model.StatsBreakdownItemResponse{
		{Key: "session", Label: "打球", Value: sessionCost, Percent: percent(sessionCost, totalCost)},
		{Key: "racket", Label: "球拍", Value: racketCost, Percent: percent(racketCost, totalCost)},
		{Key: "stringing", Label: "穿线", Value: stringingCost, Percent: percent(stringingCost, totalCost)},
	}
}

func sessionTypeBreakdown(summary model.StatsChartsSummaryResponse) []model.StatsBreakdownItemResponse {
	total := float64(summary.SessionCount)
	return []model.StatsBreakdownItemResponse{
		{Key: "training", Label: "训练", Value: float64(summary.TrainingCount), Percent: percent(float64(summary.TrainingCount), total)},
		{Key: "singles", Label: "单打", Value: float64(summary.SinglesCount), Percent: percent(float64(summary.SinglesCount), total)},
		{Key: "doubles", Label: "双打", Value: float64(summary.DoublesCount), Percent: percent(float64(summary.DoublesCount), total)},
		{Key: "match", Label: "比赛", Value: float64(summary.MatchCount), Percent: percent(float64(summary.MatchCount), total)},
	}
}

func weekdayOffset(date time.Time) int {
	weekday := int(date.Weekday())
	if weekday == 0 {
		return 6
	}
	return weekday - 1
}

func percent(value, total float64) float64 {
	if total <= 0 {
		return 0
	}
	return round1(value / total * 100)
}

func round1(value float64) float64 {
	return math.Round(value*10) / 10
}
