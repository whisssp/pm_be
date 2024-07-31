package repository

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

type LoggerRepository interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	Fatal(msg string, fields map[string]interface{})
	Start(c *gin.Context, info string) trace.Span
}