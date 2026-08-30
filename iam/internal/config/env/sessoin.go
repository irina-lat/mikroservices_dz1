package env

import (
	"os"
	"time"
)

type SessionConfig struct {
	ttl time.Duration
}

func LoadSessionConfig() (*SessionConfig, error) {
	ttlStr := os.Getenv("SESSION_TTL")
	ttl := 24 * time.Hour
	if ttlStr != "" {
		if d, err := time.ParseDuration(ttlStr); err == nil {
			ttl = d
		}
	}

	return &SessionConfig{
		ttl: ttl,
	}, nil
}

func (c *SessionConfig) TTL() time.Duration { return c.ttl }