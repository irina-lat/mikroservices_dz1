package env

import (
	"fmt"
	"os"
	"strconv"
)

type InventoryGRPCConfig struct {
	host string
	port int
}

func LoadInventoryGRPCConfig() (*InventoryGRPCConfig, error) {
	host := os.Getenv("INVENTORY_GRPC_HOST")
	if host == "" {
		host = "localhost"
	}

	portStr := os.Getenv("INVENTORY_GRPC_PORT")
	if portStr == "" {
		portStr = "50051"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid INVENTORY_GRPC_PORT: %w", err)
	}

	return &InventoryGRPCConfig{
		host: host,
		port: port,
	}, nil
}

func (c *InventoryGRPCConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.host, c.port)
}

func (c *InventoryGRPCConfig) Host() string {
	return c.host
}

func (c *InventoryGRPCConfig) Port() int {
	return c.port
}