package config

import (
	"github.com/caarlos0/env"
	zl "github.com/rs/zerolog/log"
)

type Config struct {
	Listen   string `env:"LISTEN" envDefault:"localhost:8081"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"debug"`
	PgSQL    string `env:"PGSQL" envDefault:"postgres://postgres:123456@localhost:5432/db?sslmode=disable"`
}

// LoadConfig return struct config
func LoadConfig() Config {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		zl.Fatal().Err(err).
			Msg("Can't parse env args")
	}
	return cfg
}
