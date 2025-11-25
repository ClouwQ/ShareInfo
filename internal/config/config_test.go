package config

import (
	"os"
	"testing"
)

func TestLoad_WithEnvVariables(t *testing.T) {
	// Сохраним старые значения, чтобы потом восстановить
	origEnv := map[string]string{
		"ENVIRONMENT": os.Getenv("ENVIRONMENT"),
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_PORT":     os.Getenv("DB_PORT"),
		"DB_USER":     os.Getenv("DB_USER"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
		"DBNAME":      os.Getenv("DBNAME"),
		"REDIS_HOST":  os.Getenv("REDIS_HOST"),
		"SERVER_PORT": os.Getenv("SERVER_PORT"),
	}
	defer func() {
		for k, v := range origEnv {
			os.Setenv(k, v)
		}
	}()

	// Установим моки окружения
	os.Setenv("ENVIRONMENT", "test")
	os.Setenv("DB_HOST", "mockdb")
	os.Setenv("DB_PORT", "2222")
	os.Setenv("DB_USER", "mockuser")
	os.Setenv("DB_PASSWORD", "mockpass")
	os.Setenv("DBNAME", "mockdbtest")
	os.Setenv("REDIS_HOST", "mockredis")
	os.Setenv("SERVER_PORT", "12345")

	cfg, err := NewConfig()
	if err != nil {
		t.Fatalf("unexpected error loading config: %v", err)
	}
	if cfg.Environment != "test" {
		t.Errorf("expected ENVIRONMENT=test, got %s", cfg.Environment)
	}
	if cfg.Database.Host != "mockdb" {
		t.Errorf("expected DB_HOST=mockdb, got %s", cfg.Database.Host)
	}
	if cfg.Database.Port != 2222 {
		t.Errorf("expected DB_PORT=2222, got %d", cfg.Database.Port)
	}
	if cfg.Database.User != "mockuser" {
		t.Errorf("expected DB_USER=mockuser, got %s", cfg.Database.User)
	}
	if cfg.Database.Password != "mockpass" {
		t.Errorf("expected DB_PASSWORD=mockpass, got %s", cfg.Database.Password)
	}
	if cfg.Database.DBName != "mockdbtest" {
		t.Errorf("expected DBNAME=mockdbtest, got %s", cfg.Database.DBName)
	}
	if cfg.Redis.Host != "mockredis" {
		t.Errorf("expected REDIS_HOST=mockredis, got %s", cfg.Redis.Host)
	}
	if cfg.Server.Port != 12345 {
		t.Errorf("expected SERVER_PORT=12345, got %d", cfg.Server.Port)
	}
}
