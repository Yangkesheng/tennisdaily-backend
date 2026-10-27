package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigPathPriority(t *testing.T) {
	original, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(original); err != nil {
			t.Error(err)
		}
	})
	t.Setenv("CONFIG_PATH", "")
	for _, name := range []string{"config.yaml", "config.local.example.yaml"} {
		if err := os.WriteFile(name, []byte("{}"), 0600); err != nil {
			t.Fatal(err)
		}
		if got := configPath(); got != name {
			t.Fatalf("got %s, want %s", got, name)
		}
	}
	t.Setenv("CONFIG_PATH", filepath.Join(dir, "custom.yaml"))
	if got := configPath(); got != os.Getenv("CONFIG_PATH") {
		t.Fatalf("explicit path ignored: %s", got)
	}
}

func TestLocalExampleLoads(t *testing.T) {
	t.Setenv("CONFIG_PATH", "../../config.local.example.yaml")
	cfg := loadConfigFile()
	if cfg.Database.Host != "127.0.0.1" || cfg.Storage.Enabled || cfg.Storage.RunOnStart || cfg.Storage.ShoeRunOnStart {
		t.Fatal("local example must use localhost and disable migration tasks")
	}
	if _, err := NewSessionCategoryResolver(cfg.SessionEnums); err != nil {
		t.Fatal(err)
	}
	if _, err := NewPolyesterStringHealthResolver(cfg.PolyesterStringHealth); err != nil {
		t.Fatal(err)
	}
	if _, err := NewShoeWearResolver(cfg.ShoeWear); err != nil {
		t.Fatal(err)
	}
}
