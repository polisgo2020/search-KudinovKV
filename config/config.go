package config

import (
	"os"
)

type Config struct {
	Listen   string `env:"LISTEN" envDefault:"localhost:8080"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"debug"`
}

// LoadConfig return struct config
func LoadConfig() Config {
	return Config{
		Listen:   os.Getenv("LISTEN"),
		LogLevel: os.Getenv("LOG_LEVEL"),
	}
}
