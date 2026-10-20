package service

import (
	"errors"
	"tennisdaily-backend/internal/model"
	"testing"
)

func TestLibraryVisibilityCannotBeOverriddenByClient(t *testing.T) {
	s := &RacketService{adminUserIDs: []int64{7}}
	pending, approved := model.RacketReviewPending, model.RacketReviewApproved
	for _, requested := range []*int{nil, &pending, &approved} {
		q := s.visibleLibraryQuery(8, model.RacketLibraryQuery{IsReviewed: requested, BrandID: 3})
		if q.IsReviewed == nil || *q.IsReviewed != approved || q.BrandID != 3 {
			t.Fatal("ordinary user must only see approved entries")
		}
		admin := s.visibleLibraryQuery(7, model.RacketLibraryQuery{IsReviewed: requested})
		if admin.IsReviewed != requested {
			t.Fatal("admin filter was changed")
		}
	}
}

func TestReviewLibraryRejectsUnauthorizedOrInvalidInput(t *testing.T) {
	s := &RacketService{adminUserIDs: []int64{7}}
	pending := 0
	if _, err := s.ReviewLibrary(8, 1, model.ReviewRacketLibraryRequest{IsReviewed: &pending}); !errors.Is(err, ErrForbidden) {
		t.Fatalf("unauthorized review: %v", err)
	}
	for _, value := range []int{-1, 2} {
		if _, err := s.ReviewLibrary(7, 1, model.ReviewRacketLibraryRequest{IsReviewed: &value}); !errors.Is(err, ErrInvalidRequest) {
			t.Fatalf("invalid review status: %v", err)
		}
	}
	if _, err := s.ReviewLibrary(7, 1, model.ReviewRacketLibraryRequest{}); !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("missing status: %v", err)
	}
}
