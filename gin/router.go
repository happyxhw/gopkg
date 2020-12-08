package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/happyxhw/gopkg/gin/middlewares"
	"go.uber.org/zap"
)

func NewEngine(logger *zap.Logger) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middlewares.RequestId())
	r.Use(middlewares.GinZap(logger))
	return r
}
