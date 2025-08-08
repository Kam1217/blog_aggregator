package main

import (
	"fmt"
	"log"

	"github.com/Kam1217/blog_aggregator/internal/config"
)

func main() {
	ptr, err := config.Read()
	if err != nil {
		log.Fatal("failed to read the config file: ", err)
	}

	if err := ptr.SetUser("Kamila"); err != nil {
		log.Fatal("failed to set user: ", err)
	}

	updated_conf, err := config.Read()
	if err != nil {
		log.Fatal("failed to read the config file: ", err)
	}

	fmt.Printf("%+v\n", updated_conf)
}
