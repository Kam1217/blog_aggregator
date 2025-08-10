package internal

import (
	"github.com/Kam1217/blog_aggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}
