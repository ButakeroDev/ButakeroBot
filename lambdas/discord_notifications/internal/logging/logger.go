package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Info(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
}

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger() (*ZapLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &ZapLogger{logger: logger}, nil
}

func (l *ZapLogger) Close() error {
	err := l.logger.Sync()
	if err != nil && err.Error() != "sync /dev/stderr: invalid argument" {
		return err
	}
	return nil
}

func (l *ZapLogger) Info(msg string, fields ...zapcore.Field) {
	l.logger.Info(msg, fields...)
}

func (l *ZapLogger) Error(msg string, fields ...zapcore.Field) {
	l.logger.Error(msg, fields...)
}
