package usecases

import (
	"ShareInfo/internal/domain"
	"context"
	"mime/multipart"
)

type Link interface {
	CreateLink(ctx context.Context, link domain.Link) error
	GetLinkZipFiles(ctx context.Context, linkId domain.Link) (string, error)
	DeleteLink(ctx context.Context, linkId int64) error // Удаляются все файлы внутри него
}

type UploadedFile interface {
	SaveFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) error
	GetFiles(ctx context.Context, linkId int64) error
	DeleteFile(ctx context.Context, fileName int64) error
	DeleteFiles(ctx context.Context, linkId int64) error
}

type LinkMessage interface {
	CreateMessage(ctx context.Context, messages domain.LinkMessages) error
	GetMessages(ctx context.Context, linkId int64) ([]domain.LinkMessages, error)
	DeleteMessages(ctx context.Context, linkId int64) error
}
