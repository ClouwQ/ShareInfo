package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

func Load() error {
	envFlag := flag.String("env", "dev", "environment")
	flag.Parse()

	envFile := fmt.Sprintf(".env.%s", *envFlag)

	if _, err := os.Stat(envFile); err == nil {
		fmt.Printf("Loading enviroment from: %s\n", envFile)
		if err := godotenv.Load(envFile); err != nil {
			return fmt.Errorf("error loading %s file: %w", envFile, err)
		}
	} else {
		fmt.Printf("File %s not found, trying .env\n", envFile)
		if err := godotenv.Load(".env"); err != nil {
			return fmt.Errorf("error loading %s file: %w", envFile, err)
		}
	}

	if os.Getenv("ENVIRONMENT") == "" {
		os.Setenv("ENVIRONMENT", *envFlag)
	}
	return nil
}

func NewConfig() (*Config, error) {
	err := Load()
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	// сразу парсим файл в структуру
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("error to parse config: %w", err)
	}

	// валидируем
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("error to validate config: %w", err)
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	var missing []string

	// Database validation
	if c.Database.Host == "" {
		missing = append(missing, "DB_HOST")
	}
	if c.Database.Port == 0 {
		missing = append(missing, "DB_PORT")
	}
	if c.Database.User == "" {
		missing = append(missing, "DB_USER")
	}
	if c.Database.Password == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if c.Database.DBName == "" {
		missing = append(missing, "DB_NAME")
	}

	// Server validation
	if c.Server.Port == 0 {
		missing = append(missing, "SERVER_PORT")
	}

	// Check for blocked ports
	blockedPorts := map[int]bool{
		21: true, 23: true, 135: true, 139: true, 445: true,
		3389: true, 1433: true, 3306: true, 5432: true, 5900: true, 6379: true,
	}
	if blockedPorts[c.Server.Port] {
		return fmt.Errorf("port %d is not allowed", c.Server.Port)
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required config vars: %s", strings.Join(missing, ", "))
	}
	return nil
}
