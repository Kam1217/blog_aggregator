package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Kam1217/blog_aggregator/internal/config"
	"github.com/Kam1217/blog_aggregator/internal/database"
	_ "github.com/lib/pq"
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

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatal("failed to open a connection to the database:", err)
	}
	dbQueries := database.New(db)

	programState := state{cfg: cfg, cfgManager: cfgMgr, db: dbQueries}
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)

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
