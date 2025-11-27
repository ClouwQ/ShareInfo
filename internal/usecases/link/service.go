package link

import (
	"ShareInfo/internal/domain"
	"ShareInfo/internal/repository/database"
	_ "ShareInfo/internal/utils"
	"context"
	"go.uber.org/zap"
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

func (s Service) CreateLink(ctx context.Context, link domain.Link) error {
	// Генерируем, сохраняем, возвращаем, запускаем ttl на удаление, если файлов не было добавлено
	//linkId, err := utils.GenerateLinkId(&s.LinkRepo)
	//if err != nil {
	//	return fmt.Errorf("error to create link id: %w", err)
	//}
	//link := domain.Link{
	//	linkId: linkId,
	//
	//}

}

func (s Service) GetLinkZipFiles(ctx context.Context, linkId domain.Link) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) DeleteLink(ctx context.Context, linkId int64) error {
	//TODO implement me
	panic("implement me")
}
