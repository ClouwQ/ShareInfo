package database

import (
	"ShareInfo/internal/domain"
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func setupMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return sqlx.NewDb(db, "sqlmock"), mock
}

func TestUploadedFileRepository_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewUploadedFileRepository(db)

	file := &domain.UploadedFile{
		Name:      "file.txt",
		Size:      1234,
		Type:      "text/plain",
		CreatedAt: time.Now(),
	}

	mock.ExpectExec("INSERT INTO uploaded_files").
		WithArgs(int64(1), file.Name, file.Size, file.Type, file.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create(context.Background(), file, 1)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUploadedFileRepository_GetAllByLinkId(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewUploadedFileRepository(db)

	rows := sqlmock.NewRows([]string{"link_id", "name", "size", "type", "created_at"}).
		AddRow(1, "file1.txt", 100, "text/plain", time.Now()).
		AddRow(1, "file2.txt", 200, "text/plain", time.Now())

	mock.ExpectQuery("SELECT \\* FROM uploaded_files WHERE link_id = \\$1").
		WithArgs(1).
		WillReturnRows(rows)

	files, err := repo.GetAllByLinkId(context.Background(), 1)
	require.NoError(t, err)
	require.Len(t, *files, 2)
}

func TestUploadedFileRepository_DeleteAllByLinkId(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewUploadedFileRepository(db)

	mock.ExpectExec("DELETE FROM uploaded_files WHERE link_id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 2))

	err := repo.DeleteAllByLinkId(context.Background(), 1)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
