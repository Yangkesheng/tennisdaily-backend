package service

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"tennisdaily-backend/internal/model"
	"tennisdaily-backend/internal/repository"
)

const (
	maxSignatureLength        = 30
	maxCourtNameLength        = 20
	minDefaultDurationMinutes = 30
	maxDefaultDurationMinutes = 600
)

type UserSettingsService struct {
	userSettingRepo *repository.UserSettingRepository
	userRepo        *repository.UserRepository
	contentSecurity *ContentSecurityService
}

func NewUserSettingsService(userSettingRepo *repository.UserSettingRepository, userRepo *repository.UserRepository, contentSecurity *ContentSecurityService) *UserSettingsService {
	return &UserSettingsService{
		userSettingRepo: userSettingRepo,
		userRepo:        userRepo,
		contentSecurity: contentSecurity,
	}
}

func (s *UserSettingsService) Get(userID int64) (model.UserSettingsResponse, error) {
	settings, err := s.userSettingRepo.ListByUser(userID)
	if err != nil {
		return model.UserSettingsResponse{}, err
	}

	response := model.UserSettingsResponse{}
	for _, setting := range settings {
		switch setting.SettingKey {
		case model.UserSettingSignature:
			response.Signature = setting.SettingValue
		case model.UserSettingDefaultCourtName:
			response.DefaultCourtName = setting.SettingValue
		case model.UserSettingDefaultDurationMinutes:
			if value, parseErr := strconv.Atoi(setting.SettingValue); parseErr == nil {
				response.DefaultDurationMinutes = value
			}
		}
	}
	return response, nil
}

func (s *UserSettingsService) Update(userID int64, req model.UpdateUserSettingsRequest) (model.UserSettingsResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return model.UserSettingsResponse{}, err
	}
	if user == nil {
		return model.UserSettingsResponse{}, ErrNotFound
	}

	if req.Signature != nil {
		if utf8.RuneCountInString(*req.Signature) > maxSignatureLength {
			return model.UserSettingsResponse{}, ErrInvalidRequest
		}
		if err := s.contentSecurity.CheckTexts(user.OpenID, *req.Signature); err != nil {
			return model.UserSettingsResponse{}, err
		}
		if err := s.userSettingRepo.Upsert(userID, model.UserSettingSignature, *req.Signature); err != nil {
			return model.UserSettingsResponse{}, err
		}
	}

	if req.DefaultCourtName != nil {
		courtName := strings.TrimSpace(*req.DefaultCourtName)
		if utf8.RuneCountInString(courtName) > maxCourtNameLength {
			return model.UserSettingsResponse{}, ErrInvalidRequest
		}
		if err := s.contentSecurity.CheckTexts(user.OpenID, courtName); err != nil {
			return model.UserSettingsResponse{}, err
		}
		if err := s.userSettingRepo.Upsert(userID, model.UserSettingDefaultCourtName, courtName); err != nil {
			return model.UserSettingsResponse{}, err
		}
	}

	if req.DefaultDurationMinutes != nil {
		duration := *req.DefaultDurationMinutes
		if duration != 0 && (duration < minDefaultDurationMinutes || duration > maxDefaultDurationMinutes) {
			return model.UserSettingsResponse{}, ErrInvalidRequest
		}
		if err := s.userSettingRepo.Upsert(userID, model.UserSettingDefaultDurationMinutes, strconv.Itoa(duration)); err != nil {
			return model.UserSettingsResponse{}, err
		}
	}

	return s.Get(userID)
}
