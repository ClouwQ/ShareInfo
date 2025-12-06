package http

import (
	"ShareInfo/internal/repository/database"
	"ShareInfo/internal/usecase"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	logger *zap.Logger

	// Usecases
	linkUseCase        usecase.Link
	fileUseCase        usecase.UploadedFile
	linkMessageUseCase usecase.LinkMessage

	// Repository
	linkRepo         database.LinkRepository
	uploadedFileRepo database.UploadedFileRepository
}

func NewHandler(
	logger *zap.Logger,
	linkUseCase usecase.Link,
	fileUseCase usecase.UploadedFile,
	linkMessageUseCase usecase.LinkMessage,
	linkRepo database.LinkRepository,
	uploadedFileRepo database.UploadedFileRepository,
) *Handler {
	return &Handler{
		logger:             logger,
		linkUseCase:        linkUseCase,
		fileUseCase:        fileUseCase,
		linkMessageUseCase: linkMessageUseCase,
		linkRepo:           linkRepo,
		uploadedFileRepo:   uploadedFileRepo,
	}
}

func (h *Handler) HealthCheck(c *gin.Context) {
	RespondSuccess(c, http.StatusOK, gin.H{
		"status":  "UP",
		"service": "link-service",
		"version": "1.0.0",
	})
}
