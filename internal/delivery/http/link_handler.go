package http

import (
	"ShareInfo/internal/domain"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

// CreateLink создает, сохраняет "ссылку" и присылает ее linkId, чтобы прикрепить к ней файлы
func (h *Handler) CreateLink(c *gin.Context) {
	// Создаем ссылку
	linkId, err := h.linkUseCase.Create(context.Background())
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "internal_error", "Failed to create link")
		return
	}

	h.logger.Debug("Create link request received",
		zap.Int64("linkId", linkId),
	)

	RespondSuccess(c, http.StatusOK, linkId)
}

// GetLink возвращает всю статистику по файлам
func (h *Handler) GetLink(c *gin.Context) {
	linkIdStr := c.Query("link_id")
	linkId, err := strconv.ParseInt(linkIdStr, 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "bad_request", "Failed to parse link_id")
	}

	// забираем метаданные файлов
	meta, err := h.linkUseCase.GetMeta(c.Request.Context(), linkId)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "internal_error", "Failed to get links meta")
	}

	RespondSuccess(c, http.StatusOK, meta)
}

type FreezeLinkRequest struct {
	ID             int64  `json:"id" binding:"required"`
	Description    string `json:"description" binding:"omitempty"`
	IsOnceDownload bool   `json:"is_once_download" binding:"required"`
	LifeTime       int    `json:"life_time" binding:"required"` // В часах
}

func (h *Handler) FreezeLink(c *gin.Context) {
	var req FreezeLinkRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		RespondError(c, http.StatusUnprocessableEntity, "invalid_request", "Failed to bind request body")
	}

	h.logger.Debug("Freeze link request received",
		zap.Any("request", req),
		zap.Int64("linkId", req.ID),
	)

	createdAt := time.Now()
	expiredAt := createdAt.Add(time.Duration(req.LifeTime) * time.Hour)
	link := domain.Link{
		ID:             req.ID,
		Description:    req.Description,
		IsActive:       true,
		IsOnceDownload: req.IsOnceDownload,
		IsFrozen:       true,
		CreatedAt:      createdAt,
		ExpiresAt:      expiredAt,
	}

	err := h.linkRepo.Freeze(context.Background(), link)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "internal_error", "Failed to freeze link")
	}

	RespondSuccess(c, http.StatusOK, link)
}
