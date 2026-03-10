package config

import (
	"os"
	"testing"
)

func TestMustLoad_Defaults(t *testing.T) {
	// Temporarily clear environment variables that might interfere
	os.Unsetenv("PORT")
	os.Unsetenv("DATABASE_PATH")
	os.Unsetenv("PHOTOS_DIR")
	os.Unsetenv("LOG_LEVEL")

	// This should not panic as defaults cover all required fields
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("MustLoad panicked on defaults: %v", r)
		}
	}()

	cfg := MustLoad()

	if cfg.Port != "8080" {
		t.Errorf("Expected Port to be 8080, got %s", cfg.Port)
	}
	if cfg.DatabasePath != "/data/brayat.db" {
		t.Errorf("Expected DatabasePath to be /data/brayat.db, got %s", cfg.DatabasePath)
	}
	if cfg.PhotosDir != "/data/photos" {
		t.Errorf("Expected PhotosDir to be /data/photos, got %s", cfg.PhotosDir)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("Expected LogLevel to be info, got %s", cfg.LogLevel)
	}
}

func TestMustLoad_Overrides(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("DATABASE_PATH", "/tmp/brayat.sqlite")
	os.Setenv("PHOTOS_DIR", "/tmp/photos")
	os.Setenv("LOG_LEVEL", "debug")

	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DATABASE_PATH")
		os.Unsetenv("PHOTOS_DIR")
		os.Unsetenv("LOG_LEVEL")
	}()

	cfg := MustLoad()

	if cfg.Port != "9090" {
		t.Errorf("Expected Port to be 9090, got %s", cfg.Port)
	}
	if cfg.DatabasePath != "/tmp/brayat.sqlite" {
		t.Errorf("Expected DatabasePath to be /tmp/brayat.sqlite, got %s", cfg.DatabasePath)
	}
	if cfg.PhotosDir != "/tmp/photos" {
		t.Errorf("Expected PhotosDir to be /tmp/photos, got %s", cfg.PhotosDir)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("Expected LogLevel to be debug, got %s", cfg.LogLevel)
	}
}

func TestMustLoad_PanicOnMissingRequiredOrInvalid(t *testing.T) {
	// Set an invalid port (non-numeric)
	os.Setenv("PORT", "abc")
	defer os.Unsetenv("PORT")

	panicked := false
	defer func() {
		if r := recover(); r != nil {
			panicked = true
		}
	}()

	MustLoad()

	if !panicked {
		t.Error("Expected MustLoad to panic due to invalid PORT")
	}
}
