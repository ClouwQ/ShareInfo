package database

import (
	"ShareInfo/internal/domain"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type UploadedFileRepository struct {
	db *sqlx.DB
}

func NewUploadedFileRepository(db *sqlx.DB) *UploadedFileRepository {
	return &UploadedFileRepository{db: db}
}

func (r *UploadedFileRepository) Create(ctx context.Context, file *domain.UploadedFiles, linkId int64) error {
	query := `
		INSERT INTO uploaded_files (link_id, name, size, type, created_at) VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, query, linkId, file.Name, file.Size, file.Type, file.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert uploaded file: %w", err)
	}
	return nil
}

func (r *UploadedFileRepository) GetAllByLinkId(ctx context.Context, linkId int64) (*[]domain.UploadedFiles, error) {
	var files []domain.UploadedFiles
	query := `SELECT * FROM uploaded_files WHERE link_id = $1`
	err := r.db.SelectContext(ctx, &files, query, linkId)
	if err != nil {
		return nil, fmt.Errorf("failed to query uploaded files: %w", err)
	}
	return &files, nil
}

func (r *UploadedFileRepository) DeleteAllByLinkId(ctx context.Context, linkId int64) error {
	query := `DELETE FROM uploaded_files WHERE link_id = $1`
	_, err := r.db.ExecContext(ctx, query, linkId)
	if err != nil {
		return fmt.Errorf("failed to delete uploaded files: %w", err)
	}
	return nil
}
