package config

import (
	"strings"
	"testing"
)

func TestNewPolyesterStringHealthResolverDefaults(t *testing.T) {
	resolver, err := NewPolyesterStringHealthResolver(PolyesterStringHealthConfig{})
	if err != nil {
		t.Fatalf("build resolver with defaults: %v", err)
	}

	if resolver.StandardLifeHours() != 16 {
		t.Fatalf("expected standard life hours 16, got %v", resolver.StandardLifeHours())
	}
	if resolver.RestWearPerDay() != 0.18 {
		t.Fatalf("expected rest wear per day 0.18, got %v", resolver.RestWearPerDay())
	}

	state := resolver.State(97.75)
	if state.Key != "fresh" || state.Display != "新上线 · 手感正脆" {
		t.Fatalf("expected fresh state, got %+v", state)
	}
	state = resolver.State(46.625)
	if state.Key != "good" {
		t.Fatalf("expected good state, got %+v", state)
	}
	state = resolver.State(0)
	if state.Key != "expired" {
		t.Fatalf("expected expired fallback, got %+v", state)
	}

	if got := resolver.RenderDisplay("状态良好 · 预计还可打 {remainingHours}h", 7); got != "状态良好 · 预计还可打 7h" {
		t.Fatalf("expected rendered display, got %q", got)
	}
}

func TestNewPolyesterStringHealthResolverValidation(t *testing.T) {
	base := func() PolyesterStringHealthConfig {
		return PolyesterStringHealthConfig{
			StandardLifeHours: 16,
			RestWearPerDay:    0.18,
			States: []PolyesterStringHealthStateConfig{
				{Key: "fresh", MinScore: 85, Display: "新上线"},
				{Key: "expired", MinScore: 0, Display: "已超期"},
			},
		}
	}

	tests := []struct {
		name    string
		mutate  func(*PolyesterStringHealthConfig)
		wantErr string
	}{
		{
			name: "standard life hours must be positive",
			mutate: func(cfg *PolyesterStringHealthConfig) {
				cfg.StandardLifeHours = 0
			},
			wantErr: "standardLifeHours",
		},
		{
			name: "duplicate state key",
			mutate: func(cfg *PolyesterStringHealthConfig) {
				cfg.States[1].Key = "fresh"
			},
			wantErr: "duplicate",
		},
		{
			name: "no fallback state",
			mutate: func(cfg *PolyesterStringHealthConfig) {
				cfg.States[1].MinScore = 10
			},
			wantErr: "fallback",
		},
		{
			name: "empty display",
			mutate: func(cfg *PolyesterStringHealthConfig) {
				cfg.States[0].Display = ""
			},
			wantErr: "display is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := base()
			tt.mutate(&cfg)
			_, err := NewPolyesterStringHealthResolver(cfg)
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestRepoConfigPolyesterStringHealth(t *testing.T) {
	t.Setenv("CONFIG_PATH", "../../config.yaml")

	fc := loadConfigFile()
	if fc.PolyesterStringHealth.StandardLifeHours != 16 {
		t.Fatalf("expected config.yaml standardLifeHours 16, got %v", fc.PolyesterStringHealth.StandardLifeHours)
	}
	if len(fc.PolyesterStringHealth.States) != 6 {
		t.Fatalf("expected 6 states in config.yaml, got %d", len(fc.PolyesterStringHealth.States))
	}
}
