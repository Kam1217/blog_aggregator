package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/Kam1217/blog_aggregator/internal/config"
)

func TestCommandsRegister(t *testing.T) {
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	testHandler := func(s *state, cmd command) error {
		return nil
	}

	t.Run("Register single command", func(t *testing.T) {
		cmds.register("test", testHandler)
		if len(cmds.registeredCommands) != 1 {
			t.Errorf("Expected 1 registered command but got %d", len(cmds.registeredCommands))
		}

		if cmds.registeredCommands["test"] == nil {
			t.Errorf("Expected 'test' command to be registered")
		}
	})

	t.Run("Register multiple commands", func(t *testing.T) {
		cmds.register("test_2", testHandler)
		if len(cmds.registeredCommands) != 2 {
			t.Errorf("Expected 2 registered commands but got %d", len(cmds.registeredCommands))
		}
	})
}

func TestCommandRun(t *testing.T) {
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cfg := &config.Config{DbURL: "test_db", CurrentUserName: ""}
	cfgMgr := &config.ConfigManager{Path: "/temp/test"}
	s := &state{cfg: cfg, cfgManager: cfgMgr}

	successHandler := func(s *state, cmd command) error {
		return nil
	}

	failedHandler := func(s *state, cmd command) error {
		return errors.New("failed handler")
	}

	cmds.register("success", successHandler)
	cmds.register("fail", failedHandler)

	t.Run("No command", func(t *testing.T) {
		cmd := command{name: "nonexistent", args: []string{}}
		err := cmds.run(s, cmd)
		if err == nil {
			t.Error("Expected error for non-existent command")
		}
	})

	t.Run("Success command", func(t *testing.T) {
		cmd := command{name: "success", args: []string{}}
		err := cmds.run(s, cmd)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
	})

	t.Run("Fail command", func(t *testing.T) {
		cmd := command{name: "fail", args: []string{}}
		err := cmds.run(s, cmd)
		if err == nil {
			t.Errorf("Expected error from failer handler")
		}
	})
}

func TestHandlerLogin(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test_gatorconfig.jason")
	cfgMgr := &config.ConfigManager{Path: configPath}

	expectedConf := &config.Config{DbURL: "postgres://example", CurrentUserName: ""}
	data, err := json.Marshal(expectedConf)
	if err != nil {
		t.Errorf("Failed to marshal:%v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := cfgMgr.Read()
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	s := &state{cfg: cfg, cfgManager: cfgMgr}

	t.Run("No arguments", func(t *testing.T) {
		cmd := command{name: "login", args: []string{}}
		err := handlerLogin(s, cmd)
		if err == nil {
			t.Errorf("Expected error when no arguments provided")
		}
	})
}
