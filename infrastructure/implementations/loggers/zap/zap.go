package zap

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	zap "go.uber.org/zap"
	"os"
)

type ZapRepo struct {
	sugar *zap.SugaredLogger
}

// NewZapRepository creates a new Zap logger repository
func NewZapRepository() *ZapRepo {
	// Initialize Zap
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	return &ZapRepo{zapLogger.Sugar()}
}

// Debug logs a debug message
func (z *ZapRepo) Debug(msg string, fields map[string]interface{}) {
	if fields == nil {
		z.sugar.Debug(msg)
		return
	}
	z.sugar.Debugw(msg, z.convertFields(fields)...)
}

// Info logs an info message
func (z *ZapRepo) Info(msg string, fields map[string]interface{}) {
	if fields == nil {
		z.sugar.Info(msg)
		return
	}
	z.sugar.Infow(msg, z.convertFields(fields)...)
}

// Warn logs a warning message
func (z *ZapRepo) Warn(msg string, fields map[string]interface{}) {
	if fields == nil {
		z.sugar.Warn(msg)
		return
	}
	z.sugar.Warnw(msg, z.convertFields(fields)...)
}

// Error logs an error message
func (z *ZapRepo) Error(msg string, fields map[string]interface{}) {
	if fields == nil {
		z.sugar.Error(msg)
		return
	}
	z.sugar.Errorw(msg, z.convertFields(fields)...)
}

// Fatal logs a fatal message
func (z *ZapRepo) Fatal(msg string, fields map[string]interface{}) {
	if fields == nil {
		z.sugar.Fatal(msg)
		os.Exit(1)
		return
	}
	z.sugar.Fatalw(msg, z.convertFields(fields)...)
	os.Exit(1)
}

// convertFields converts fields into Zap-compatible fields
func (z *ZapRepo) convertFields(fields map[string]interface{}) []interface{} {
	// Convert fields to Zap fields
	var zapFields []interface{}
	for key, value := range fields {
		zapFields = append(zapFields, key, value)
	}
	return zapFields
}

func (z *ZapRepo) Start(c *gin.Context, info string) (*gin.Context, trace.Span) { return nil, nil }

func (z *ZapRepo) End() { return }
