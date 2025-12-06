package file

import (
	"ShareInfo/internal/domain"
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// buildZipForLink собирает файлы в zip для скачивания по ссылке
func buildZipForLink(ctx context.Context, linkId int64, files *[]domain.UploadedFile) (string, error) {
	dirPath, err := getDirPath(linkId)
	if err != nil {
		return "", err
	}
	zipPath := dirPath + ".zip"

	// Создаем файл архива
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return "", fmt.Errorf("create zip file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Добавляем каждый файл в архив
	for _, file := range *files {
		filePath := filepath.Join(dirPath, file.Name)
		if err := addFileToZip(ctx, zipWriter, filePath, file.Name); err != nil {
			return "", fmt.Errorf("add file to zip: %w", err)
		}
	}

	// Закрываем writer
	if err := zipWriter.Close(); err != nil {
		return "", fmt.Errorf("close zip writer: %w", err)
	}

	return zipPath, nil
}

// addFileToZip сохраняем файл в архив
func addFileToZip(ctx context.Context, zipWriter *zip.Writer, filePath string, fileName string) error {
	src, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer src.Close()

	// Имя внутри архива — только имя файла, без полного пути
	w, err := zipWriter.Create(fileName)
	if err != nil {
		return fmt.Errorf("create zip entry: %w", err)
	}

	if _, err := io.Copy(w, src); err != nil {
		return fmt.Errorf("copy to zip: %w", err)
	}

	return nil
}
