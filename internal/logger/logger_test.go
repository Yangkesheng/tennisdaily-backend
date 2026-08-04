package logger

import "testing"

func TestSetLevel(t *testing.T) {
	cases := []struct {
		level string
		want  Level
	}{
		{"debug", LevelDebug},
		{"DEBUG", LevelDebug},
		{"info", LevelInfo},
		{"warn", LevelWarn},
		{"warning", LevelWarn},
		{"error", LevelError},
	}
	for _, tc := range cases {
		if err := SetLevel(tc.level); err != nil {
			t.Fatalf("SetLevel(%q) error: %v", tc.level, err)
		}
		if currentLevel != tc.want {
			t.Fatalf("SetLevel(%q) = %d, want %d", tc.level, currentLevel, tc.want)
		}
	}
}

func TestSetLevelInvalid(t *testing.T) {
	if err := SetLevel("verbose"); err == nil {
		t.Fatal("SetLevel(verbose) should error")
	}
}

func TestEnabled(t *testing.T) {
	if err := SetLevel("info"); err != nil {
		t.Fatalf("SetLevel(info) error: %v", err)
	}
	if Enabled(LevelDebug) {
		t.Fatal("debug should be disabled at info level")
	}
	if !Enabled(LevelInfo) || !Enabled(LevelWarn) || !Enabled(LevelError) {
		t.Fatal("info/warn/error should be enabled at info level")
	}
}
