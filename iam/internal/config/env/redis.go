package env

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type RedisConfig struct {
	host              string
	port              int
	connectionTimeout time.Duration
	maxIdle           int
	idleTimeout       time.Duration
}

func LoadRedisConfig() (*RedisConfig, error) {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}

	portStr := os.Getenv("REDIS_PORT")
	if portStr == "" {
		portStr = "6333"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_PORT: %w", err)
	}

	connTimeoutStr := os.Getenv("REDIS_CONNECTION_TIMEOUT")
	connTimeout := 10 * time.Second
	if connTimeoutStr != "" {
		if d, err := time.ParseDuration(connTimeoutStr); err == nil {
			connTimeout = d
		}
	}

	maxIdleStr := os.Getenv("REDIS_MAX_IDLE")
	maxIdle := 10
	if maxIdleStr != "" {
		if v, err := strconv.Atoi(maxIdleStr); err == nil {
			maxIdle = v
		}
	}

	idleTimeoutStr := os.Getenv("REDIS_IDLE_TIMEOUT")
	idleTimeout := 10 * time.Second
	if idleTimeoutStr != "" {
		if d, err := time.ParseDuration(idleTimeoutStr); err == nil {
			idleTimeout = d
		}
	}

	return &RedisConfig{
		host:              host,
		port:              port,
		connectionTimeout: connTimeout,
		maxIdle:           maxIdle,
		idleTimeout:       idleTimeout,
	}, nil
}

func (c *RedisConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.host, c.port)
}
func (c *RedisConfig) ConnectionTimeout() time.Duration { return c.connectionTimeout }
func (c *RedisConfig) MaxIdle() int                     { return c.maxIdle }
func (c *RedisConfig) IdleTimeout() time.Duration       { return c.idleTimeout }