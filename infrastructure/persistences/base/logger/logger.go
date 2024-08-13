package logger

import (
	"os"
	"pm/infrastructure/implementations/loggers"
	"strings"
)

var logger *loggers.LoggerRepo

func NewLogger() {
	logChannels := os.Getenv("LOGGER_CHANNELS")
	logChannelsList := strings.Split(logChannels, ",")
	logger = loggers.NewLoggerRepository(logChannelsList)
}

func GetLogger() *loggers.LoggerRepo {
	return logger
}
