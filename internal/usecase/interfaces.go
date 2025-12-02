package usecase

import (
	"ShareInfo/internal/domain"
	"ShareInfo/internal/usecase/file"
	"context"
)

type Link interface {
	Create(ctx context.Context) (int64, error)
	Freeze(ctx context.Context, linkId int64) error

	GetLinkZipFiles(ctx context.Context, linkId domain.Link) (string, error)
	DeleteLink(ctx context.Context, linkId int64) error // Удаляются все файлы внутри него
}

type UploadedFile interface {
	SaveFile(ctx context.Context, file file.File) error

	GetFiles(ctx context.Context, linkId int64) error
	DeleteFile(ctx context.Context, fileName int64) error
	DeleteFiles(ctx context.Context, linkId int64) error
}

type LinkMessage interface {
	CreateMessage(ctx context.Context, messages domain.LinkMessages) error
	GetMessages(ctx context.Context, linkId int64) ([]domain.LinkMessages, error)
	DeleteMessages(ctx context.Context, linkId int64) error
}
