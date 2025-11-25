package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery — покрытие паник
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Логируем trace stack
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.String("stack", string(debug.Stack())),
				)

				// Отправляем ответ клиенту
				c.JSON(500, gin.H{
					"success": false,
					"error":   "internal_server_error",
					"message": "Internal server error occurred",
				})

				// TODO отправить уведомление админу

				c.Abort()
			}
		}()

		c.Next()
	}
}

// TimeoutMiddleware добавляет таймаут к запросам
func TimeoutMiddleware(timeout time.Duration, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		finished := make(chan struct{})
		go func() {
			c.Next()
			finished <- struct{}{}
		}()

		select {
		case <-finished:
			return
		case <-ctx.Done():
			logger.Warn("Request timeout",
				zap.String("path", c.Request.URL.Path),
				zap.Duration("timeout", timeout),
			)
			// 408 код
			c.JSON(http.StatusRequestTimeout, gin.H{
				"success": false,
				"error":   "request_timeout",
				"message": fmt.Sprintf("Request timeout after %v", timeout),
			})
			c.Abort()
		}
	}
}
