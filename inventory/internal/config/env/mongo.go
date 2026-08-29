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
	// Используем переменные с префиксом INVENTORY_
	host := os.Getenv("INVENTORY_MONGO_HOST")
	if host == "" {
		host = "localhost"
	}

	portStr := os.Getenv("INVENTORY_MONGO_PORT")
	if portStr == "" {
		portStr = "27018"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid INVENTORY_MONGO_PORT: %w", err)
	}

	database := os.Getenv("INVENTORY_MONGO_INITDB_DATABASE")
	if database == "" {
		database = "inventory"
	}

	authDB := os.Getenv("INVENTORY_MONGO_AUTH_DB")
	if authDB == "" {
		authDB = "admin"
	}

	username := os.Getenv("INVENTORY_MONGO_INITDB_ROOT_USERNAME")
	if username == "" {
		username = "admin"
	}

	password := os.Getenv("INVENTORY_MONGO_INITDB_ROOT_PASSWORD")
	if password == "" {
		password = "admin"
	}

	// Отладка: выводим значения
	fmt.Printf("🔍 MONGO_HOST: %s\n", host)
	fmt.Printf("🔍 MONGO_PORT: %d\n", port)
	fmt.Printf("🔍 MONGO_DATABASE: %s\n", database)
	fmt.Printf("🔍 MONGO_AUTH_DB: %s\n", authDB)
	fmt.Printf("🔍 MONGO_USERNAME: %s\n", username)
	fmt.Printf("🔍 MONGO_PASSWORD: %s\n", password)

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