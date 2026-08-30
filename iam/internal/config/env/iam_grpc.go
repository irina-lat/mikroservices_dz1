package env

import (
	"fmt"
	"os"
	"strconv"
)

type IAMGRPCConfig struct {
	host string
	port int
}

func LoadIAMGRPCConfig() (*IAMGRPCConfig, error) {
	host := os.Getenv("IAM_GRPC_HOST")
	if host == "" {
		host = "localhost"
	}

	portStr := os.Getenv("IAM_GRPC_PORT")
	if portStr == "" {
		portStr = "50053"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid IAM_GRPC_PORT: %w", err)
	}

	return &IAMGRPCConfig{
		host: host,
		port: port,
	}, nil
}

func (c *IAMGRPCConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.host, c.port)
}
func (c *IAMGRPCConfig) Host() string { return c.host }
func (c *IAMGRPCConfig) Port() int    { return c.port }