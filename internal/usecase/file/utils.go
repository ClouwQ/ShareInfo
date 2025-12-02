package file

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func getFileType(filename string) string {
	return filepath.Ext(filename)
}

func getDirPath(linkId int64) (string, error) {
	rootDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(rootDir, "files", strconv.FormatInt(linkId, 10)), nil
}

func getFilePath(linkId int64, filename string) (string, error) {
	dirPath, err := getDirPath(linkId)
	if err != nil {
		return "", err
	}
	return filepath.Join(dirPath, filename), nil
}

// makeLinkDir проверяет и создает директорию ссылки с адресом files/{linkId}
func makeLinkDir(linkId int64) (string, bool, error) {
	// Путь к новой папке
	dirPath, err := getDirPath(linkId)
	if err != nil {
		return "", false, err
	}
	info, err := os.Stat(dirPath)

	// Директория существует
	if err == nil && info.IsDir() {
		return dirPath, true, nil
	}

	if err != nil && !os.IsNotExist(err) {
		// Ошибка проверки существования
		return "", false, fmt.Errorf("error checking directory: %w", err)
	}

	// Директория не существует, создаём
	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		return "", false, fmt.Errorf("failed to create directory: %w", err)
	}

	return dirPath, false, nil
}

// saveFileOnDisk сохраняет файл в директорию
func saveFileOnDisk(ctx context.Context, file File) error {
	// Проверяем/создаем папку
	dir, _, err := makeLinkDir(file.LinkID)
	if err != nil {
		return err
	}
	// Формируем полный путь до файла
	filePath := filepath.Join(dir, file.Name)

	// Сохраняем на диск
	if err := os.WriteFile(filePath, file.file, 0644); err != nil {
		return fmt.Errorf("failed to save file %s: %w", file.Name, err)
	}
	return nil
}
