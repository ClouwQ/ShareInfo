package http

import (
	usecaseFileSerice "ShareInfo/internal/usecase/file"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
)

type FileUploadMetaRequest struct {
	LinkID int64 `form:"link_id" binding:"required"`
}

type FileDeleteMetaRequest struct {
	LinkID   int64  `form:"link_id" binding:"required"`
	FileName string `form:"file_name" binding:"required"`
}

// UploadFile godoc
// @Summary Загрузить файл к ссылке
// @Description Прикрепляет файл к указанной ссылке. Поддерживает multipart/form-data
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param link_id formData int64 true "ID ссылки"
// @Param file formData file true "Файл для загрузки"
// @Success 200 {object} map[string]interface{} "Файл загружен"
// @Failure 400 {object} map[string]interface{} "Ошибка валидации или файл не найден"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /files/upload [post]
func (h *Handler) UploadFile(c *gin.Context) {
	var savedFile usecaseFileSerice.File
	var req FileUploadMetaRequest
	if err := c.ShouldBind(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "bad_request", "Failed to bind request body (need link_id)")
		return
	}

	// Забираем файл
	file, err := c.FormFile("file")
	if err != nil {
		RespondError(c, http.StatusBadRequest, "bad_request", "Failed to get file from form")
		return
	}

	// Пакуем в []byte
	f, err := file.Open()
	if err != nil {
		RespondError(c, http.StatusBadRequest, "bad_request", "Failed to get file from form")
		return
	}
	defer f.Close()

	fileByte, err := io.ReadAll(f)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "bad_request", "Failed to get file from form")
		return
	}

	// Забираем Headers файла
	savedFile.LinkID = req.LinkID
	savedFile.Name = file.Filename
	savedFile.Size = int64(len(fileByte))

	// Проверяем размер файла
	// TODO пробросить конфиг с максимальным размерам файла

	// Кидаем в usecase
	err = h.fileUseCase.SaveFile(context.Background(), savedFile)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "internal_server", "Failed to save file")
		return
	}
	RespondSuccess(c, http.StatusOK, "file_uploaded")
}

// DeleteFile godoc
// @Summary Удалить файл из ссылки
// @Description Удаляет конкретный файл по имени из ссылки
// @Tags files
// @Accept json
// @Produce json
// @Param link_id formData int64 true "ID ссылки"
// @Param file_name formData string true "Имя файла для удаления"
// @Success 200 {object} map[string]interface{} "Файл удален"
// @Failure 400 {object} map[string]interface{} "Ошибка валидации параметров"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /files/delete [delete]
func (h *Handler) DeleteFile(c *gin.Context) {
	// Забираем параметры запроса
	var fileMeta FileDeleteMetaRequest
	if err := c.ShouldBind(&fileMeta); err != nil {
		RespondError(c, http.StatusBadRequest, "bad_request", "Failed to bind request body (need link_id and file_name)")
		return
	}

	// Кидаем запрос в usecase
	err := h.fileUseCase.DeleteFile(context.Background(), fileMeta.LinkID, fileMeta.FileName)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "internal_server", "Failed to delete file")
		return
	}
	RespondSuccess(c, http.StatusOK, "file_deleted")
}

// GetZipFiles godoc
// @Summary Скачать все файлы ссылки в ZIP архиве
// @Description Собирает все файлы ссылки в ZIP и инициирует скачивание
// @Tags files
// @Produce application/zip
// @Param link_id query int64 true "ID ссылки"
// @Success 200 {file} application/zip {filename}="file_{link_id}.zip" "ZIP архив файлов"
// @Failure 400 {object} map[string]interface{} "Неверный link_id"
// @Failure 500 {object} map[string]interface{} "Ошибка создания архива"
// @Router /files/{id} [get]
func (h *Handler) GetZipFiles(c *gin.Context) {
	// получаем linkId
	linkIdStr := c.Param("link_id")
	linkId, err := strconv.ParseInt(linkIdStr, 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "bad_request", "Failed to parse link_id")
		return
	}
	// проверяем и собираем все файлы
	zipPath, err := h.fileUseCase.DownloadFiles(c.Request.Context(), linkId)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "internal_server", "Failed to download file")
		return
	}
	// отправляем файл
	fileName := fmt.Sprintf("file_%d.zip", linkId)
	c.FileAttachment(zipPath, fileName)
}
