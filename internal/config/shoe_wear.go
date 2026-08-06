package config

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// ShoeWearConfig 是球鞋磨损度配置，对应 config.yaml 的 shoeWear 段。
// 配置缺失时使用代码内置默认值，配置存在时校验并通过 resolver 提供给业务层。
type ShoeWearConfig struct {
	StandardLifeHours float64               `yaml:"standardLifeHours"`
	RestWearPerDay    float64               `yaml:"restWearPerDay"`
	States            []ShoeWearStateConfig `yaml:"states"`
}

type ShoeWearStateConfig struct {
	Key      string  `yaml:"key"`
	MinScore float64 `yaml:"minScore"`
	Label    string  `yaml:"label"`
	Display  string  `yaml:"display"`
	Color    string  `yaml:"color"`
}

type ShoeWearState struct {
	Key     string
	Label   string
	Display string
	Color   string
}

type ShoeWearResolver struct {
	standardLifeHours float64
	restWearPerDay    float64
	states            []ShoeWearStateConfig
}

func NewShoeWearResolver(cfg ShoeWearConfig) (*ShoeWearResolver, error) {
	cfg = applyShoeWearDefaults(cfg)
	if cfg.StandardLifeHours <= 0 {
		return nil, fmt.Errorf("shoeWear.standardLifeHours must be > 0")
	}
	if cfg.RestWearPerDay < 0 {
		return nil, fmt.Errorf("shoeWear.restWearPerDay must be >= 0")
	}
	if len(cfg.States) == 0 {
		return nil, fmt.Errorf("shoeWear.states is required")
	}

	seenKeys := make(map[string]struct{}, len(cfg.States))
	seenScores := make(map[float64]struct{}, len(cfg.States))
	states := make([]ShoeWearStateConfig, 0, len(cfg.States))
	for _, state := range cfg.States {
		if state.Key == "" {
			return nil, fmt.Errorf("shoeWear.states[].key is required")
		}
		if _, exists := seenKeys[state.Key]; exists {
			return nil, fmt.Errorf("duplicate shoeWear state key: %q", state.Key)
		}
		seenKeys[state.Key] = struct{}{}
		if state.MinScore < 0 || state.MinScore > 100 {
			return nil, fmt.Errorf("shoeWear state %q minScore must be in [0, 100]", state.Key)
		}
		if _, exists := seenScores[state.MinScore]; exists {
			return nil, fmt.Errorf("duplicate shoeWear state minScore: %v", state.MinScore)
		}
		seenScores[state.MinScore] = struct{}{}
		if state.Display == "" {
			return nil, fmt.Errorf("shoeWear state %q display is required", state.Key)
		}
		states = append(states, state)
	}

	sort.Slice(states, func(i, j int) bool {
		return states[i].MinScore > states[j].MinScore
	})
	if states[len(states)-1].MinScore != 0 {
		return nil, fmt.Errorf("shoeWear.states must include a fallback state with minScore 0")
	}

	return &ShoeWearResolver{
		standardLifeHours: cfg.StandardLifeHours,
		restWearPerDay:    cfg.RestWearPerDay,
		states:            states,
	}, nil
}

// StandardLifeHours 返回球鞋标准可用寿命（小时）。
func (r *ShoeWearResolver) StandardLifeHours() float64 {
	return r.standardLifeHours
}

// RestWearPerDay 返回球鞋每日静置老化系数。
func (r *ShoeWearResolver) RestWearPerDay() float64 {
	return r.restWearPerDay
}

// State 按 score 返回命中的状态配置，score 小于所有 minScore 时返回兜底状态。
func (r *ShoeWearResolver) State(score float64) ShoeWearState {
	for _, state := range r.states {
		if score >= state.MinScore {
			return toShoeWearState(state)
		}
	}
	return toShoeWearState(r.states[len(r.states)-1])
}

// RenderDisplay 将 display 模板中的 {remainingHours} 替换为实际剩余小时数。
func (r *ShoeWearResolver) RenderDisplay(display string, remainingHours int) string {
	return strings.ReplaceAll(display, "{remainingHours}", strconv.Itoa(remainingHours))
}

func toShoeWearState(state ShoeWearStateConfig) ShoeWearState {
	return ShoeWearState{
		Key:     state.Key,
		Label:   state.Label,
		Display: state.Display,
		Color:   state.Color,
	}
}

func applyShoeWearDefaults(cfg ShoeWearConfig) ShoeWearConfig {
	if len(cfg.States) == 0 {
		cfg.StandardLifeHours = 60
		cfg.RestWearPerDay = 0.06
		cfg.States = []ShoeWearStateConfig{
			{Key: "fresh", MinScore: 85, Label: "全新", Display: "全新 · 缓震充足", Color: "green"},
			{Key: "good", MinScore: 65, Label: "状态良好", Display: "状态良好 · 预计还可打 {remainingHours}h", Color: "teal"},
			{Key: "worn", MinScore: 40, Label: "开始磨损", Display: "开始磨损 · 建议关注中底", Color: "orange"},
			{Key: "tired", MinScore: 20, Label: "磨损明显", Display: "磨损明显 · 建议考虑换鞋", Color: "deep_orange"},
			{Key: "dead", MinScore: 8, Label: "缓震衰减", Display: "缓震衰减 · 建议尽快换鞋", Color: "red"},
			{Key: "expired", MinScore: 0, Label: "已超期", Display: "已超期 · 不建议继续使用", Color: "red"},
		}
	}
	return cfg
}
