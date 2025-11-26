package database

import (
	"ShareInfo/internal/config"
	"ShareInfo/internal/domain"
	"ShareInfo/internal/infrastructure/database"
	_ "ShareInfo/internal/infrastructure/database"
	"context"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestLinkAdd(t *testing.T) {
	link := &domain.Link{
		Description:    "Description",
		IsActive:       true,
		IsOnceDownload: false,
		ExpiresAt:      time.Now().Add(time.Hour),
		CreatedAt:      time.Now(),
	}

	cfg := &config.DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "development",
		DBName:   "shareinfo",
		Password: "21421412412421412",
		SSLMode:  "disable",
	}
	logger, _ := zap.NewDevelopment()
	db, err := database.NewPostgresConnection(cfg, logger)
	logger.Info("connect to postgres")
	if err != nil {
		t.Errorf("Error connecting to database: %s", err)
	}
	logger.Info("successful connect to postgres")
	r := NewLinkRepository(db)

	err = r.Create(context.Background(), link)
	if err != nil {
		t.Errorf("failed to create link: %v", err)
	}

	if link.ID == 0 {
		t.Errorf("link id should be set")
	}
}
