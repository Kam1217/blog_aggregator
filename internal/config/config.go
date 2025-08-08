package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

type ConfigManager struct {
	Path string
}

func (cm *ConfigManager) Read() (*Config, error) {
	data, err := os.ReadFile(cm.Path)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	var conf Config
	if err := json.Unmarshal(data, &conf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}
	return &conf, nil
}

func (cm *ConfigManager) SetUser(c *Config, userName string) error {
	c.CurrentUserName = userName
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal config:%w", err)
	}

	if err := os.WriteFile(cm.Path, data, 0644); err != nil {
		return fmt.Errorf("could not write config: %w", err)
	}
	return nil
}
