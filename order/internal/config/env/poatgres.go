package env

import (
	"fmt"
	"os"
	"strconv"
)

type PostgresConfig struct {
	host     string
	port     int
	user     string
	password string
	database string
	sslMode  string
}

func LoadPostgresConfig() (*PostgresConfig, error) {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}

	portStr := os.Getenv("POSTGRES_PORT")
	if portStr == "" {
		portStr = "5435"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid POSTGRES_PORT: %w", err)
	}

	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		user = "order_user"
	}

	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		password = "order_password"
	}

	database := os.Getenv("POSTGRES_DB")
	if database == "" {
		database = "order"
	}

	sslMode := os.Getenv("POSTGRES_SSL_MODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	return &PostgresConfig{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		database: database,
		sslMode:  sslMode,
	}, nil
}

func (c *PostgresConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.host, c.port, c.user, c.password, c.database, c.sslMode)
}
func (c *PostgresConfig) Host() string     { return c.host }
func (c *PostgresConfig) Port() int        { return c.port }
func (c *PostgresConfig) User() string     { return c.user }
func (c *PostgresConfig) Password() string { return c.password }
func (c *PostgresConfig) Database() string { return c.database }
func (c *PostgresConfig) SSLMode() string  { return c.sslMode }