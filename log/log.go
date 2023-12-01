package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func init() {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cfg.EncoderConfig.EncodeTime = customTimeEncoder
	cfg.Encoding = "console"

	l, _ := cfg.Build()
	logger = l.Sugar()
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func Debugf(msg string, args ...interface{}) {
	logger.Debugf(msg, args...)
}

func Infof(msg string, args ...interface{}) {
	logger.Infof(msg, args...)
}

func Warnf(msg string, args ...interface{}) {
	logger.Warnf(msg, args...)
}

func Errorf(msg string, args ...interface{}) {
	logger.Errorf(msg, args...)
}

func Fatalf(msg string, args ...interface{}) {
	logger.Fatalf(msg, args...)
}
