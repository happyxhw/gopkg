[![Build Status](https://github.com/happyxhw/gopkg/workflows/gopkg/badge.svg)](https://github.com/happyxhw/gopkg/workflows/gopkg/badge.svg)

[GITHUB](https://github.com/happyxhw/gopkg)

# 常用的 golang 库封装

> 本代码仅符合个人需求，建议有需要的可以下载最新 master 代码后按需修改



### 日志，基于 zap

默认使用 info 级别，‘console’ 输出，带调用行号
```go
func TestConsoleLogger(t *testing.T) {
	Info("test", zap.String("1", "2"))
}
```
可以自定义配置
```go
func TestConsoleLogger(t *testing.T) {
	c := Config{
		Level:       "info",
		FileName:    "xxx.log",   
		EncoderType: "console",
		Caller:      true,
	}
	InitLogger(&c)
	Info("test", zap.String("1", "2"))
}
```

生成新的 zap 实例：

```go
SetUp(zapcore.InfoLevel, "filename.log", "console", opts...)
```

配置：

```bash
Level: 日志级别：debug, info, warn, error
FileName: 不为空则输出到文件，warn以上（包括warn）日志输出到 xxx_err.log，不输出到终端，按天分割，只保存最近 7 天
EncoderType: console, json
Caller: 是否启用行号
```



### 数据库，基于 gorm v2

```go
func TestPostgres(t *testing.T) {
	l := logger.SetUp(zapcore.InfoLevel, "", "console")
	db, err := NewPostgresDb(&Config{
		User:         "xxx",
		Password:     "xxx",
		Host:         "127.0.0.1",
		Port:         "5432",
		DB:           "stravadb",
		MaxIdleConns: 10,
		MaxOpenConns: 10,
		Logger:       l.WithOptions(zap.AddCallerSkip(3), zap.AddCaller()),
		Level:        "info",
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	db.Exec("select idx from t_test")
}
```



### Redis，基于 go-redis v7

```go
func TestNewRedis(t *testing.T) {
	client, err := NewRedis(&Config{
		Host: "127.0.0.1:6379",
	})
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}
```



### Dispatcher，简单的并发工作池

```go
func TestNewDispatcher(t *testing.T) {
	d := NewDispatcher(10, -1)

	go func() {
		for {
			x, ok := <-d.ResultCh()
			if ok {
				fmt.Println(x)
			} else {
				fmt.Println("end")
				return
			}
		}
	}()

	for i := 0; i < 30; i++ {
		y := i
		err := d.Send(func() (interface{}, error) {
			time.Sleep(time.Second * 3)
			return y, nil
		})
		if err == ErrStopped {
			fmt.Println(err)
		}
	}
	d.Stop()
}
```



### GIN, 快速生成 HTTP 服务

```go
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
		User:         "xx",
		Password:     "xxxx",
		Host:         "127.0.0.1",
		Port:         "5432",
		DB:           "stravadb",
		MaxIdleConns: 1,
		MaxOpenConns: 1,
		Logger:       logger.GetLogger().WithOptions(zap.AddCallerSkip(2)),
		Level:        "info",
	})
	_ = db.AutoMigrate(&models.BaseUser{})
	key, identityKey := "test_key", "email"
	userHandler := user.NewUser(db, identityKey)
	jwtHandler := middlewares.NewJwt(key, identityKey, userHandler)

	v1 := r.Group("/api/v1")
	v1.POST("/auth/register", userHandler.Registry)
	v1.POST("/auth/login", jwtHandler.LoginHandler)
	v1.GET("/auth/refresh", jwtHandler.RefreshHandler)

	Serve(r, &c)
}
```

