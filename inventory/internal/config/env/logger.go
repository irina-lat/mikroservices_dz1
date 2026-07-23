package env

import (
	"os"
	"strconv"
)

type LoggerConfig struct {
	level  string
	asJSON bool
}

func LoadLoggerConfig() (*LoggerConfig, error) {
	level := os.Getenv("LOGGER_LEVEL")
	if level == "" {
		level = "info"
	}

	asJSONStr := os.Getenv("LOGGER_AS_JSON")
	asJSON := true
	if asJSONStr != "" {
		asJSON, _ = strconv.ParseBool(asJSONStr)
	}

	return &LoggerConfig{
		level:  level,
		asJSON: asJSON,
	}, nil
}

func (c *LoggerConfig) Level() string {
	return c.level
}

func (c *LoggerConfig) AsJSON() bool {
	return c.asJSON
}