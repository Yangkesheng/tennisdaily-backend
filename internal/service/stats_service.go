package service

import (
	"fmt"
	"math"
	"time"

	"tennisdaily-backend/internal/config"
	"tennisdaily-backend/internal/model"
	"tennisdaily-backend/internal/repository"
)

const statsRatingTrendLimit = 8

type StatsService struct {
	sessionRepo             *repository.SessionRepository
	racketRepo              *repository.RacketRepository
	shoeRepo                *repository.ShoeRepository
	sessionCategoryResolver *config.SessionCategoryResolver
	loc                     *time.Location
}

func NewStatsService(sessionRepo *repository.SessionRepository, racketRepo *repository.RacketRepository, shoeRepo *repository.ShoeRepository, sessionCategoryResolver *config.SessionCategoryResolver) *StatsService {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.Local
	}
	return &StatsService{sessionRepo: sessionRepo, racketRepo: racketRepo, shoeRepo: shoeRepo, sessionCategoryResolver: sessionCategoryResolver, loc: loc}
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

	breakdownRows, err := s.sessionRepo.SessionBreakdownByCategory(userID, start, end)
	if err != nil {
		return model.StatsChartsResultResponse{}, err
	}

	return model.StatsChartsResultResponse{
		Period:    model.StatsPeriodMonth,
		Year:      year,
		Month:     month,
		RangeText: fmt.Sprintf("%d年%d月", year, month),
		Summary:   summary,
		Charts:    s.buildCharts(start, end, days, ratingTrend, breakdownRows, summary),
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

	breakdownRows, err := s.sessionRepo.SessionBreakdownByCategory(userID, start, end)
	if err != nil {
		return model.StatsChartsResultResponse{}, err
	}

	return model.StatsChartsResultResponse{
		Period:    model.StatsPeriodYear,
		Year:      year,
		Month:     0,
		RangeText: fmt.Sprintf("%d年", year),
		Summary:   summary,
		Charts:    s.buildChartsForFrequency(frequency, ratingTrend, breakdownRows, summary),
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

	shoeCost, err := s.shoeRepo.SumPurchaseCostByRange(userID, start, end)
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
		ShoeCost:       shoeCost,
		TotalCost:      sessionStats.SessionCost + racketCost + stringingCost + shoeCost,
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

func expenseBreakdown(sessionCost, racketCost, stringingCost, shoeCost float64) []model.StatsBreakdownItemResponse {
	totalCost := sessionCost + racketCost + stringingCost + shoeCost
	return []model.StatsBreakdownItemResponse{
		{Key: "session", Label: "打球", Value: sessionCost, Percent: percent(sessionCost, totalCost)},
		{Key: "racket", Label: "球拍", Value: racketCost, Percent: percent(racketCost, totalCost)},
		{Key: "stringing", Label: "穿线", Value: stringingCost, Percent: percent(stringingCost, totalCost)},
		{Key: "shoe", Label: "球鞋", Value: shoeCost, Percent: percent(shoeCost, totalCost)},
	}
}

func (s *StatsService) buildCharts(start, end time.Time, days []model.SessionCalendarDay, ratingTrend []model.StatsRatingTrendItemResponse, rows []model.StatsSessionBreakdownRow, summary model.StatsChartsSummaryResponse) model.StatsChartsResponse {
	return s.buildChartsForFrequency(s.weeklyFrequency(start, end, days), ratingTrend, rows, summary)
}

func (s *StatsService) buildChartsForFrequency(frequency []model.StatsFrequencyChartItemResponse, ratingTrend []model.StatsRatingTrendItemResponse, rows []model.StatsSessionBreakdownRow, summary model.StatsChartsSummaryResponse) model.StatsChartsResponse {
	return model.StatsChartsResponse{
		Frequency:                           frequency,
		RatingTrend:                         ratingTrend,
		ExpenseBreakdown:                    expenseBreakdown(summary.SessionCost, summary.RacketCost, summary.StringingCost, summary.ShoeCost),
		SessionCategoryCountBreakdown:       s.sessionCategoryBreakdown(rows, float64(summary.SessionCount), sessionBreakdownMetricCount),
		SessionSubCategoryCountBreakdown:    s.sessionSubCategoryBreakdown(rows, float64(summary.SessionCount), sessionBreakdownMetricCount),
		SessionCategoryDurationBreakdown:    s.sessionCategoryBreakdown(rows, float64(summary.TotalMinutes), sessionBreakdownMetricMinutes),
		SessionSubCategoryDurationBreakdown: s.sessionSubCategoryBreakdown(rows, float64(summary.TotalMinutes), sessionBreakdownMetricMinutes),
		SessionCategoryCostBreakdown:        s.sessionCategoryBreakdown(rows, summary.SessionCost, sessionBreakdownMetricCost),
		SessionSubCategoryCostBreakdown:     s.sessionSubCategoryBreakdown(rows, summary.SessionCost, sessionBreakdownMetricCost),
	}
}

type sessionBreakdownMetric string

const (
	sessionBreakdownMetricCount   sessionBreakdownMetric = "count"
	sessionBreakdownMetricMinutes sessionBreakdownMetric = "minutes"
	sessionBreakdownMetricCost    sessionBreakdownMetric = "cost"
)

func (s *StatsService) sessionCategoryBreakdown(rows []model.StatsSessionBreakdownRow, total float64, metric sessionBreakdownMetric) []model.StatsBreakdownItemResponse {
	valueByCategory := make(map[model.SessionCategory]float64, len(rows))
	for _, row := range rows {
		valueByCategory[row.Category] += sessionBreakdownValue(row, metric)
	}

	categories := s.sessionCategoryResolver.Categories()
	breakdown := make([]model.StatsBreakdownItemResponse, 0, len(categories))
	for _, category := range categories {
		categoryValue := model.SessionCategory(category.Value)
		value := valueByCategory[categoryValue]
		breakdown = append(breakdown, model.StatsBreakdownItemResponse{
			Key:      sessionCategoryKey(categoryValue),
			Label:    category.Label,
			Value:    value,
			Percent:  percent(value, total),
			Category: category.Value,
		})
	}
	return breakdown
}

func (s *StatsService) sessionSubCategoryBreakdown(rows []model.StatsSessionBreakdownRow, total float64, metric sessionBreakdownMetric) []model.StatsBreakdownItemResponse {
	valueBySub := make(map[string]float64, len(rows))
	for _, row := range rows {
		valueBySub[sessionSubCategoryKey(row.Category, row.SubCategory)] += sessionBreakdownValue(row, metric)
	}

	categories := s.sessionCategoryResolver.Categories()
	breakdown := make([]model.StatsBreakdownItemResponse, 0)
	for _, category := range categories {
		categoryValue := model.SessionCategory(category.Value)
		for _, subCategory := range category.SubCategories {
			subCategoryValue := model.SessionSubCategory(subCategory.Value)
			key := sessionSubCategoryKey(categoryValue, subCategoryValue)
			value := valueBySub[key]
			typeText, _ := s.sessionCategoryResolver.TypeText(categoryValue, subCategoryValue)
			breakdown = append(breakdown, model.StatsBreakdownItemResponse{
				Key:         key,
				Label:       typeText,
				Value:       value,
				Percent:     percent(value, total),
				Category:    category.Value,
				SubCategory: subCategory.Value,
			})
		}
	}
	return breakdown
}

func sessionBreakdownValue(row model.StatsSessionBreakdownRow, metric sessionBreakdownMetric) float64 {
	switch metric {
	case sessionBreakdownMetricCount:
		return float64(row.Count)
	case sessionBreakdownMetricMinutes:
		return float64(row.Minutes)
	case sessionBreakdownMetricCost:
		return row.Cost
	default:
		return 0
	}
}

func sessionCategoryKey(category model.SessionCategory) string {
	return fmt.Sprintf("category_%d", category)
}

func sessionSubCategoryKey(category model.SessionCategory, subCategory model.SessionSubCategory) string {
	return fmt.Sprintf("category_%d_sub_%d", category, subCategory)
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
