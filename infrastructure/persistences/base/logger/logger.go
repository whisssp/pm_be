package logger

import (
	"os"
	"pm/infrastructure/implementations/loggers"
	"strings"
)

func NewLogger() *loggers.LoggerRepo {

	logChannels := os.Getenv("LOGGER_CHANNELS")
	logChannelsList := strings.Split(logChannels, ",")
	logger := loggers.NewLoggerRepository(logChannelsList)

	return logger
}
