package usecase

import (
	"ShareInfo/internal/domain"
	"ShareInfo/internal/usecase/file"
	"context"
)

type Link interface {
	Create(ctx context.Context) (int64, error)
	Freeze(ctx context.Context, link domain.Link) error
	GetMeta(ctx context.Context, linkId int64) (domain.LinkWithFilesMeta, error)
	GetLinkZipFiles(ctx context.Context, linkId int64) (string, error)
	DeleteLink(ctx context.Context, linkId int64) error
}

type UploadedFile interface {
	SaveFile(ctx context.Context, file file.File) error
	DeleteFile(ctx context.Context, linkId int64, fileName string) error
	GetFiles(ctx context.Context, linkId int64) error
	DeleteFiles(ctx context.Context, file domain.UploadedFile) error
	DownloadFiles(ctx context.Context, linkId int64) (string, error)
}

type LinkMessage interface {
}
