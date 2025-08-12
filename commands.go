package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Kam1217/blog_aggregator/internal/config"
	"github.com/Kam1217/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

type state struct {
	cfg        *config.Config
	cfgManager *config.ConfigManager
	db         *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
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

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.registeredCommands[cmd.name]
	if !exists {
		return fmt.Errorf("command does not exist: %v", cmd.name)
	}
	if err := handler(s, cmd); err != nil {
		return fmt.Errorf("error calling the command: %w", err)
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("login handler expects a single argument but got an empty slice")
	}
	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	//TODO: double check we dont need err != nil
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if user.Name != "" {
		os.Exit(1)
	}

	newUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	if err := s.cfgManager.SetUser(s.cfg, newUser.Name); err != nil {
		return fmt.Errorf("error setting the config username: %w", err)
	}
	fmt.Printf("Username %s has been registered:\n", newUser.Name)

	return nil
}
