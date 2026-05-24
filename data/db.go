package data

import (
	"log"
	"os"
	"path/filepath"

	"task-manager/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

const defaultDBPath = "data/task-manager.db"

func SetupDB() *gorm.DB {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		log.Fatalf("failed to create db directory: %v", err)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&models.Assignment{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	return db
}
