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

func Read() (*Config, error) {
	home_path, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not return current home directory: %w", err)
	}
	path := home_path + "/.gatorconfig.json"

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	var conf Config
	if err := json.Unmarshal(data, &conf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}
	return &conf, nil
}

func (c *Config) SetUser(user_name string) error {
	c.CurrentUserName = user_name
	data, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return fmt.Errorf("could not marshal config:%w", err)
	}
	home_path, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get home directory:%w", err)
	}
	path := home_path + "/.gatorconfig.json"

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("could not write config: %w", err)
	}
	return nil
}
