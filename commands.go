package main

import (
	"fmt"

	"github.com/Kam1217/blog_aggregator/internal/config"
)

type state struct {
	cfg        *config.Config
	cfgManager *config.ConfigManager
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("login handler expects a single argument but got an empty slice")
	}
	if err := s.cfgManager.SetUser(s.cfg, cmd.args[0]); err != nil {
		return fmt.Errorf("error setting the username to config: %w", err)
	}
	fmt.Printf("Username has been set to: %s\n", s.cfg.CurrentUserName)
	return nil
}
