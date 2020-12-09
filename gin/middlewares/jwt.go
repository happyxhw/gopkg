package middlewares

// https://github.com/appleboy/gin-jwt

import (
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/happyxhw/gopkg/logger"
	"go.uber.org/zap"
)

type Auth interface {
	Registry(*gin.Context)
	Authenticator(c *gin.Context) (interface{}, error)
	Authorizator(data interface{}, c *gin.Context) bool
	PayloadFunc(data interface{}) jwt.MapClaims
	Unauthorized(*gin.Context, int, string)
}

func NewJwt(
	key, identityKey string,
	auth Auth,
) *jwt.GinJWTMiddleware {
	jwtMid := jwt.GinJWTMiddleware{
		Realm:            "happy-token",
		SigningAlgorithm: "HS512",
		Key:              []byte(key),
		Timeout:          time.Hour,
		MaxRefresh:       time.Hour,
		Authenticator:    auth.Authenticator,
		Authorizator:     auth.Authorizator,
		PayloadFunc:      auth.PayloadFunc,
		Unauthorized:     auth.Unauthorized,
		IdentityKey:      identityKey,
		TokenLookup:      "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:    "Bearer",
		TimeFunc:         time.Now,
	}
	r, err := jwt.New(&jwtMid)
	if err != nil {
		logger.Fatal("init jwt middleware", zap.Error(err))
	}
	return r
}
