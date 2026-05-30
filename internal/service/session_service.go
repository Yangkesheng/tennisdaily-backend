package service

import (
	"time"

	"tennisdaily-backend/internal/model"
	"tennisdaily-backend/internal/repository"
)

const (
	defaultDurationMinutes = 120
	maxDurationMinutes     = 600
	defaultRating          = 3
)

type SessionService struct {
	repo *repository.SessionRepository
	loc  *time.Location
}

func NewSessionService(repo *repository.SessionRepository) *SessionService {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.Local
	}
	return &SessionService{repo: repo, loc: loc}
}

func (s *SessionService) List(userID int64) ([]model.SessionResponse, error) {
	sessions, err := s.repo.List(userID)
	if err != nil {
		return nil, err
	}
	return toSessionResponses(sessions), nil
}

func (s *SessionService) Create(userID int64, req model.CreateSessionRequest) (model.SessionResponse, error) {
	session, err := s.buildSession(userID, req.Date, req.DurationMinutes, req.Rating, req.Type, req.MatchRank, req.CourtName, req.Cost, req.RacketID, req.RacketName, req.ShoeName, req.Note)
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

	updated, err := s.buildSession(userID, req.Date, req.DurationMinutes, req.Rating, req.Type, req.MatchRank, req.CourtName, req.Cost, req.RacketID, req.RacketName, req.ShoeName, req.Note)
	if err != nil {
		return model.SessionResponse{}, err
	}

	existing.Date = updated.Date
	existing.DurationMinutes = updated.DurationMinutes
	existing.Rating = updated.Rating
	existing.Type = updated.Type
	existing.MatchRank = updated.MatchRank
	existing.CourtName = updated.CourtName
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

func (s *SessionService) buildSession(userID int64, date string, duration int, rating int16, sessionType model.SessionType, matchRank model.MatchRank, courtName string, cost float64, racketID int64, racketName string, shoeName string, note string) (model.TennisSession, error) {
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
