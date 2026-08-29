package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"notification/internal/app"
	"notification/internal/config"
)

func main() {
	// Загружаем .env файл
	if err := godotenv.Load("../deploy/env/.env"); err != nil {
		log.Println("⚠️ No .env file found, using system env")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	app, err := app.New(cfg)
	if err != nil {
		log.Fatalf("create app: %v", err)
	}

	if err := app.Run(context.Background()); err != nil {
		log.Printf("app error: %v", err)
		os.Exit(1)
	}
}