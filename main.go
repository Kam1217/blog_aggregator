package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Kam1217/blog_aggregator/internal/config"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("error getting home directory: ", err)
	}
	configPath := filepath.Join(home, ".gatorconfig.json")
	cfgMgr := &config.ConfigManager{Path: configPath}

	cfg, err := cfgMgr.Read()
	if err != nil {
		log.Fatal("failed to read config: ", err)
	}

	programState := state{cfg: cfg, cfgManager: cfgMgr}
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("not enough arguments provided")
	}
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]
	cmd := command{name: cmdName, args: cmdArgs}
	if err := cmds.run(&programState, cmd); err != nil {
		log.Fatal("username is required")
	}

	updatedCfg, err := cfgMgr.Read()
	if err != nil {
		log.Fatal("failed to read updated config: ", err)
	}
	fmt.Printf("%+v\n", updatedCfg)
}
