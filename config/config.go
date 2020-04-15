package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

type Config struct {
	Listen   string `env:"LISTEN" envDefault:"localhost:8080"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"debug"`
}

// LoadConfig return struct config
func LoadConfig() Config {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	return cfg
}
