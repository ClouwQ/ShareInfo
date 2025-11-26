package database

import (
	"ShareInfo/internal/domain"
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestLinkMessageRepository_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewLinkMessageRepository(db)

	message := &domain.LinkMessages{
		LinkID: 1,
		Text:   "test message",
		SentAt: time.Now(),
	}

	mock.ExpectExec("INSERT INTO link_messages").
		WithArgs(message.LinkID, message.Text, message.SentAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create(context.Background(), message)
	require.NoError(t, err)
	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestLinkMessageRepository_GetByLinkId(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewLinkMessageRepository(db)

	rows := sqlmock.NewRows([]string{"link_id", "text", "sent_at"}).
		AddRow(1, "msg1", time.Now()).
		AddRow(1, "msg2", time.Now())

	mock.ExpectQuery("SELECT \\* FROM link_messages WHERE link_id = \\$1").
		WithArgs(1).
		WillReturnRows(rows)

	messages, err := repo.GetByLinkId(context.Background(), 1)
	require.NoError(t, err)
	require.Len(t, messages, 2)
}

func TestLinkMessageRepository_DeleteAllByLinkId(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewLinkMessageRepository(db)

	mock.ExpectExec("DELETE FROM link_messages WHERE link_id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.DeleteAllByLinkId(context.Background(), 1)
	require.NoError(t, err)
}
