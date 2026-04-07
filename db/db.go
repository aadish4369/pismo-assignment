package db

import (
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"pismo-assignment/models"
)

var DB *gorm.DB

// ConnectDatabase opens SQLite using DATABASE_PATH or "data/app.db".
func ConnectDatabase() {
	dsn := os.Getenv("DATABASE_PATH")
	if dsn == "" {
		dsn = "data/app.db"
	}
	if err := ConnectWithDSN(dsn); err != nil {
		log.Fatal("failed to connect database: ", err)
	}
	log.Println("Database connected")
}

// ConnectWithDSN opens SQLite at the given DSN (e.g. file path or :memory:).
func ConnectWithDSN(dsn string) error {
	if dsn != ":memory:" && filepath.Base(dsn) != "" {
		if dir := filepath.Dir(dsn); dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return err
			}
		}
	}
	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	return err
}

// Migrate runs AutoMigrate for all domain models.
func Migrate() error {
	return DB.AutoMigrate(
		&models.Account{},
		&models.Transaction{},
		&models.Installment{},
	)
}
