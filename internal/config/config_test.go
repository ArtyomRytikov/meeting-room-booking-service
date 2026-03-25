package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("APP_PORT", "8080")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "user")
	os.Setenv("DB_PASSWORD", "pass")
	os.Setenv("DB_NAME", "db")

	cfg := Load()

	if cfg.AppPort != "8080" {
		t.Fatalf("expected port 8080, got %s", cfg.AppPort)
	}
}

func TestDatabaseURL(t *testing.T) {
	cfg := Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "user",
		DBPassword: "pass",
		DBName:     "db",
	}

	url := cfg.DatabaseURL()

	if url == "" {
		t.Fatal("expected non-empty db url")
	}
}
