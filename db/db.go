package db

import (
	"log"
	"os"
	"path/filepath"
	"strings"

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
	dsn = sqliteWritableDSN(dsn)
	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	return err
}

// sqliteWritableDSN forces read-write mode for file paths. Avoids rare "readonly database"
// cases with relative paths and makes journal files land next to the DB file predictably.
func sqliteWritableDSN(dsn string) string {
	if dsn == ":memory:" || strings.HasPrefix(dsn, "file:") {
		return dsn
	}
	abs, err := filepath.Abs(dsn)
	if err != nil {
		return dsn
	}
	// mode=rwc: open for reading and writing, create if missing
	// journal_mode=DELETE: avoid -wal/-shm sidecars (simpler permissions story)
	return "file:" + filepath.ToSlash(abs) + "?mode=rwc&_pragma=busy_timeout(5000)&_pragma=journal_mode(DELETE)"
}

// Migrate runs AutoMigrate for all domain models.
func Migrate() error {
	return DB.AutoMigrate(
		&models.Account{},
		&models.Transaction{},
		&models.InstallmentPlan{},
	)
}
