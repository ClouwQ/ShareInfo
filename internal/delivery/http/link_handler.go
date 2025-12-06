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

// CreateLink godoc
// @Summary Создать новую ссылку для загрузки файлов
// @Description Создает ссылку и возвращает linkId для прикрепления файлов
// @Tags links
// @Accept json
// @Produce json
// @Success 200 {object} int64 "linkId создан"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /links [post]
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

// GetLink godoc
// @Summary Получить метаданные и статистику ссылки
// @Description Возвращает информацию о файлах, скачиваниях и статусе ссылки
// @Tags links
// @Accept json
// @Produce json
// @Param id path int64 true "ID ссылки"
// @Success 200 {object} object "Метаданные ссылки"
// @Failure 400 {object} map[string]interface{} "Неверный ID ссылки"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /links/{id} [get]
func (h *Handler) GetLink(c *gin.Context) {
	linkIdStr := c.Param("id")
	linkId, err := strconv.ParseInt(linkIdStr, 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "bad_request", "Failed to parse link_id")
		return
	}

	// забираем метаданные файлов
	meta, err := h.linkUseCase.GetMeta(c.Request.Context(), linkId)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "internal_error", "Failed to get links meta")
		return
	}

	RespondSuccess(c, http.StatusOK, meta)
}

// FreezeLinkRequest godoc
// @Summary Параметры заморозки ссылки
type FreezeLinkRequest struct {
	ID             int64  `json:"id" binding:"required" example:"123"`
	Description    string `json:"description" binding:"omitempty" example:"Документы проекта"`
	IsOnceDownload bool   `json:"is_once_download" binding:"required" example:"false"`
	LifeTime       int    `json:"life_time" binding:"required" example:"24"` // В часах
}

// FreezeLink godoc
// @Summary Заморозить/настроить ссылку
// @Description Финализирует ссылку: устанавливает описание, время жизни, режим одноразового скачивания
// @Tags links
// @Accept json
// @Produce json
// @Param body body FreezeLinkRequest true "Параметры заморозки"
// @Success 200 {object} domain.Link "Ссылка заморожена"
// @Failure 422 {object} map[string]interface{} "Неверный JSON"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /links/freeze [post]
func (h *Handler) FreezeLink(c *gin.Context) {
	var req FreezeLinkRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		RespondError(c, http.StatusUnprocessableEntity, "invalid_request", "Failed to bind request body")
		return
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
		return
	}

	RespondSuccess(c, http.StatusOK, link)
}
