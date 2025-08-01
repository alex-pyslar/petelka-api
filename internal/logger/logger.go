package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger предоставляет интерфейс для ведения журнала с различными уровнями.
type Logger struct {
	*zap.Logger
}

// NewLogger создает и инициализирует новый экземпляр Logger.
func NewLogger() (*Logger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"} // Вывод в консоль
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return &Logger{logger}, nil
}

// Info записывает информационное сообщение.
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

// Infof записывает форматированное информационное сообщение.
func (l *Logger) Infof(msg string, args ...interface{}) {
	l.Logger.Sugar().Infof(msg, args...)
}

// Warning записывает предупреждающее сообщение.
func (l *Logger) Warning(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}

// Warningf записывает форматированное предупреждающее сообщение.
func (l *Logger) Warningf(msg string, args ...interface{}) {
	l.Logger.Sugar().Warnf(msg, args...)
}

// Error записывает сообщение об ошибке.
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}

// Errorf записывает форматированное сообщение об ошибке.
func (l *Logger) Errorf(msg string, args ...interface{}) {
	l.Logger.Sugar().Errorf(msg, args...)
}

// Fatal записывает сообщение о фатальной ошибке и завершает работу приложения.
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.Logger.Fatal(msg, fields...)
}

// Fatalf записывает форматированное сообщение о фатальной ошибке и завершает работу.
func (l *Logger) Fatalf(msg string, args ...interface{}) {
	l.Logger.Sugar().Fatalf(msg, args...)
}
