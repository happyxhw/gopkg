package dbgo

import (
	"context"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	gLogger "gorm.io/gorm/logger"
)

type gormLogger struct {
	logger           *zap.Logger
	LogLevel         gLogger.LogLevel
	SlowThreshold    time.Duration
	SkipCallerLookup bool
}

func newLogger(l *zap.Logger, level string) gormLogger {
	ll := gLogger.Warn
	switch strings.ToLower(level) {
	case "info":
		ll = gLogger.Info
	case "error":
		ll = gLogger.Error
	case "silent":
		ll = gLogger.Silent
	}
	gl := gormLogger{
		logger:           l,
		LogLevel:         ll,
		SlowThreshold:    time.Second,
		SkipCallerLookup: true,
	}
	return gl
}

func (gl gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return gormLogger{
		logger:           gl.logger,
		LogLevel:         level,
		SlowThreshold:    gl.SlowThreshold,
		SkipCallerLookup: gl.SkipCallerLookup,
	}
}

func (gl gormLogger) Info(_ context.Context, s string, i ...interface{}) {
	if gl.LogLevel < gLogger.Info {
		return
	}
	gl.logger.Sugar().Infof(s, i...)
}

func (gl gormLogger) Warn(_ context.Context, s string, i ...interface{}) {
	if gl.LogLevel < gLogger.Warn {
		return
	}
	gl.logger.Sugar().Warnf(s, i...)
}

func (gl gormLogger) Error(_ context.Context, s string, i ...interface{}) {
	if gl.LogLevel < gLogger.Error {
		return
	}
	gl.logger.Sugar().Errorf(s, i...)
}

func (gl gormLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if gl.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && gl.LogLevel >= gLogger.Error:
		sql, rows := fc()
		gl.logger.Error("trace", zap.Error(err), zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	case gl.SlowThreshold != 0 && elapsed > gl.SlowThreshold && gl.LogLevel >= gLogger.Warn:
		sql, rows := fc()
		gl.logger.Warn("trace", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	case gl.LogLevel >= gLogger.Info:
		sql, rows := fc()
		gl.logger.Info("trace", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	}
}
