package main

import (
	"testing"
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
