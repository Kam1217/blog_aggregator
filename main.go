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

	if err := cfgMgr.SetUser(cfg, "Kamila"); err != nil {
		log.Fatal("failers to set user: ", err)
	}

	updatedCfg, err := cfgMgr.Read()
	if err != nil {
		log.Fatal("failed to read updated config: ", err)
	}
	fmt.Printf("%+v\n", updatedCfg)
}
