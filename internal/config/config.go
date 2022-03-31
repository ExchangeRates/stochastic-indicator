package config

import "github.com/sirupsen/logrus"

type Config struct {
	Port         int          `toml:"port"`
	LogLevel     logrus.Level `toml:"log_level"`
	EmaClientURL string       `toml:"ema_cleint_url"`
}

func NewConfig() *Config {
	return &Config{
		Port:         8082,
		LogLevel:     logrus.DebugLevel,
		EmaClientURL: "http://localhost:8080",
	}
}
