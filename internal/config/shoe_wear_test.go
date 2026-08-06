package config

import (
	"strings"
	"testing"
)

func TestNewShoeWearResolverDefaults(t *testing.T) {
	resolver, err := NewShoeWearResolver(ShoeWearConfig{})
	if err != nil {
		t.Fatalf("build resolver with defaults: %v", err)
	}

	if resolver.StandardLifeHours() != 60 {
		t.Fatalf("expected standard life hours 60, got %v", resolver.StandardLifeHours())
	}
	if resolver.RestWearPerDay() != 0.06 {
		t.Fatalf("expected rest wear per day 0.06, got %v", resolver.RestWearPerDay())
	}

	state := resolver.State(95)
	if state.Key != "fresh" || state.Display != "全新 · 缓震充足" {
		t.Fatalf("expected fresh state, got %+v", state)
	}
	state = resolver.State(50)
	if state.Key != "worn" {
		t.Fatalf("expected worn state, got %+v", state)
	}
	state = resolver.State(0)
	if state.Key != "expired" {
		t.Fatalf("expected expired fallback, got %+v", state)
	}

	if got := resolver.RenderDisplay("状态良好 · 预计还可打 {remainingHours}h", 12); got != "状态良好 · 预计还可打 12h" {
		t.Fatalf("expected rendered display, got %q", got)
	}
}

func TestNewShoeWearResolverValidation(t *testing.T) {
	base := func() ShoeWearConfig {
		return ShoeWearConfig{
			StandardLifeHours: 60,
			RestWearPerDay:    0.12,
			States: []ShoeWearStateConfig{
				{Key: "fresh", MinScore: 85, Display: "全新"},
				{Key: "expired", MinScore: 0, Display: "已超期"},
			},
		}
	}

	tests := []struct {
		name    string
		mutate  func(*ShoeWearConfig)
		wantErr string
	}{
		{
			name: "standard life hours must be positive",
			mutate: func(cfg *ShoeWearConfig) {
				cfg.StandardLifeHours = 0
			},
			wantErr: "standardLifeHours",
		},
		{
			name: "duplicate state key",
			mutate: func(cfg *ShoeWearConfig) {
				cfg.States[1].Key = "fresh"
			},
			wantErr: "duplicate",
		},
		{
			name: "no fallback state",
			mutate: func(cfg *ShoeWearConfig) {
				cfg.States[1].MinScore = 10
			},
			wantErr: "fallback",
		},
		{
			name: "empty display",
			mutate: func(cfg *ShoeWearConfig) {
				cfg.States[0].Display = ""
			},
			wantErr: "display is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := base()
			tt.mutate(&cfg)
			_, err := NewShoeWearResolver(cfg)
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestRepoConfigShoeWear(t *testing.T) {
	t.Setenv("CONFIG_PATH", "../../config.yaml")

	fc := loadConfigFile()
	if fc.ShoeWear.StandardLifeHours != 60 {
		t.Fatalf("expected config.yaml standardLifeHours 60, got %v", fc.ShoeWear.StandardLifeHours)
	}
	if len(fc.ShoeWear.States) != 6 {
		t.Fatalf("expected 6 states in config.yaml, got %d", len(fc.ShoeWear.States))
	}
}
