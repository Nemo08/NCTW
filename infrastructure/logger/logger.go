package logger

import (
	//"github.com/Nemo08/NCTW/services/api"
	"os"

	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/wantedly/gorm-zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel  = zap.DebugLevel
	InfoLevel   = zap.InfoLevel
	WarnLevel   = zap.WarnLevel
	ErrorLevel  = zap.ErrorLevel
	DPanicLevel = zap.DPanicLevel
	PanicLevel  = zap.PanicLevel
	FatalLevel  = zap.FatalLevel
)

type Logr struct {
	zap.SugaredLogger
	atom *zap.AtomicLevel
}

var Log Logr

func newLogger() Logr {
	atom := zap.NewAtomicLevel()
	encoderCfg := zap.NewProductionEncoderConfig()

	zl := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.Lock(os.Stdout),
			atom,
		),
		zap.AddCaller(),
	)

	defer zl.Sync()

	return Logr{
		SugaredLogger: *zl.Sugar(),
		atom:          &atom,
	}
}

func (lg *Logr) WithField(key string, value interface{}) *Logr {
	return &Logr{
		SugaredLogger: *lg.With(zap.String(key, value.(string))),
		atom:          lg.atom,
	}
}

func (lg *Logr) SetLevel(l zapcore.Level) {
	lg.atom.SetLevel(l)
}

func (lg *Logr) GormLogger() *gormzap.Logger {
	return gorm-zap.New(lg.Desugar())

func (lg *Logr) EchoLogger() echo.MiddlewareFunc{
	return echozap.ZapLogger(lg.Desugar())
}

func init() {
	Log = newLogger()
}
