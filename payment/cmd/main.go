package main

import (
	"context"
	"log"
	"os"

	"payment/internal/app"
	"payment/internal/config"
)

func main() {
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