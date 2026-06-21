package service

import (
	"fmt"
	"time"

	"tennisdaily-backend/internal/model"
	"tennisdaily-backend/internal/repository"
)

const (
	defaultDurationMinutes = 120
	maxDurationMinutes     = 600
	defaultRating          = 3
	defaultSessionPage     = 1
	defaultSessionPageSize = 20
	maxSessionPageSize     = 100
)

type SessionService struct {
	repo       *repository.SessionRepository
	racketRepo *repository.RacketRepository
	loc        *time.Location
}

func NewSessionService(repo *repository.SessionRepository, racketRepo *repository.RacketRepository) *SessionService {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.Local
	}
	return &SessionService{repo: repo, racketRepo: racketRepo, loc: loc}
}

func (s *SessionService) List(userID int64) ([]model.SessionResponse, error) {
	sessions, err := s.repo.List(userID)
	if err != nil {
		return nil, err
	}
	return toSessionResponses(sessions), nil
}

func (s *SessionService) ListPage(userID int64, query model.SessionListQuery) (model.SessionListPageResponse, error) {
	if query.Date != "" {
		responses, err := s.ListByDate(userID, query.Date)
		if err != nil {
			return model.SessionListPageResponse{}, err
		}
		total := int64(len(responses))
		return model.SessionListPageResponse{
			List:       responses,
			Total:      total,
			Page:       1,
			PageSize:   len(responses),
			TotalPages: 1,
			HasMore:    false,
		}, nil
	}

	page := query.Page
	if page == 0 {
		page = defaultSessionPage
	}
	pageSize := query.PageSize
	if pageSize == 0 {
		pageSize = defaultSessionPageSize
	}
	if page < 1 || pageSize < 1 || pageSize > maxSessionPageSize {
		return model.SessionListPageResponse{}, ErrInvalidRequest
	}

	sessions, total, err := s.repo.ListPage(userID, page, pageSize)
	if err != nil {
		return model.SessionListPageResponse{}, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}

	return model.SessionListPageResponse{
		List:       toSessionResponses(sessions),
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		HasMore:    int64(page*pageSize) < total,
	}, nil
}

func (s *SessionService) ListByDate(userID int64, date string) ([]model.SessionResponse, error) {
	sessionDate, err := time.ParseInLocation("2006-01-02", date, s.loc)
	if err != nil {
		return nil, ErrInvalidRequest
	}

	sessions, err := s.repo.ListByDate(userID, sessionDate)
	if err != nil {
		return nil, err
	}
	return toSessionResponses(sessions), nil
}

func (s *SessionService) Create(userID int64, req model.CreateSessionRequest) (model.SessionResponse, error) {
	session, err := s.buildSession(userID, req.Date, req.DurationMinutes, req.Rating, req.Type, req.MatchRank, req.CourtName, req.Partner, req.Cost, req.RacketID, req.RacketName, req.ShoeName, req.Note)
	if err != nil {
		return model.SessionResponse{}, err
	}
	if err := s.repo.Create(&session); err != nil {
		return model.SessionResponse{}, err
	}
	return model.NewSessionResponse(session), nil
}

func (s *SessionService) FindByID(userID, id int64) (model.SessionResponse, error) {
	session, err := s.repo.FindByID(userID, id)
	if err != nil {
		return model.SessionResponse{}, err
	}
	if session == nil {
		return model.SessionResponse{}, ErrNotFound
	}
	return model.NewSessionResponse(*session), nil
}

func (s *SessionService) Update(userID, id int64, req model.UpdateSessionRequest) (model.SessionResponse, error) {
	existing, err := s.repo.FindByID(userID, id)
	if err != nil {
		return model.SessionResponse{}, err
	}
	if existing == nil {
		return model.SessionResponse{}, ErrNotFound
	}

	updated, err := s.buildSession(userID, req.Date, req.DurationMinutes, req.Rating, req.Type, req.MatchRank, req.CourtName, req.Partner, req.Cost, req.RacketID, req.RacketName, req.ShoeName, req.Note)
	if err != nil {
		return model.SessionResponse{}, err
	}

	existing.Date = updated.Date
	existing.DurationMinutes = updated.DurationMinutes
	existing.Rating = updated.Rating
	existing.Type = updated.Type
	existing.MatchRank = updated.MatchRank
	existing.CourtName = updated.CourtName
	existing.Partner = updated.Partner
	existing.Cost = updated.Cost
	existing.RacketID = updated.RacketID
	existing.RacketName = updated.RacketName
	existing.ShoeName = updated.ShoeName
	existing.Note = updated.Note

	if err := s.repo.Update(existing); err != nil {
		return model.SessionResponse{}, err
	}
	return model.NewSessionResponse(*existing), nil
}

func (s *SessionService) Delete(userID, id int64) error {
	deleted, err := s.repo.SoftDelete(userID, id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrNotFound
	}
	return nil
}

func (s *SessionService) Latest(userID int64) (*model.SessionResponse, error) {
	session, err := s.repo.Latest(userID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, nil
	}
	resp := model.NewSessionResponse(*session)
	return &resp, nil
}

func (s *SessionService) Calendar(userID int64, year, month int) (model.SessionCalendarResponse, error) {
	if year < 2000 || year > 2100 || month < 0 || month > 12 {
		return model.SessionCalendarResponse{}, ErrInvalidRequest
	}

	start := time.Date(year, time.January, 1, 0, 0, 0, 0, s.loc)
	end := start.AddDate(1, 0, 0)
	if month > 0 {
		start = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, s.loc)
		end = start.AddDate(0, 1, 0)
	}
	days, err := s.repo.CalendarDays(userID, start, end)
	if err != nil {
		return model.SessionCalendarResponse{}, err
	}

	summary, err := s.calendarSummary(userID, start, end)
	if err != nil {
		return model.SessionCalendarResponse{}, err
	}

	ratingTrend, err := s.calendarRatingTrend(userID, start, end)
	if err != nil {
		return model.SessionCalendarResponse{}, err
	}

	return model.SessionCalendarResponse{
		Year:           year,
		Month:          month,
		ActiveDayCount: len(days),
		Days:           days,
		Summary:        summary,
		Charts: model.SessionCalendarChartsResponse{
			WeeklySessions:       s.calendarWeeklySessions(start, end, days),
			RatingTrend:          ratingTrend,
			ExpenseBreakdown:     calendarExpenseBreakdown(summary.SessionCost, summary.RacketCost, summary.StringingCost),
			SessionTypeBreakdown: calendarSessionTypeBreakdown(summary),
		},
	}, nil
}

func (s *SessionService) calendarSummary(userID int64, start, end time.Time) (model.SessionCalendarSummaryResponse, error) {
	sessionStats, err := s.repo.StatsAggregate(userID, start, end)
	if err != nil {
		return model.SessionCalendarSummaryResponse{}, err
	}

	racketCost, err := s.racketRepo.SumPurchaseCostByRange(userID, start, end)
	if err != nil {
		return model.SessionCalendarSummaryResponse{}, err
	}

	stringingCost, err := s.racketRepo.SumStringingCostByRange(userID, start, end)
	if err != nil {
		return model.SessionCalendarSummaryResponse{}, err
	}

	return model.SessionCalendarSummaryResponse{
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

func (s *SessionService) calendarWeeklySessions(start, end time.Time, days []model.SessionCalendarDay) []model.CalendarWeeklySessionChartItemResponse {
	mondayStart := start.AddDate(0, 0, -weekdayOffset(start))
	weeks := 0
	for weekStart := mondayStart; weekStart.Before(end); weekStart = weekStart.AddDate(0, 0, 7) {
		weeks++
	}

	weeklySessions := make([]model.CalendarWeeklySessionChartItemResponse, weeks)
	for index := range weeklySessions {
		weeklySessions[index] = model.CalendarWeeklySessionChartItemResponse{
			Label: fmt.Sprintf("第%d周", index+1),
			Count: 0,
		}
	}

	for _, day := range days {
		date, err := time.ParseInLocation("2006-01-02", day.Date, s.loc)
		if err != nil {
			continue
		}
		index := int(date.Sub(mondayStart).Hours() / 24 / 7)
		if index >= 0 && index < len(weeklySessions) {
			weeklySessions[index].Count += int64(day.Count)
		}
	}

	return weeklySessions
}

func (s *SessionService) calendarRatingTrend(userID int64, start, end time.Time) ([]model.CalendarRatingTrendChartItemResponse, error) {
	sessions, err := s.repo.ListRatingTrendByRange(userID, start, end, statsRatingTrendLimit)
	if err != nil {
		return nil, err
	}

	trend := make([]model.CalendarRatingTrendChartItemResponse, 0, len(sessions))
	for index := len(sessions) - 1; index >= 0; index-- {
		session := sessions[index]
		trend = append(trend, model.CalendarRatingTrendChartItemResponse{
			Date:   session.Date.Format("2006-01-02"),
			Rating: session.Rating,
		})
	}
	return trend, nil
}

func calendarExpenseBreakdown(sessionCost, racketCost, stringingCost float64) []model.CalendarExpenseChartItemResponse {
	return []model.CalendarExpenseChartItemResponse{
		{Label: "打球", Value: sessionCost},
		{Label: "球拍", Value: racketCost},
		{Label: "穿线", Value: stringingCost},
	}
}

func calendarSessionTypeBreakdown(summary model.SessionCalendarSummaryResponse) []model.StatsBreakdownItemResponse {
	total := float64(summary.SessionCount)
	return []model.StatsBreakdownItemResponse{
		{Key: "training", Label: "训练", Value: float64(summary.TrainingCount), Percent: percent(float64(summary.TrainingCount), total)},
		{Key: "singles", Label: "单打", Value: float64(summary.SinglesCount), Percent: percent(float64(summary.SinglesCount), total)},
		{Key: "doubles", Label: "双打", Value: float64(summary.DoublesCount), Percent: percent(float64(summary.DoublesCount), total)},
		{Key: "match", Label: "比赛", Value: float64(summary.MatchCount), Percent: percent(float64(summary.MatchCount), total)},
	}
}

func (s *SessionService) buildSession(userID int64, date string, duration int, rating int16, sessionType model.SessionType, matchRank model.MatchRank, courtName string, partner string, cost float64, racketID int64, racketName string, shoeName string, note string) (model.TennisSession, error) {
	sessionDate, err := time.ParseInLocation("2006-01-02", date, s.loc)
	if err != nil {
		return model.TennisSession{}, ErrInvalidRequest
	}

	if duration == 0 {
		duration = defaultDurationMinutes
	}
	if duration <= 0 || duration > maxDurationMinutes {
		return model.TennisSession{}, ErrInvalidRequest
	}

	if rating == 0 {
		rating = defaultRating
	}
	if rating < 1 || rating > 5 {
		return model.TennisSession{}, ErrInvalidRequest
	}

	if !sessionType.IsValid() || !matchRank.IsValid() {
		return model.TennisSession{}, ErrInvalidRequest
	}
	if !sessionType.IsMatch() {
		matchRank = model.MatchRankNone
	}

	return model.TennisSession{
		UserID:          userID,
		Date:            sessionDate,
		DurationMinutes: duration,
		Rating:          rating,
		Type:            sessionType,
		MatchRank:       matchRank,
		CourtName:       courtName,
		Partner:         partner,
		Cost:            cost,
		RacketID:        racketID,
		RacketName:      racketName,
		ShoeName:        shoeName,
		Note:            note,
	}, nil
}

func toSessionResponses(sessions []model.TennisSession) []model.SessionResponse {
	responses := make([]model.SessionResponse, 0, len(sessions))
	for _, session := range sessions {
		responses = append(responses, model.NewSessionResponse(session))
	}
	return responses
}
