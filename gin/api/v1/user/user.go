package user

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/happyxhw/gopkg/gin/models"
	"github.com/happyxhw/gopkg/logger"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrParam      = errors.New("invalid parameters")
	ErrExists     = errors.New("email exists")
	ErrDbInternal = errors.New("internal db error")
)

type User struct {
	db          *gorm.DB
	identityKey string
}

func NewUser(db *gorm.DB, identityKey string) *User {
	u := User{
		db:          db,
		identityKey: identityKey,
	}
	return &u
}

func (u User) Registry(c *gin.Context) {
	var user models.BaseUser
	if err := c.ShouldBindJSON(&user); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, ErrParam)
		return
	}

	passwordByte, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(passwordByte)

	var dbUser models.BaseUser
	_ = u.db.Select("id").Where("email = ?", user.Email).First(&dbUser).Error
	if dbUser.ID > 0 {
		_ = c.AbortWithError(http.StatusBadRequest, ErrExists)
		return
	}

	if err := u.db.Create(&user).Error; err != nil {
		logger.Error("create user", zap.Error(err))
		_ = c.AbortWithError(http.StatusInternalServerError, ErrDbInternal)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"code": http.StatusCreated,
		"msg":  "ok",
	})
}

func (u User) Authenticator(c *gin.Context) (interface{}, error) {
	var user models.BaseUser
	var dbUser models.BaseUser
	if err := c.ShouldBindJSON(&user); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	_ = u.db.Select("user_name, email, password, created_at").Where("email = ?", user.Email).Find(&dbUser).Error
	if dbUser.Password != "" {
		err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
		if err == nil {
			return dbUser, nil
		}
	}
	return nil, jwt.ErrFailedAuthentication
}

func (u User) Authorizator(data interface{}, c *gin.Context) bool {
	return true
}

func (u User) PayloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(*models.BaseUser); ok {
		return jwt.MapClaims{
			u.identityKey: v.Email,
			"user_name":   v.UserName,
			"create_time": v.CreatedAt,
		}
	}
	return jwt.MapClaims{}
}

func (u User) Unauthorized(c *gin.Context, code int, msg string) {
	_ = c.AbortWithError(code, errors.New(msg))
}
