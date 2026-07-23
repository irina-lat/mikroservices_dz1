package env

import (
	"fmt"
	"os"
	"strconv"
)

type MongoConfig struct {
	host     string
	port     int
	database string
	authDB   string
	username string
	password string
}

func LoadMongoConfig() (*MongoConfig, error) {
	host := os.Getenv("MONGO_HOST")
	if host == "" {
		host = "localhost"
	}

	portStr := os.Getenv("MONGO_PORT")
if portStr == "" {
		portStr = "27018"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid MONGO_PORT: %w", err)
	}

	database := os.Getenv("MONGO_DATABASE")
	if database == "" {
		database = "inventory"
	}

	authDB := os.Getenv("MONGO_AUTH_DB")
	if authDB == "" {
		authDB = "admin"
	}

	username := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	if username == "" {
		username = "inventory_admin"
	}

	password := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	if password == "" {
		password = "inventory_secret"
	}

	return &MongoConfig{
		host:     host,
		port:     port,
		database: database,
		authDB:   authDB,
		username: username,
		password: password,
	}, nil
}

func (c *MongoConfig) URI() string {
	if c.username != "" && c.password != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=%s",
			c.username, c.password, c.host, c.port, c.database, c.authDB)
	}
	return fmt.Sprintf("mongodb://%s:%d/%s", c.host, c.port, c.database)
}

func (c *MongoConfig) Database() string {
	return c.database
}

func (c *MongoConfig) Host() string {
	return c.host
}

func (c *MongoConfig) Port() int {
	return c.port
}