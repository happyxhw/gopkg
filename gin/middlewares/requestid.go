package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

const key = "X-Request-Id"

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for incoming header, use it if exists
		requestID := c.Request.Header.Get("X-Request-Id")

		// Create request id with UUID4
		if requestID == "" {
			uuid4, _ := uuid.NewV4()
			requestID = uuid4.String()
		}

		// Expose it for use in the application
		c.Set(key, requestID)

		// Set X-Request-Id header
		c.Writer.Header().Set("X-Request-Id", requestID)
		c.Next()
	}
}

func GetReqID(c *gin.Context) string {
	return c.GetString(key)
}
