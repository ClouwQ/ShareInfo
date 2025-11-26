package database

import (
	"ShareInfo/internal/domain"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type LinkMessageRepository struct {
	db *sqlx.DB
}

func NewLinkMessageRepository(db *sqlx.DB) *LinkMessageRepository {
	return &LinkMessageRepository{
		db: db,
	}
}

func (r *LinkMessageRepository) Create(ctx context.Context, message *domain.LinkMessages) error {
	query := `INSERT INTO link_messages (link_id, text, sent_at) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, message.LinkID, message.Text, message.SentAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *LinkMessageRepository) GetByLinkId(ctx context.Context, linkID int64) ([]*domain.LinkMessages, error) {
	var messages []*domain.LinkMessages
	query := `SELECT * FROM link_messages WHERE link_id = $1`
	err := r.db.SelectContext(ctx, &messages, query, linkID)
	if err != nil {
		return nil, fmt.Errorf("failed get link messages: %w", err)
	}
	return messages, nil
}

func (r *LinkMessageRepository) DeleteAllByLinkId(ctx context.Context, linkID int64) error {
	query := `DELETE FROM link_messages WHERE link_id = $1`
	_, err := r.db.ExecContext(ctx, query, linkID)
	if err != nil {
		return fmt.Errorf("failed delete link messages: %w", err)
	}
	return nil
}
