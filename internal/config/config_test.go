package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestReadConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test_gatorconfig.json")
	cfgMgr := &ConfigManager{Path: configPath}

	expectedConf := &Config{DbURL: "postgres://example", CurrentUserName: "example_user_name"}
	data, err := json.Marshal(expectedConf)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("failed to write: %v", err)
	}
	t.Run("valid config file", func(t *testing.T) {
		conf, err := cfgMgr.Read()
		if err != nil {
			t.Fatalf("failed to read config: %v", err)
		}
		if conf.DbURL != expectedConf.DbURL || conf.CurrentUserName != expectedConf.CurrentUserName {
			t.Errorf("Read got %+v, want %+v", conf, expectedConf)
		}
	})
	t.Run("missing config file", func(t *testing.T) {
		missingMgr := &ConfigManager{Path: filepath.Join(tempDir, "nonexistent.json")}
		_, err := missingMgr.Read()
		if err == nil {
			t.Errorf("expected error for missing config file, got nil")
		}
	})
}
