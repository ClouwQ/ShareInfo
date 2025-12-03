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

// DeleteFile удаляет файл из ссылки: с диска, с бд, с кэша + проверяет не заморожена ли ссылка
func (s Service) DeleteFile(ctx context.Context, linkId int64, fileName string) error {
	// Проверяем доступность к удалению
	link, err := s.LinkRepo.GetByID(ctx, linkId)
	if err != nil {
		return err
	}

	if !link.IsActive || link.IsFrozen {
		return fmt.Errorf("link is not active or already frozen")
	}

	// Удаляем файл из бд
	err = s.UploadedFileRepo.DeleteOneByName(ctx, fileName)
	if err != nil {
		return err
	}

	// Удаляем файл с диска
	err = deleteFile(ctx, linkId, fileName)
	if err != nil {
		return err
	}

	// TODO Удалить файл из кэша
	return nil
}

func (s Service) GetFiles(ctx context.Context, linkId int64) error {
	return nil
}

func (s Service) DeleteFiles(ctx context.Context, file domain.UploadedFile) error {
	return nil
}
