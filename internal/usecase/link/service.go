package link

import (
	"ShareInfo/internal/domain"
	"ShareInfo/internal/repository/database"
	"ShareInfo/internal/utils"
	_ "ShareInfo/internal/utils"
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

func NewLinkUseCase(logger *zap.Logger, linkRepo database.LinkRepository, fileRepo database.UploadedFileRepository) *Service {
	return &Service{
		LinkRepo:         linkRepo,
		UploadedFileRepo: fileRepo,
		logger:           logger,
	}
}

func (s Service) CreateLink(ctx context.Context) (int64, error) {
	// Генерируем, сохраняем, возвращаем, запускаем ttl на удаление, если файлов не было добавлено
	linkId, err := utils.GenerateLinkId(&s.LinkRepo)
	if err != nil {
		return 0, fmt.Errorf("error to create link id: %w", err)

	}
	link := domain.Link{
		ID: linkId,
		// Даем пользователю 15 мин на рассуждение
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Minute * 15),
	}

	// TODO Запускаем функцию по отчистке ссылок (если в ссылке ничего не изменилось за 15 мин — удаляем)

	err = s.LinkRepo.Create(ctx, &link)
	if err != nil {
		return 0, fmt.Errorf("error to create link: %w", err)
	}
	return linkId, nil
}

func (s Service) Freeze(ctx context.Context, link domain.Link) error {
	// Закидываем в большую бд
	err := s.LinkRepo.Freeze(ctx, link)
	if err != nil {
		return fmt.Errorf("error to freeze link: %w", err)
	}

	// Кэшируем
	// TODO добавить функционал кэширования

	// TODO: откладываем удаление линки, если изменений не произойдет (не будут добавлены файлы)

	return nil
}

// GetMeta собираем все метаданные о ссылки, но без чата
func (s Service) GetMeta(ctx context.Context, linkId int64) (domain.LinkWithFilesMeta, error) {
	var meta domain.LinkWithFilesMeta
	linkMeta, err := s.LinkRepo.GetByID(ctx, linkId)
	if err != nil {
		return domain.LinkWithFilesMeta{}, fmt.Errorf("error to get link meta: %w", err)
	}
	filesMeta, err := s.UploadedFileRepo.GetAllByLinkId(ctx, linkId)
	if err != nil {
		return domain.LinkWithFilesMeta{}, fmt.Errorf("error to get files meta: %w", err)
	}

	meta.LinkMeta = *linkMeta
	meta.FilesMeta = *filesMeta
	return meta, nil
}

func (s Service) GetLinkZipFiles(ctx context.Context, linkId int64) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) DeleteLink(ctx context.Context, linkId int64) error {
	//TODO implement me
	panic("implement me")
}
