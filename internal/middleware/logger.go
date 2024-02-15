package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggingMiddleware logs the output of each call, as we know in releaseMode gin will turn off the logger by development mode.
func LoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		statusCode := c.Writer.Status()
		elapsed := time.Since(startTime)

		if err := c.Err(); err != nil {
			logger.Error("Error during request processing", zap.Error(err), zap.String("path", path), zap.String("method", method), zap.Duration("elapsed", elapsed))
		} else {
			logger.Info("Completed request", zap.String("path", path), zap.String("method", method), zap.Int("status", statusCode), zap.Duration("elapsed", elapsed))
		}
	}
}
