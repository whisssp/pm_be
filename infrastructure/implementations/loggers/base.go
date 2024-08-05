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
	c       *gin.Context
	//span          trace.Span
	honeycombRepo *honeycomb.HoneycombRepository
}

const (
	Zap       = "Zap"
	Honeycomb = "Honeycomb"
	Slack     = "Slack"
)

type Option func(*LoggerRepo)

// NewLoggerRepository creates a new logger repository based on the specified channels
func NewLoggerRepository(channels []string) *LoggerRepo {
	var loggers []loggers.LoggerRepository
	var honeycombRepo *honeycomb.HoneycombRepository

	if len(channels) < 1 {
		return &LoggerRepo{
			loggers:       loggers,
			c:             nil,
			honeycombRepo: honeycombRepo,
		}
	}

	for _, channel := range channels {
		switch channel {
		case Zap:
			loggers = append(loggers, zap.NewZapRepository())
		case Honeycomb:
			honeycombRepo = honeycomb.NewHoneycombRepository()
			loggers = append(loggers, honeycombRepo)
		default:
			// You might want to log or handle unsupported channels
			continue
		}
	}

	//var span trace.Span
	//if honeycombRepo != nil {
	//	span = honeycombRepo.GetSpan()
	//}
	return &LoggerRepo{loggers: loggers, c: nil, honeycombRepo: honeycombRepo}
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
func (l *LoggerRepo) Start(c *gin.Context, info string, options ...Option) trace.Span {
	l.c = c

	var span trace.Span
	for _, logger := range l.loggers {
		_, ok := logger.(*honeycomb.HoneycombRepository)
		if ok {
			span = logger.Start(c, info)
		} else {
			logger.Start(c, info)
		}
	}

	// Execute optional functions
	for _, opts := range options {
		opts(l)
	}

	return span
}

func (l *LoggerRepo) End() {
	l.honeycombRepo.GetSpan().End()
}

func (l *LoggerRepo) SetContextWithSpanFunc() Option {
	return func(l *LoggerRepo) {
		l.c.Set("otel_context", trace.ContextWithSpan(l.c, l.honeycombRepo.GetSpan()))
	}
}

func (l *LoggerRepo) SetContextWithSpan(span trace.Span) {
	l.c.Set("otel_context", trace.ContextWithSpan(l.c, span))
}

func (l *LoggerRepo) UseGivenSpan(span trace.Span) Option {
	return func(l *LoggerRepo) {
		l.honeycombRepo.UseGivenSpan(span)
	}
}