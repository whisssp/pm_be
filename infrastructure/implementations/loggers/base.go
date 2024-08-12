package loggers

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"pm/domain/repository/loggers"
	honeycomb "pm/infrastructure/implementations/loggers/honeycomb"
	"pm/infrastructure/implementations/loggers/zap"
)

type LoggerRepo struct {
	loggers []loggers.LoggerRepository
}

const (
	Zap       = "Zap"
	Honeycomb = "Honeycomb"
	Slack     = "Slack"
)

type Option func(*LoggerRepo)

var logger *LoggerRepo

// NewLoggerRepository creates a new logger repository based on the specified channels
func NewLoggerRepository(channels []string) *LoggerRepo {
	var loggers []loggers.LoggerRepository

	if len(channels) < 1 {
		return &LoggerRepo{
			loggers: loggers,
		}
	}

	for _, channel := range channels {
		switch channel {
		case Zap:
			loggers = append(loggers, zap.NewZapRepository())
		case Honeycomb:
			loggers = append(loggers, honeycomb.NewHoneycombRepository())
		default:
			// You might want to log or handle unsupported channels
			continue
		}
	}

	return &LoggerRepo{loggers: loggers}
}

// Debug logs a debug message
func (l *LoggerRepo) Debug(msg string, fields map[string]interface{}, options ...Option) {
	for _, logger := range l.loggers {
		logger.Debug(msg, fields)
	}
}

// Info logs an info message
func (l *LoggerRepo) Info(msg string, fields map[string]interface{}, options ...Option) {
	for _, logger := range l.loggers {
		logger.Info(msg, fields)
	}
}

// Warn logs a warning message
func (l *LoggerRepo) Warn(msg string, fields map[string]interface{}, options ...Option) {
	for _, logger := range l.loggers {
		logger.Warn(msg, fields)
	}
}

// Error logs an error message
func (l *LoggerRepo) Error(msg string, fields map[string]interface{}, options ...Option) {
	for _, logger := range l.loggers {
		logger.Error(msg, fields)
	}
}

// Fatal logs a fatal message
func (l *LoggerRepo) Fatal(msg string, fields map[string]interface{}, options ...Option) {
	for _, logger := range l.loggers {
		logger.Fatal(msg, fields)
	}
}

// Start function
func (l *LoggerRepo) Start(c *gin.Context, info string, options ...Option) (*gin.Context, trace.Span) {
	var ctx *gin.Context
	var span trace.Span
	for _, logger := range l.loggers {
		_, ok := logger.(*honeycomb.HoneycombRepository)
		if ok {
			ctx, span = logger.Start(c, info)
		} else {
			logger.Start(c, info)
		}
	}

	return ctx, span
}

func (l *LoggerRepo) End() {
	for _, logger := range l.loggers {
		logger.End()
	}
}
