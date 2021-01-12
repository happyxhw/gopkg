package gin

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/happyxhw/gopkg/dbgo"
	"github.com/happyxhw/gopkg/gin/api/v1/user"
	"github.com/happyxhw/gopkg/gin/middlewares"
	"github.com/happyxhw/gopkg/gin/models"
	"github.com/happyxhw/gopkg/goredis"
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

	db, _ := dbgo.NewPostgresDb(&dbgo.Config{
		User:         "happyxhw",
		Password:     "808258XXxx",
		Host:         "127.0.0.1",
		Port:         5432,
		DB:           "stravadb",
		MaxIdleConns: 1,
		MaxOpenConns: 1,
		Logger:       logger.GetLogger().WithOptions(zap.AddCallerSkip(2)),
		Level:        "info",
	})
	red, _ := goredis.NewRedis(&goredis.Config{
		Host: "127.0.0.1:6379",
	})
	_ = db.AutoMigrate(&models.BaseUser{})
	key, identityKey := "test_key", "email"
	userHandler := user.NewUser(db, red, identityKey)
	jwtHandler := middlewares.NewJwt(key, identityKey, userHandler)

	v1 := r.Group("/api/v1")
	v1.POST("/auth/register", userHandler.Register)
	v1.POST("/auth/login", jwtHandler.LoginHandler)
	v1.GET("/auth/refresh", jwtHandler.RefreshHandler)
	v1.POST("/auth/request-pass", userHandler.RequestPass)
	v1.POST("/auth/reset-pass", userHandler.ResetPass)

	Serve(r, &c)
}
