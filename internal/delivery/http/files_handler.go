package http

import (
	usecaseFileSerice "ShareInfo/internal/usecase/file"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type FileMetaRequest struct {
	LinkID int64 `form:"link_id" binding:"required"`
}

func (h *Handler) UploadFile(c *gin.Context) {
	var savedFile usecaseFileSerice.File
	var req FileMetaRequest
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

	// Кидаем в usecase
	err = h.fileUseCase.SaveFile(context.Background(), savedFile)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "internal_server", "Failed to save file")
		return
	}
	RespondSuccess(c, http.StatusOK, "file_uploaded")
}
