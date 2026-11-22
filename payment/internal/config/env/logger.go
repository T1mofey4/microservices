package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type loggerEnvConfig struct {
	Level  string `env:"LOGGER_LEVEL,required"`
	AsJSON bool   `env:"LOGGER_AS_JSON,required"`
}

type loggerConfig struct {
	raw loggerEnvConfig
}

func NewLoggerConfig() (*loggerConfig, error) {
	var raw loggerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("logger parse config failed: %w", err)
	}

	return &loggerConfig{raw: raw}, nil
}

func (c *loggerConfig) Level() string {
	return c.raw.Level
}

func (c *loggerConfig) AsJSON() bool {
	return c.raw.AsJSON
}
