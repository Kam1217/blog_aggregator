package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestReadConfig(t *testing.T) {
	tempDir := t.TempDir()
	configFileName = filepath.Join(tempDir, "test_gatorconfig.json")
	expectedConf := &Config{DbURL: "postgres://example", CurrentUserName: "example_user_name"}
	data, err := json.Marshal(expectedConf)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	if err := os.WriteFile(configFileName, data, 0644); err != nil {
		t.Fatalf("failed to write: %v", err)
	}
	t.Run("valid config file", func(t *testing.T) {
		conf, err := Read()
		if err != nil {
			t.Fatalf("failed to read config: %v", err)
		}
		if conf.DbURL != expectedConf.DbURL || conf.CurrentUserName != expectedConf.CurrentUserName {
			t.Errorf("Read returned %s, expected %s", conf, expectedConf)
		}
	})
}
