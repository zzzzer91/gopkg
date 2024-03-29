package logx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var loggerLevel = zap.NewAtomicLevel()

var logger = func() *zap.SugaredLogger {
	config := zap.NewProductionConfig()
	config.DisableStacktrace = true
	config.Level = loggerLevel
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig = encoderConfig
	l, _ := config.Build(zap.AddCallerSkip(1))
	sugar := l.Sugar()
	return sugar
}()

func SetLevel(level int) {
	loggerLevel.SetLevel(zapcore.Level(level))
}

func Debugf(format string, args ...any) {
	logger.Debugf(format, args...)
}

func Infof(format string, args ...any) {
	logger.Infof(format, args...)
}

func Warnf(format string, args ...any) {
	logger.Warnf(format, args...)
}

func Errorf(format string, args ...any) {
	logger.Errorf(format, args...)
}

func Fatalf(format string, args ...any) {
	logger.Fatalf(format, args...)
}

func Debug(args ...any) {
	logger.Debugln(args...)
}

func Info(args ...any) {
	logger.Infoln(args...)
}

func Warn(args ...any) {
	logger.Warnln(args...)
}

func Error(args ...any) {
	logger.Errorln(args...)
}

func Fatal(args ...any) {
	logger.Fatalln(args...)
}
