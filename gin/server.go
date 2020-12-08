package gin

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/happyxhw/gopkg/logger"

	"go.uber.org/zap"
)

type Config struct {
	Addr string
	Mode string
}

func Serve(router *gin.Engine, c *Config) {
	logger.Info("http server start")

	// init router
	gin.SetMode(c.Mode)
	pprof.Register(router, "dev/pprof")

	server := &http.Server{
		Addr:           c.Addr,
		Handler:        router,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	logger.Info("start http server listening", zap.String("addr", c.Addr))

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutdown server")

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("server shutdown err", zap.Error(err))
	}
	logger.Info("server exited")
}
