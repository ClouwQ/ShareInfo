package file

import (
	"ShareInfo/internal/domain"
	"ShareInfo/internal/repository/database"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"time"
)

type Service struct {
	LinkRepo         database.LinkRepository
	UploadedFileRepo database.UploadedFileRepository
	logger           *zap.Logger
}

func NewService(logger *zap.Logger, linkRepo database.LinkRepository, fileRepo database.UploadedFileRepository) *Service {
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

// DownloadFiles скачать все файлы в zip
func (s Service) DownloadFiles(ctx context.Context, linkId int64) (string, error) {
	// Проверяем возможность скачивания файла
	link, err := s.LinkRepo.GetByID(ctx, linkId)
	if err != nil {
		return "", err
	}
	if !link.IsActive {
		return "", errors.New("link is not active")
	}

	// Проверяем наличие собранного zip файла
	haveZip, zipPath, err := checkZipFile(linkId)
	if err != nil {
		return "", err
	}
	// Если нет собранного zip, то создаем :)
	if !haveZip {
		// Получаем список файлов
		files, err := s.UploadedFileRepo.GetAllByLinkId(ctx, linkId)
		if err != nil {
			return "", fmt.Errorf("failed to fetch files: %w", err)
		}

		if len(*files) == 0 {
			return "", fmt.Errorf("no files found")
		}

		// Отправляем на создание
		zipPath, err = buildZipForLink(ctx, linkId, files)
		if err != nil {
			return "", fmt.Errorf("failed to build zip: %w", err)
		}
	}

	// Инкриминируем счетчик в таблице аналитики
	err = s.LinkRepo.IncrementDownloadsCount(ctx, linkId)
	if err != nil {
		return "", fmt.Errorf("failed to increment downloads count: %w", err)
	}

	return zipPath, nil
}
