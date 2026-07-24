package service

import (
	"errors"
	"testing"

	"tennisdaily-backend/internal/config"
	"tennisdaily-backend/internal/model"
)

func TestSessionServiceCreateRejectsDurationOverMax(t *testing.T) {
	resolver := newTestSessionCategoryResolver(t)
	service := NewSessionService(nil, nil, nil, resolver, nil)

	_, err := service.Create(29, model.CreateSessionRequest{
		Date:            "2026-07-25 03:00",
		DurationMinutes: 1800,
		Rating:          3,
		Category:        model.SessionCategoryDaily,
		SubCategory:     model.SessionSubCategoryDoubles,
		MatchRank:       model.MatchRankNone,
		RacketID:        15,
		RacketName:      "Pure Aero 98 Tour racket unstrung",
	})
	if !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("expected ErrInvalidRequest, got %v", err)
	}
	if got, want := InvalidRequestMessage(err), "打球时长必须在 1-600 分钟之间"; got != want {
		t.Fatalf("expected message %q, got %q", want, got)
	}
}

func newTestSessionCategoryResolver(t *testing.T) *config.SessionCategoryResolver {
	t.Helper()

	resolver, err := config.NewSessionCategoryResolver(config.SessionEnumsConfig{
		Categories: []config.SessionCategoryConfig{
			{
				Value: int16(model.SessionCategoryDaily),
				Label: "日常球局",
				SubCategories: []config.SessionSubCategoryConfig{
					{Value: int16(model.SessionSubCategorySingles), Label: "打单"},
					{Value: int16(model.SessionSubCategoryDoubles), Label: "双打"},
				},
			},
			{
				Value: int16(model.SessionCategoryTraining),
				Label: "训练",
				SubCategories: []config.SessionSubCategoryConfig{
					{Value: int16(model.SessionSubCategoryServe), Label: "发球"},
					{Value: int16(model.SessionSubCategoryOther), Label: "其他"},
				},
			},
			{
				Value: int16(model.SessionCategoryMatch),
				Label: "比赛",
				SubCategories: []config.SessionSubCategoryConfig{
					{Value: int16(model.SessionSubCategorySingles), Label: "单打"},
					{Value: int16(model.SessionSubCategoryDoubles), Label: "双打"},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("new session category resolver: %v", err)
	}
	return resolver
}
