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

func ConnectDatabase() {
	dsn := os.Getenv("DATABASE_PATH")
	if dsn == "" {
		dsn = "data/app.db"
	}
	if err := Connect(dsn); err != nil {
		log.Fatal("failed to connect database: ", err)
	}
	log.Println("Database connected")
}

func Connect(dsn string) error {
	if dsn != ":memory:" {
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

func Migrate() error {
	return DB.AutoMigrate(
		&models.Account{},
		&models.Transaction{},
		&models.InstallmentPlan{},
	)
}
