package service

import (
	"math"
	"testing"
	"time"

	"tennisdaily-backend/internal/config"
)

func TestCalculateStringHealthFresh(t *testing.T) {
	loc := mustShanghaiLocation(t)
	health := newTestPolyesterStringHealthResolver(t)
	stringDate := time.Date(2026, 8, 2, 19, 30, 0, 0, loc)
	now := time.Date(2026, 8, 4, 9, 0, 0, 0, loc)

	result := calculateStringHealthAt(now, stringDate, 0, loc, health)

	if result.State != "fresh" {
		t.Fatalf("expected state fresh, got %q", result.State)
	}
	if result.Display != "新上线 · 手感正脆" {
		t.Fatalf("expected fresh display, got %q", result.Display)
	}
	if math.Abs(result.Score-97.75) > 0.001 {
		t.Fatalf("expected score 97.75, got %f", result.Score)
	}
	if result.RemainingHours != 16 {
		t.Fatalf("expected remaining hours 16, got %d", result.RemainingHours)
	}
}

func TestCalculateStringHealthGood(t *testing.T) {
	loc := mustShanghaiLocation(t)
	health := newTestPolyesterStringHealthResolver(t)
	stringDate := time.Date(2026, 8, 1, 10, 0, 0, 0, loc)
	now := time.Date(2026, 8, 4, 23, 0, 0, 0, loc)

	result := calculateStringHealthAt(now, stringDate, 480, loc, health)

	if result.State != "good" {
		t.Fatalf("expected state good, got %q", result.State)
	}
	if result.Display != "状态良好 · 预计还可打 7h" {
		t.Fatalf("expected good display, got %q", result.Display)
	}
	if math.Abs(result.Score-46.625) > 0.001 {
		t.Fatalf("expected score 46.625, got %f", result.Score)
	}
	if result.RemainingHours != 7 {
		t.Fatalf("expected remaining hours 7, got %d", result.RemainingHours)
	}
}

func TestCalculateStringHealthUsesConfiguredThresholds(t *testing.T) {
	loc := mustShanghaiLocation(t)
	health, err := config.NewPolyesterStringHealthResolver(config.PolyesterStringHealthConfig{
		StandardLifeHours: 20,
		RestWearPerDay:    0.1,
		States: []config.PolyesterStringHealthStateConfig{
			{Key: "fresh", MinScore: 90, Display: "刚上线", Color: "green"},
			{Key: "expired", MinScore: 0, Display: "已报废", Color: "red"},
		},
	})
	if err != nil {
		t.Fatalf("build resolver: %v", err)
	}

	// 3 天静置、未打球：effective_wear = 0.3，score = 100 * (1 - 0.3/20) = 98.5，落在 fresh。
	stringDate := time.Date(2026, 8, 1, 10, 0, 0, 0, loc)
	now := time.Date(2026, 8, 4, 9, 0, 0, 0, loc)
	result := calculateStringHealthAt(now, stringDate, 0, loc, health)

	if result.State != "fresh" {
		t.Fatalf("expected state fresh, got %q", result.State)
	}
	if result.Display != "刚上线" {
		t.Fatalf("expected configured display, got %q", result.Display)
	}
	if result.RemainingHours != 20 {
		t.Fatalf("expected remaining hours 20, got %d", result.RemainingHours)
	}
}

func mustShanghaiLocation(t *testing.T) *time.Location {
	t.Helper()

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatalf("load Asia/Shanghai: %v", err)
	}
	return loc
}

func newTestPolyesterStringHealthResolver(t *testing.T) *config.PolyesterStringHealthResolver {
	t.Helper()

	health, err := config.NewPolyesterStringHealthResolver(config.PolyesterStringHealthConfig{
		StandardLifeHours: 16,
		RestWearPerDay:    0.18,
		States: []config.PolyesterStringHealthStateConfig{
			{Key: "fresh", MinScore: 85, Label: "新上线", Display: "新上线 · 手感正脆", Color: "green"},
			{Key: "peak", MinScore: 70, Label: "巅峰期", Display: "巅峰期 · 预计还可打 {remainingHours}h", Color: "green"},
			{Key: "good", MinScore: 45, Label: "状态良好", Display: "状态良好 · 预计还可打 {remainingHours}h", Color: "teal"},
			{Key: "decline", MinScore: 25, Label: "开始衰减", Display: "开始衰减 · 建议近期重穿", Color: "orange"},
			{Key: "dead", MinScore: 10, Label: "手感变死", Display: "手感变死 · 建议重穿", Color: "deep_orange"},
			{Key: "expired", MinScore: 0, Label: "已超期", Display: "已超期 · 不建议比赛使用", Color: "red"},
		},
	})
	if err != nil {
		t.Fatalf("build resolver: %v", err)
	}
	return health
}
