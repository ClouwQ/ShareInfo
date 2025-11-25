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
		err := os.Setenv("ENVIRONMENT", *envFlag)
		if err != nil {
		}
	}
	return nil
}

func NewConfig() (*Config, error) {
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

	if c.Database.Host == "" {
		missing = append(missing, "DB_HOST")
	}
	notAllowedPorts := []int{21, 23, 135, 139, 445, 3389, 1433, 3306, 5432, 5900, 6379}
	if c.Database.Port == 0 {
		for port := range notAllowedPorts {
			if c.Database.Port == port {
				return fmt.Errorf("not allowed port %d", c.Database.Port)
			}
		}
		missing = append(missing, "DB_PORT")
	}
	if c.Database.User == "" {
		missing = append(missing, "DB_USER")
	}
	if c.Database.Password == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if c.Database.DBName == "" {
		missing = append(missing, "DBNAME")
	}
	if c.Redis.Host == "" {
		missing = append(missing, "REDIS_HOST")
	}
	if c.Server.Port == 0 {
		missing = append(missing, "SERVER_PORT")
	}
	// Проверь другие "обязательные" параметры по аналогии

	if len(missing) > 0 {
		return fmt.Errorf("missing required config vars: %v", strings.Join(missing, ", "))
	}
	return nil
}
