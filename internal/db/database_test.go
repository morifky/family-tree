package db

import (
	"os"
	"testing"
)

func TestMustOpen(t *testing.T) {
	// Create a temporary file for the database
	tmpFile, err := os.CreateTemp("", "brayat_test_*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	dbPath := tmpFile.Name()
	tmpFile.Close()

	defer os.Remove(dbPath)

	// This should not panic
	db := MustOpen(dbPath)

	if db == nil {
		t.Error("Expected valid gorm.DB pointer, got nil")
	}

	// Verify PRAGMAs via a query
	var journalMode string
	err = db.Raw("PRAGMA journal_mode").Scan(&journalMode).Error
	if err != nil {
		t.Fatalf("failed to query journal_mode: %v", err)
	}
	// Note: It might be 'wal' (lowercase) or 'WAL'
	if journalMode != "wal" && journalMode != "WAL" {
		t.Errorf("Expected journal_mode = wal, got %s", journalMode)
	}
}

func TestMustOpen_PanicOnInvalidPath(t *testing.T) {
	// Trying to open a DB in a path that doesn't exist to cause panic
	invalidPath := "/path/that/does/not/exist/db.sqlite"

	panicked := false
	defer func() {
		if r := recover(); r != nil {
			panicked = true
		}
	}()

	MustOpen(invalidPath)

	if !panicked {
		t.Error("Expected MustOpen to panic with invalid path, but it didn't")
	}
}
