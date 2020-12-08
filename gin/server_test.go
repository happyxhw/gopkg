package gin

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/happyxhw/gopkg/logger"
	"go.uber.org/zap"
)

func TestServe(t *testing.T) {
	r := NewEngine(logger.GetLogger().WithOptions(zap.AddCallerSkip(-1)))
	r.GET("/api/v1/greeter", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "reply",
		})
	})
	c := Config{
		Addr: "0.0.0.0:8080",
		Mode: "debug",
	}
	Serve(r, &c)
}
