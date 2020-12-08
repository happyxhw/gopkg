package gin

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/happyxhw/gopkg/gin/middlewares"
	"go.uber.org/zap"
)

func NewEngine(logger *zap.Logger) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middlewares.RequestId())
	r.Use(middlewares.GinZap(logger))
	// init cor middleware
	corConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
	corConfig.AllowAllOrigins = true

	r.Use(cors.New(corConfig))
	return r
}
