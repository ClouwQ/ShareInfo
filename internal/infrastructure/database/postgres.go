package database

import (
	"ShareInfo/internal/config"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

func NewPostgresConnection(cfg *config.DatabaseConfig, logger *zap.Logger) (*sqlx.DB, error) {
	logger.Info("Connecting to PostgreSQL",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("user", cfg.User),
		zap.String("database", cfg.Password))

	dsn := cfg.DSN()
	logger.Info(fmt.Sprintf("Database URL: %s", dsn))

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Postgres: %w", err)
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxLifetime)

	logger.Debug("Successfully setting PostgreSQL")

	// healthcheck
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		err = db.Close()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}
	return db, nil
}

func Close(db *sqlx.DB, logger *zap.Logger) error {
	logger.Info("Closing PostgreSQL")
	if err := db.Close(); err != nil {
		logger.Error("Failed to close PostgreSQL", zap.Error(err))
		return err
	}
	logger.Info("Successfully closed PostgreSQL")
	return nil
}
