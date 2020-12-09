package user

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/happyxhw/gopkg/gin/models"
	"github.com/happyxhw/gopkg/logger"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrParam         = errors.New("invalid parameters")
	ErrExists        = errors.New("email exists")
	ErrNotExists     = errors.New("user not exists")
	ErrDbInternal    = errors.New("internal db error")
	ErrRedisInternal = errors.New("internal redis error")
	ErrCode          = errors.New("validate code error")
)

type User struct {
	db          *gorm.DB
	red         *redis.Client
	identityKey string
}

func NewUser(db *gorm.DB, red *redis.Client, identityKey string) *User {
	u := User{
		db:          db,
		red:         red,
		identityKey: identityKey,
	}
	return &u
}

func (u User) Register(c *gin.Context) {
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

func (u User) RequestPass(c *gin.Context) {
	var user models.BaseUser
	if err := c.ShouldBindJSON(&user); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, ErrParam)
		return
	}

	var dbUser models.BaseUser
	_ = u.db.Select("id").Where("email = ?", user.Email).First(&dbUser).Error
	if dbUser.ID == 0 {
		_ = c.AbortWithError(http.StatusBadRequest, ErrNotExists)
		return
	}

	code := getRandomString(7)
	err := u.red.Set(fmt.Sprintf("validation_%d", dbUser.ID), code, time.Minute*5).Err()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, ErrRedisInternal)
		return
	}
	/*
	   sending code to email
	*/
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "ok",
	})
}

func (u User) ResetPass(c *gin.Context) {
	user := models.BaseUser{}
	if err := c.ShouldBindJSON(&user); err != nil || user.Password == "" {
		_ = c.AbortWithError(http.StatusBadRequest, ErrParam)
		return
	}

	var dbUser models.BaseUser
	_ = u.db.Select("id").Where("email = ?", user.Email).First(&dbUser).Error
	if dbUser.ID == 0 {
		_ = c.AbortWithError(http.StatusBadRequest, ErrNotExists)
		return
	}

	code, err := u.red.Get(fmt.Sprintf("validation_%d", dbUser.ID)).Result()
	if err != nil || code != user.Code {
		_ = c.AbortWithError(http.StatusBadRequest, ErrCode)
		return
	}

	passwordByte, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	newPassword := string(passwordByte)
	err = u.db.Model(&models.BaseUser{}).Where("email = ?", user.Email).Update("password", newPassword).Error
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, ErrDbInternal)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
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

// getRandomString 随机生成大写字母和数字组合
func getRandomString(l int) string {
	str := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	n := int64(len(bytes))
	var result []byte
	for i := 0; i < l; i++ {
		t, _ := rand.Int(rand.Reader, big.NewInt(n))
		result = append(result, bytes[t.Int64()])
	}
	return string(result)

}
