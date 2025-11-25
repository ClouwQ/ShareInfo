package http

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	logger *zap.Logger

	// Usecases

}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) HealthCheck(c *gin.Context) {
	RespondSuccess(c, http.StatusOK, gin.H{
		"status":  "UP",
		"service": "link-service",
		"version": "1.0.0",
	})
}
