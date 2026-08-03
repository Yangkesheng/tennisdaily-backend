package config

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// PolyesterStringHealthConfig 是聚酯线健康度配置，对应 config.yaml 的 polyesterStringHealth 段。
// 配置缺失时使用代码内置默认值，配置存在时校验并通过 resolver 提供给业务层。
type PolyesterStringHealthConfig struct {
	StandardLifeHours float64                            `yaml:"standardLifeHours"`
	RestWearPerDay    float64                            `yaml:"restWearPerDay"`
	States            []PolyesterStringHealthStateConfig `yaml:"states"`
}

type PolyesterStringHealthStateConfig struct {
	Key      string  `yaml:"key"`
	MinScore float64 `yaml:"minScore"`
	Label    string  `yaml:"label"`
	Display  string  `yaml:"display"`
	Color    string  `yaml:"color"`
}

type PolyesterStringHealthState struct {
	Key     string
	Label   string
	Display string
	Color   string
}

type PolyesterStringHealthResolver struct {
	standardLifeHours float64
	restWearPerDay    float64
	states            []PolyesterStringHealthStateConfig
}

func NewPolyesterStringHealthResolver(cfg PolyesterStringHealthConfig) (*PolyesterStringHealthResolver, error) {
	cfg = applyPolyesterStringHealthDefaults(cfg)
	if cfg.StandardLifeHours <= 0 {
		return nil, fmt.Errorf("polyesterStringHealth.standardLifeHours must be > 0")
	}
	if cfg.RestWearPerDay < 0 {
		return nil, fmt.Errorf("polyesterStringHealth.restWearPerDay must be >= 0")
	}
	if len(cfg.States) == 0 {
		return nil, fmt.Errorf("polyesterStringHealth.states is required")
	}

	seenKeys := make(map[string]struct{}, len(cfg.States))
	seenScores := make(map[float64]struct{}, len(cfg.States))
	states := make([]PolyesterStringHealthStateConfig, 0, len(cfg.States))
	for _, state := range cfg.States {
		if state.Key == "" {
			return nil, fmt.Errorf("polyesterStringHealth.states[].key is required")
		}
		if _, exists := seenKeys[state.Key]; exists {
			return nil, fmt.Errorf("duplicate polyesterStringHealth state key: %q", state.Key)
		}
		seenKeys[state.Key] = struct{}{}
		if state.MinScore < 0 || state.MinScore > 100 {
			return nil, fmt.Errorf("polyesterStringHealth state %q minScore must be in [0, 100]", state.Key)
		}
		if _, exists := seenScores[state.MinScore]; exists {
			return nil, fmt.Errorf("duplicate polyesterStringHealth state minScore: %v", state.MinScore)
		}
		seenScores[state.MinScore] = struct{}{}
		if state.Display == "" {
			return nil, fmt.Errorf("polyesterStringHealth state %q display is required", state.Key)
		}
		states = append(states, state)
	}

	sort.Slice(states, func(i, j int) bool {
		return states[i].MinScore > states[j].MinScore
	})
	if states[len(states)-1].MinScore != 0 {
		return nil, fmt.Errorf("polyesterStringHealth.states must include a fallback state with minScore 0")
	}

	return &PolyesterStringHealthResolver{
		standardLifeHours: cfg.StandardLifeHours,
		restWearPerDay:    cfg.RestWearPerDay,
		states:            states,
	}, nil
}

// StandardLifeHours 返回聚酯线标准可用寿命（小时）。
func (r *PolyesterStringHealthResolver) StandardLifeHours() float64 {
	return r.standardLifeHours
}

// RestWearPerDay 返回聚酯线每日静置衰减系数。
func (r *PolyesterStringHealthResolver) RestWearPerDay() float64 {
	return r.restWearPerDay
}

// State 按 score 返回命中的状态配置，score 小于所有 minScore 时返回兜底状态。
func (r *PolyesterStringHealthResolver) State(score float64) PolyesterStringHealthState {
	for _, state := range r.states {
		if score >= state.MinScore {
			return toPolyesterStringHealthState(state)
		}
	}
	return toPolyesterStringHealthState(r.states[len(r.states)-1])
}

// RenderDisplay 将 display 模板中的 {remainingHours} 替换为实际剩余小时数。
func (r *PolyesterStringHealthResolver) RenderDisplay(display string, remainingHours int) string {
	return strings.ReplaceAll(display, "{remainingHours}", strconv.Itoa(remainingHours))
}

func toPolyesterStringHealthState(state PolyesterStringHealthStateConfig) PolyesterStringHealthState {
	return PolyesterStringHealthState{
		Key:     state.Key,
		Label:   state.Label,
		Display: state.Display,
		Color:   state.Color,
	}
}

func applyPolyesterStringHealthDefaults(cfg PolyesterStringHealthConfig) PolyesterStringHealthConfig {
	if len(cfg.States) == 0 {
		cfg.StandardLifeHours = 16
		cfg.RestWearPerDay = 0.18
		cfg.States = []PolyesterStringHealthStateConfig{
			{Key: "fresh", MinScore: 85, Label: "新上线", Display: "新上线 · 手感正脆", Color: "green"},
			{Key: "peak", MinScore: 70, Label: "巅峰期", Display: "巅峰期 · 预计还可打 {remainingHours}h", Color: "green"},
			{Key: "good", MinScore: 45, Label: "状态良好", Display: "状态良好 · 预计还可打 {remainingHours}h", Color: "teal"},
			{Key: "decline", MinScore: 25, Label: "开始衰减", Display: "开始衰减 · 建议近期重穿", Color: "orange"},
			{Key: "dead", MinScore: 10, Label: "手感变死", Display: "手感变死 · 建议重穿", Color: "deep_orange"},
			{Key: "expired", MinScore: 0, Label: "已超期", Display: "已超期 · 不建议比赛使用", Color: "red"},
		}
	}
	return cfg
}
