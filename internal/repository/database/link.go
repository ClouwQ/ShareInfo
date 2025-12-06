package database

import (
	"ShareInfo/internal/domain"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"
	"time"
)

type LinkRepository struct {
	db *sqlx.DB
}

func NewLinkRepository(db *sqlx.DB) *LinkRepository {
	return &LinkRepository{
		db: db,
	}
}

func (r *LinkRepository) Create(ctx context.Context, link *domain.Link) error {

	// Запрос в Link
	{
		query := `
			INSERT INTO links (description, is_active, is_once_download, expires_at, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
			`

		err := r.db.QueryRowxContext(ctx, query,
			link.Description,
			link.IsActive,
			link.IsOnceDownload,
			link.ExpiresAt,
			link.CreatedAt,
		).Scan(&link.ID)

		if err != nil {
			return fmt.Errorf("failed to create link: %w", err)
		}
	}

	// Запрос в LinkAnalytics
	{
		query := `
			INSERT INTO links_analytics (link_id, downloads, last_download_at) VALUES ($1, $2, $3)
			`
		_, err := r.db.ExecContext(ctx, query, link.ID, 0, link.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to create link: %w", err)
		}
	}
	return nil
}

func (r *LinkRepository) GetByID(ctx context.Context, id int64) (*domain.Link, error) {
	var link *domain.Link
	query := `
		SELECT * FROM links WHERE id = $1
	`
	err := r.db.QueryRowxContext(ctx, query, id).Scan(&link)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch link: %w", err)
	}
	return link, nil
}

func (r *LinkRepository) Delete(ctx context.Context, id int64) error {

	errG := errgroup.Group{}

	errG.Go(func() error {
		queryLink := `DELETE FROM links WHERE id = $1`
		_, err := r.db.ExecContext(ctx, queryLink, id)
		if err != nil {
			return fmt.Errorf("failed to delete link: %w", err)
		}
		return nil
	})

	errG.Go(func() error {
		queryLinkAnalytics := `DELETE FROM links_analytics WHERE links_analytics.link_id = $1`
		_, err := r.db.ExecContext(ctx, queryLinkAnalytics, id)
		if err != nil {
			return fmt.Errorf("failed to delete links_analytics: %w", err)
		}
		return nil
	})

	return errG.Wait()
}

// Freeze замораживает ссылку, говоря о том, что ее нельзя изменить и время удаления пошло
func (r *LinkRepository) Freeze(ctx context.Context, link domain.Link) error {
	query := `
		UPDATE links SET is_frozen = true, is_active = true, is_once_download = $1, expires_at = $2, created_at = $3 WHERE id = $4
`
	_, err := r.db.ExecContext(ctx, query, link.IsOnceDownload, link.ExpiresAt, link.CreatedAt, link.ID)
	if err != nil {
		return fmt.Errorf("failed to freeze link: %w", err)
	}
	return nil
}

func (r *LinkRepository) UpdateDescription(ctx context.Context, id int64, description string) error {
	query := `
	UPDATE links SET description = $2 WHERE id = $1
`
	_, err := r.db.ExecContext(ctx, query, id, description)
	if err != nil {
		return fmt.Errorf("failed to update links description: %w", err)
	}
	return nil
}

// NewDownload инкриминируем счетчик загрузок у аналитики и выставляем новую дату
func (r *LinkRepository) IncrementDownloadsCount(ctx context.Context, linkId int64) error {
	query := `
		UPDATE links_analytics
		SET downloads = downloads + 1, last_download_at = $1
		WHERE link_id = $2
		`
	_, err := r.db.ExecContext(ctx, query, linkId, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update links analytics: %w", err)
	}
	return nil
}
