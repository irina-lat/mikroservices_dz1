package migrator

import (
	"database/sql"
	"log"
	"path/filepath"

	"github.com/pressly/goose/v3"
)

// Run выполняет все миграции из папки migrations
func Run(db *sql.DB) error {
	log.Println("🚀 Starting migrations...")

	// Определяем путь к папке с миграциями
	migrationsDir, err := filepath.Abs("migrations")
	if err != nil {
		log.Printf("❌ Failed to get migrations path: %v", err)
		return err
	}

	log.Printf("📁 Migrations dir: %s", migrationsDir)

	// Применяем миграции
	if err := goose.Up(db, migrationsDir); err != nil {
		log.Printf("❌ Failed to apply migrations: %v", err)
		return err
	}

	log.Println("✅ Migrations completed successfully!")
	return nil
}