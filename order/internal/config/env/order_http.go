package env

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type OrderHTTPConfig struct {
	host        string
	port        int
	readTimeout time.Duration
}

func LoadOrderHTTPConfig() (*OrderHTTPConfig, error) {
	host := os.Getenv("HTTP_HOST")
	if host == "" {
		host = "localhost"
	}

	portStr := os.Getenv("HTTP_PORT")
	if portStr == "" {
		portStr = "8080"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP_PORT: %w", err)
	}

	readTimeoutStr := os.Getenv("HTTP_READ_TIMEOUT")
	readTimeout := 5 * time.Second
	if readTimeoutStr != "" {
		if d, err := time.ParseDuration(readTimeoutStr); err == nil {
			readTimeout = d
		}
	}

	return &OrderHTTPConfig{
		host:        host,
		port:        port,
		readTimeout: readTimeout,
	}, nil
}

func (c *OrderHTTPConfig) Address() string            { return fmt.Sprintf("%s:%d", c.host, c.port) }
func (c *OrderHTTPConfig) Host() string               { return c.host }
func (c *OrderHTTPConfig) Port() int                  { return c.port }
func (c *OrderHTTPConfig) ReadTimeout() time.Duration { return c.readTimeout }