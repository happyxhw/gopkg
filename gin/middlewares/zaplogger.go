package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinZap returns a gin.HandlerFunc (middleware) that logs requests using uber-go/zap.
//
// Requests with errors are logged using zap.Error().
// Requests without errors are logged using zap.Info().
func GinZap(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		reqID := GetReqID(c)

		switch {
		case statusCode >= 400 && statusCode <= 499:
			logger.Warn("[GIN]",
				zap.Int("code", statusCode),
				zap.String("latency", latency.String()),
				zap.String("ip", clientIP),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ua", c.Request.UserAgent()),
				zap.String("err", c.Errors.String()),
				zap.String("X-Request-Id", reqID),
			)
			c.JSON(statusCode, gin.H{"code": statusCode, "msg": c.Errors.String()})
		case statusCode >= 500:
			logger.Error("[GIN]",
				zap.Int("code", statusCode),
				zap.String("latency", latency.String()),
				zap.String("ip", clientIP),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ua", c.Request.UserAgent()),
				zap.String("err", c.Errors.String()),
				zap.String("X-Request-Id", reqID),
			)
			c.JSON(statusCode, gin.H{"code": statusCode, "msg": c.Errors.String()})
		default:
			logger.Info("[GIN]",
				zap.Int("code", statusCode),
				zap.String("latency", latency.String()),
				zap.String("ip", clientIP),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ua", c.Request.UserAgent()),
				zap.String("err", c.Errors.String()),
				zap.String("X-Request-Id", reqID),
			)
		}
	}
}
