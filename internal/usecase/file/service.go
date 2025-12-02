package file

import (
	"ShareInfo/internal/domain"
	"ShareInfo/internal/repository/database"
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"
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

type File struct {
	file   []byte
	LinkID int64
	Name   string
	Size   int64
}

// Create создаем файл
func (s Service) SaveFile(ctx context.Context, file File) error {
	fileObj := &domain.UploadedFile{
		LinkId:    file.LinkID,
		Name:      file.Name,
		Size:      file.Size,
		Type:      getFileType(file.Name),
		CreatedAt: time.Now(),
	}

	// Сохраняем на диск
	if err := saveFileOnDisk(ctx, file); err != nil {
		return err
	}

	// Закидываем в базу
	err := s.UploadedFileRepo.Create(ctx, fileObj, fileObj.LinkId)
	if err != nil {
		return fmt.Errorf("failed to save files metadata: %w", err)
	}

	// TODO кэшируем

	return nil
}

func (s Service) GetFiles(ctx context.Context, linkId int64) error {
	return nil
}

func (s Service) DeleteFile(ctx context.Context, linkId int64) error {
	return nil
}

func (s Service) DeleteFiles(ctx context.Context, file domain.UploadedFile) error {
	return nil
}
