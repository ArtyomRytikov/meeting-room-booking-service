package config

import (
	"os"
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	_ = os.Setenv("APP_PORT", "8080")
	_ = os.Setenv("DB_HOST", "localhost")
	_ = os.Setenv("DB_PORT", "5432")
	_ = os.Setenv("DB_USER", "user")
	_ = os.Setenv("DB_PASSWORD", "pass")
	_ = os.Setenv("DB_NAME", "db")

	cfg := Load()

	if cfg.AppPort != "8080" {
		t.Fatalf("expected port 8080, got %s", cfg.AppPort)
	}
	if cfg.DBHost != "localhost" {
		t.Fatalf("expected host localhost, got %s", cfg.DBHost)
	}
	if cfg.DBPort != "5432" {
		t.Fatalf("expected port 5432, got %s", cfg.DBPort)
	}
	if cfg.DBUser != "user" {
		t.Fatalf("expected user, got %s", cfg.DBUser)
	}
	if cfg.DBPassword != "pass" {
		t.Fatalf("expected password pass, got %s", cfg.DBPassword)
	}
	if cfg.DBName != "db" {
		t.Fatalf("expected db name db, got %s", cfg.DBName)
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

	if !strings.Contains(url, "postgres://user:pass@localhost:5432/db") {
		t.Fatalf("unexpected db url: %s", url)
	}
}

func TestGetEnvFallback(t *testing.T) {
	_ = os.Unsetenv("SOME_UNKNOWN_ENV")

	got := getEnv("SOME_UNKNOWN_ENV", "fallback")
	if got != "fallback" {
		t.Fatalf("expected fallback, got %s", got)
	}
}
