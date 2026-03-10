package db

import (
	"fmt"
	"log"
	"time"

	"brayat/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MustOpen initializes the database connection, applies required PRAGMAs,
// and auto-migrates the models. It panics if the database cannot be opened.
func MustOpen(dbPath string) *gorm.DB {
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// Since we are using standard sqlite3 driver, PRAGMAs can be passed in DSN
	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_synchronous=NORMAL&_foreign_keys=ON&_busy_timeout=5000", dbPath)

	database, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(fmt.Errorf("failed to open database: %w", err))
	}

	// Ensure underlying sql.DB is accessible
	sqlDB, err := database.DB()
	if err != nil {
		panic(fmt.Errorf("failed to get sql.DB: %w", err))
	}

	// Executing PRAGMAs directly just to be 100% compliant with the spec rules
	database.Exec("PRAGMA journal_mode = WAL;")
	database.Exec("PRAGMA synchronous = NORMAL;")
	database.Exec("PRAGMA foreign_keys = ON;")
	database.Exec("PRAGMA busy_timeout = 5000;")

	// Set connection pool settings to avoid busy locks in WAL mode
	sqlDB.SetMaxOpenConns(1) // SQLite works best with 1 open conn for writes to prevent locked database errors

	// AutoMigrate all models here (will be populated in MOR-12)
	err = database.AutoMigrate(
		&model.Session{},
		&model.AccessLink{},
		&model.Person{},
		&model.Relationship{},
	)
	if err != nil {
		panic(fmt.Errorf("failed to run migrations: %w", err))
	}

	return database
}
