package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"tennisdaily-backend/internal/model"
	"tennisdaily-backend/internal/repository"
	"unicode/utf8"
)

const (
	maxSignatureLength        = 30
	minDefaultDurationMinutes = 30
	maxDefaultDurationMinutes = 600
)

type UserSettingsService struct {
	userSettingRepo *repository.UserSettingRepository
	userRepo        *repository.UserRepository
	contentSecurity *ContentSecurityService
}

func NewUserSettingsService(settings *repository.UserSettingRepository, users *repository.UserRepository, security *ContentSecurityService) *UserSettingsService {
	return &UserSettingsService{settings, users, security}
}
func settingsResponse(rows []model.UserSetting) (model.UserSettingsResponse, error) {
	result := model.UserSettingsResponse{}
	seen := make(map[int16]bool)
	for _, row := range rows {
		if seen[row.SettingKey] {
			return result, errors.New("duplicate active user setting")
		}
		seen[row.SettingKey] = true
		switch row.SettingKey {
		case model.UserSettingSignature:
			result.Signature = row.SettingValue
		case model.UserSettingDefaultCourtName:
			result.DefaultCourtName = row.SettingValue
		case model.UserSettingDefaultDurationMinutes:
			n, err := strconv.Atoi(row.SettingValue)
			if err != nil {
				return result, errors.New("invalid stored duration")
			}
			result.DefaultDurationMinutes = n
		case model.UserSettingDefaultCourtLocation:
			if err := json.Unmarshal([]byte(row.SettingValue), &result.DefaultCourtLocation); err != nil {
				return result, errors.New("invalid stored default location")
			}
		}
	}
	if l := result.DefaultCourtLocation; l != nil && (!l.Valid() || model.TrimCourtText(result.DefaultCourtName) == "") {
		return result, errors.New("invalid stored default court")
	}
	return result, nil
}
func (s *UserSettingsService) Get(ctx context.Context, userID int64) (model.UserSettingsResponse, error) {
	rows, err := s.userSettingRepo.ListByUser(ctx, userID)
	if err != nil {
		return model.UserSettingsResponse{}, err
	}
	return settingsResponse(rows)
}
func (s *UserSettingsService) Update(ctx context.Context, userID int64, req model.UpdateUserSettingsRequest) (model.UserSettingsResponse, error) {
	var result model.UserSettingsResponse
	updates := make([]model.UserSetting, 0, 4)
	texts := make([]string, 0, 3)
	if req.Signature != nil {
		if utf8.RuneCountInString(*req.Signature) > maxSignatureLength {
			return result, NewInvalidRequestError("个性签名最多 30 字")
		}
		texts = append(texts, *req.Signature)
		updates = append(updates, model.UserSetting{SettingKey: model.UserSettingSignature, SettingValue: *req.Signature})
	}
	if req.DefaultCourtName != nil || req.DefaultCourtLocation.Present {
		court, err := normalizeCourt(req.DefaultCourtName, req.DefaultCourtLocation)
		if err != nil {
			return result, err
		}
		value, err := json.Marshal(court.CourtLocation)
		if err != nil {
			return result, errors.New("invalid default location serialization")
		}
		updates = append(updates, model.UserSetting{SettingKey: model.UserSettingDefaultCourtName, SettingValue: court.CourtName}, model.UserSetting{SettingKey: model.UserSettingDefaultCourtLocation, SettingValue: string(value)})
		texts = append(texts, court.CourtName, courtAddress(court))
	}
	if req.DefaultDurationMinutes != nil {
		n := *req.DefaultDurationMinutes
		if n != 0 && (n < minDefaultDurationMinutes || n > maxDefaultDurationMinutes) {
			return result, NewInvalidRequestError("默认时长需在 30-600 分钟之间")
		}
		updates = append(updates, model.UserSetting{SettingKey: model.UserSettingDefaultDurationMinutes, SettingValue: strconv.Itoa(n)})
	}
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return result, err
	}
	if user == nil {
		return result, ErrNotFound
	}
	if err := s.contentSecurity.CheckTexts(user.OpenID, texts...); err != nil {
		return result, err
	}
	err = s.userSettingRepo.Apply(ctx, userID, updates, func(rows []model.UserSetting) error {
		var parseErr error
		result, parseErr = settingsResponse(rows)
		return parseErr
	})
	return result, err
}
