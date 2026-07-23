package env

import (
	"fmt"
	"os"
	"strconv"
)

type PaymentGRPCConfig struct {
	host string
	port int
}

func LoadPaymentGRPCConfig() (*PaymentGRPCConfig, error) {
	host := os.Getenv("GRPC_HOST")
	if host == "" {
		host = "localhost"
	}

	portStr := os.Getenv("GRPC_PORT")
	if portStr == "" {
		portStr = "50052"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid GRPC_PORT: %w", err)
	}

	return &PaymentGRPCConfig{
		host: host,
		port: port,
	}, nil
}

func (c *PaymentGRPCConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.host, c.port)
}
func (c *PaymentGRPCConfig) Host() string { return c.host }
func (c *PaymentGRPCConfig) Port() int    { return c.port }