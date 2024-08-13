package loggers

import (
	"github.com/gin-gonic/gin"
	loggers2 "pm/domain/repository/loggers"
	honeycomb "pm/infrastructure/implementations/loggers/honeycomb"
	"pm/infrastructure/implementations/loggers/zap"
)

// LoggerRepo is a logger repository that can use multiple channels
type LoggerRepo struct {
	loggers         []loggers2.LoggerRepository
	loggerMap       map[string]loggers2.LoggerRepository
	logChannelsList []string
}

const (
	Zap       = "Zap"
	Honeycomb = "Honeycomb"
	Slack     = "Slack"
)

var loggers *LoggerRepo

// NewLoggerRepository creates a new logger repository based on the specified channels
func NewLoggerRepository(channels []string) *LoggerRepo {
	var loggerMap = make(map[string]loggers2.LoggerRepository)

	if len(channels) == 0 {
		return &LoggerRepo{loggers: nil, loggerMap: loggerMap, logChannelsList: channels}
	}
	for _, channel := range channels {
		switch channel {
		case Zap:
			loggerMap[Zap] = zap.NewZapRepository()
		case Honeycomb:
		default:
			// You might want to log or handle unsupported channels
			continue
		}
	}

	return &LoggerRepo{loggers: nil, loggerMap: loggerMap, logChannelsList: channels}
}

// Debug logs a debug message
func (l *LoggerRepo) Debug(msg string, fields map[string]interface{}) {
	for _, logger := range l.loggers {
		logger.Debug(msg, fields)
	}
}

// Info logs an info message
func (l *LoggerRepo) Info(msg string, fields map[string]interface{}) {
	for _, logger := range l.loggers {
		logger.Info(msg, fields)
	}
}

// Warn logs a warning message
func (l *LoggerRepo) Warn(msg string, fields map[string]interface{}) {
	for _, logger := range l.loggers {
		logger.Warn(msg, fields)
	}
}

// Error logs an error message
func (l *LoggerRepo) Error(msg string, fields map[string]interface{}) {
	for _, logger := range l.loggers {
		logger.Error(msg, fields)
	}
}

// Fatal logs a fatal message
func (l *LoggerRepo) Fatal(msg string, fields map[string]interface{}) {
	for _, logger := range l.loggers {
		logger.Fatal(msg, fields)
	}
}

// Start function
func (l *LoggerRepo) Start(c *gin.Context, info string) (*gin.Context, *LoggerRepo) {
	ctx := c
	var loggerList []loggers2.LoggerRepository
	for _, logger := range l.logChannelsList {
		log, ok := l.loggerMap[logger]
		if ok {
			log.Start(c, info)
			loggerList = append(loggerList, log)
		} else {
			switch logger {
			case Honeycomb:
				honeycombLog := honeycomb.NewHoneycombRepository()
				ctx, _ = honeycombLog.Start(c, info)
				loggerList = append(loggerList, honeycombLog)
			default:
				continue
			}
		}
	}

	return ctx, &LoggerRepo{
		loggers:         loggerList,
		logChannelsList: l.logChannelsList,
		loggerMap:       l.loggerMap,
	}
}

func (l *LoggerRepo) End() {
	for _, logger := range l.loggers {
		logger.End()
	}
}
