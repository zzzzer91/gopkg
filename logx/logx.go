package logx

import (
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger = func() *zap.SugaredLogger {
	config := zap.NewProductionConfig()
	config.DisableStacktrace = true
	config.Level.SetLevel(zapcore.InfoLevel)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeCaller = func(ec zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		if !ec.Defined {
			enc.AppendString("undefined")
		} else {
			enc.AppendString(ec.Function + ":" + strconv.Itoa(ec.Line))
		}
	}
	config.EncoderConfig = encoderConfig
	l, _ := config.Build(zap.AddCallerSkip(1))
	sugar := l.Sugar()
	return sugar
}()

func Infof(format string, args ...any) {
	logger.Infof(format, args...)
}

func Warnf(format string, args ...any) {
	logger.Warnf(format, args...)
}

func Errorf(format string, args ...any) {
	logger.Errorf(format, args...)
}

func Info(args ...any) {
	logger.Info(args...)
}

func Warn(args ...any) {
	logger.Warn(args...)
}

func Error(args ...any) {
	logger.Error(args...)
}

func Infoln(args ...any) {
	logger.Infoln(args...)
}

func Warnln(args ...any) {
	logger.Warnln(args...)
}

func Errorln(args ...any) {
	logger.Errorln(args...)
}
