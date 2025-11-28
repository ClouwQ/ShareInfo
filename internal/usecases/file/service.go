package file

import (
	"ShareInfo/internal/domain"
	"ShareInfo/internal/repository/database"
	"context"
	"fmt"
	"go.uber.org/zap"
)

type Service struct {
	LinkRepo         database.LinkRepository
	UploadedFileRepo database.UploadedFileRepository
	logger           *zap.Logger
}

func NewService(logger *zap.Logger) *Service {
	return &Service{
		logger: logger,
	}
}

// Create создаем файл
func (s Service) Create(ctx context.Context, file domain.UploadedFile) error {
	err := s.UploadedFileRepo.Create(ctx, &file, file.LinkId)
	if err != nil {
		return fmt.Errorf("failed to save files metadata: %w", err)
	}
	return nil
}
