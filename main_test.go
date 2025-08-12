package main

import (
	"errors"
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
}
