package logger

import (
	"github.com/Nemo08/NCTW/services/api"
	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	gormzap "github.com/wantedly/gorm-zap"
	"go.uber.org/zap"
)

type Logr struct {
	z     zap.Logger
	level int //0 - all, 1 - warning and up, 2 - erorrs
}

func NewLogger() *Logr {
	zl, _ := zap.NewProduction()
	defer zl.Sync()
	return &Logr{
		z: *zl,
	}
}

func (lg *Logr) WithField(key string, value interface{}) *Logr {
	return &Logr{
		z: *lg.z.With(zap.String(key, value.(string))),
	}
}

func (lg *Logr) WithContext(ctx api.Context) *Logr {
	return &Logr{
		z: *lg.z.With(zap.String("request_id", ctx.Response().Header().Get("X-Request-ID"))),
	}
}

func (lg *Logr) Info(args ...interface{}) {
	if lg.level == 0 {
		lg.z.Sugar().Info(args)
	}
}

func (lg *Logr) Warn(args ...interface{}) {
	if lg.level <= 1 {
		lg.z.Sugar().Warn(args)
	}
}

func (lg *Logr) Error(args ...interface{}) {
	if lg.level <= 2 {
		lg.z.Sugar().Error(args)
	}
}

func (lg *Logr) SetLogLevel(l int) {
	if l > 2 {
		l = 2
	}
	if l < 0 {
		l = 0
	}
	lg.level = l
}

func (lg *Logr) GormLogger() *gormzap.Logger {
	return gormzap.New(&lg.z)
}

func (lg *Logr) EchoLogger() echo.MiddlewareFunc {
	return echozap.ZapLogger(&lg.z)
}
