package service

import (
	"math"
	"testing"
	"time"

	"tennisdaily-backend/internal/config"
)

func TestCalculateShoeWearFresh(t *testing.T) {
	loc := mustShanghaiLocation(t)
	wear := newTestShoeWearResolver(t)
	purchaseDate := time.Date(2026, 8, 1, 10, 0, 0, 0, loc)
	now := time.Date(2026, 8, 2, 9, 0, 0, 0, loc)

	result := calculateShoeWearAt(now, purchaseDate, 0, loc, wear)

	if result.State != "fresh" {
		t.Fatalf("expected state fresh, got %q", result.State)
	}
	if result.Display != "全新 · 缓震充足" {
		t.Fatalf("expected fresh display, got %q", result.Display)
	}
	// 1 天静置：effective_wear = 0.12，score = 100 * (1 - 0.12/60) = 99.8
	if math.Abs(result.Score-99.8) > 0.001 {
		t.Fatalf("expected score 99.8, got %f", result.Score)
	}
	if result.RemainingHours != 60 {
		t.Fatalf("expected remaining hours 60, got %d", result.RemainingHours)
	}
}

func TestCalculateShoeWearWorn(t *testing.T) {
	loc := mustShanghaiLocation(t)
	wear := newTestShoeWearResolver(t)
	purchaseDate := time.Date(2026, 7, 5, 10, 0, 0, 0, loc)
	now := time.Date(2026, 8, 4, 23, 0, 0, 0, loc)

	// 30 天上场 30 小时：effective_wear = 30 + 30*0.12 = 33.6，score = 44
	result := calculateShoeWearAt(now, purchaseDate, 1800, loc, wear)

	if result.State != "worn" {
		t.Fatalf("expected state worn, got %q", result.State)
	}
	if math.Abs(result.Score-44) > 0.001 {
		t.Fatalf("expected score 44, got %f", result.Score)
	}
	if result.RemainingHours != 26 {
		t.Fatalf("expected remaining hours 26, got %d", result.RemainingHours)
	}
}

func TestCalculateShoeWearUsesConfiguredThresholds(t *testing.T) {
	loc := mustShanghaiLocation(t)
	wear, err := config.NewShoeWearResolver(config.ShoeWearConfig{
		StandardLifeHours: 20,
		RestWearPerDay:    0.1,
		States: []config.ShoeWearStateConfig{
			{Key: "fresh", MinScore: 90, Display: "刚入手", Color: "green"},
			{Key: "expired", MinScore: 0, Display: "已报废", Color: "red"},
		},
	})
	if err != nil {
		t.Fatalf("build resolver: %v", err)
	}

	// 3 天静置、未上场：effective_wear = 0.3，score = 98.5，落在 fresh。
	purchaseDate := time.Date(2026, 8, 1, 10, 0, 0, 0, loc)
	now := time.Date(2026, 8, 4, 9, 0, 0, 0, loc)
	result := calculateShoeWearAt(now, purchaseDate, 0, loc, wear)

	if result.State != "fresh" {
		t.Fatalf("expected state fresh, got %q", result.State)
	}
	if result.Display != "刚入手" {
		t.Fatalf("expected configured display, got %q", result.Display)
	}
	if result.RemainingHours != 20 {
		t.Fatalf("expected remaining hours 20, got %d", result.RemainingHours)
	}
}

func newTestShoeWearResolver(t *testing.T) *config.ShoeWearResolver {
	t.Helper()

	wear, err := config.NewShoeWearResolver(config.ShoeWearConfig{
		StandardLifeHours: 60,
		RestWearPerDay:    0.12,
		States: []config.ShoeWearStateConfig{
			{Key: "fresh", MinScore: 85, Label: "全新", Display: "全新 · 缓震充足", Color: "green"},
			{Key: "good", MinScore: 65, Label: "状态良好", Display: "状态良好 · 预计还可打 {remainingHours}h", Color: "teal"},
			{Key: "worn", MinScore: 40, Label: "开始磨损", Display: "开始磨损 · 建议关注中底", Color: "orange"},
			{Key: "tired", MinScore: 20, Label: "磨损明显", Display: "磨损明显 · 建议考虑换鞋", Color: "deep_orange"},
			{Key: "dead", MinScore: 8, Label: "缓震衰减", Display: "缓震衰减 · 建议尽快换鞋", Color: "red"},
			{Key: "expired", MinScore: 0, Label: "已超期", Display: "已超期 · 不建议继续使用", Color: "red"},
		},
	})
	if err != nil {
		t.Fatalf("build resolver: %v", err)
	}
	return wear
}
