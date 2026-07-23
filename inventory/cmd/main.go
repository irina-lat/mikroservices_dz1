package main

import (
	"context"
	"log"
	"os"

	"inventory/internal/app"
	"inventory/internal/config"
)

func main() {
	// 1. Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Создаём приложение
	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	// 3. Запускаем приложение
	ctx := context.Background()
	if err := application.Run(ctx); err != nil {
		log.Printf("Application error: %v", err)
		os.Exit(1)
	}
}