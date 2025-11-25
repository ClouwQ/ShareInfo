package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

// Logger для всех запросов
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		// Время ответа
		duration := time.Since(start)

		fields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.Duration("duration", duration),
			zap.Int("size", c.Writer.Size()),
		}

		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		switch {
		case c.Writer.Status() >= 500:

			logger.Error("Server error", fields...)
		case c.Writer.Status() >= 400:
			logger.Warn("Client error", fields...)
		case c.Writer.Status() >= 300:
			logger.Info("Redirect", fields...)
		default:

			logger.Info("Request", fields...)
		}
	}
}
