package honeycomb

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"os"
	"runtime"
)

type HoneycombRepository struct {
	c    *gin.Context
	span trace.Span
}

func NewHoneycombRepository() *HoneycombRepository {
	return &HoneycombRepository{nil, nil}
}

// Start Honeycomb
func (h *HoneycombRepository) Start(c *gin.Context, info string) (*gin.Context, trace.Span) {
	_, file, line, _ := runtime.Caller(2)

	tracer := otel.Tracer("")
	ctx, span := tracer.Start(c.Request.Context(), info, trace.WithAttributes(
		attribute.String("file", file),
		attribute.String("client_ip", c.ClientIP()),
		attribute.Int("line", line),
	))
	h.span = span
	temp := c.Copy()
	temp.Request = c.Request.WithContext(ctx)
	h.c = temp
	return temp, span
}

func (h *HoneycombRepository) End() {
	h.span.End()
	return
}

func createTraceAttributes(level, jsonData, file string, line int, clientIP string) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("level", level),
		attribute.String("data", jsonData),
		attribute.String("file", file),
		attribute.String("client_ip", clientIP),
		attribute.Int("line", line),
	}
}

// Debug logs a debug message
func (h *HoneycombRepository) Debug(msg string, fields map[string]interface{}) {
	var jsonData []byte
	var jsonDataErr error

	if fields != nil {
		jsonData, jsonDataErr = json.Marshal(fields)
		if jsonDataErr != nil {
			return
		}
	}

	_, file, line, _ := runtime.Caller(2)

	attrs := createTraceAttributes("Debug", string(jsonData), file, line, h.c.ClientIP())
	h.span.AddEvent(msg, trace.WithAttributes(attrs...))
}

// Info logs an info message
func (h *HoneycombRepository) Info(msg string, fields map[string]interface{}) {
	var jsonData []byte
	var jsonDataErr error

	if fields != nil {
		jsonData, jsonDataErr = json.Marshal(fields)
		if jsonDataErr != nil {
			return
		}
	}
	_, file, line, _ := runtime.Caller(2)

	attrs := createTraceAttributes("Info", string(jsonData), file, line, h.c.ClientIP())
	h.span.AddEvent(msg, trace.WithAttributes(attrs...))
}

func (h *HoneycombRepository) Error(msg string, fields map[string]interface{}) {
	var jsonData []byte
	var jsonDataErr error

	if fields != nil {
		jsonData, jsonDataErr = json.Marshal(fields)
		if jsonDataErr != nil {
			return
		}
	}
	_, file, line, _ := runtime.Caller(2)

	attrs := createTraceAttributes("Error", string(jsonData), file, line, h.c.ClientIP())
	h.span.RecordError(errors.New(msg), trace.WithAttributes(attrs...))
	h.span.SetStatus(codes.Error, msg)
}

// Warn logs a warning message
func (h *HoneycombRepository) Warn(msg string, fields map[string]interface{}) {
	var jsonData []byte
	var jsonDataErr error

	if fields != nil {
		jsonData, jsonDataErr = json.Marshal(fields)
		if jsonDataErr != nil {
			return
		}
	}

	_, file, line, _ := runtime.Caller(2)

	attrs := createTraceAttributes("Warn", string(jsonData), file, line, h.c.ClientIP())
	h.span.AddEvent(msg, trace.WithAttributes(attrs...))
}

// Fatal logs a fatal message
func (h *HoneycombRepository) Fatal(msg string, fields map[string]interface{}) {
	// Implement fatal logging with Honeycomb
	// This function can be implemented similar to the Info function
	var jsonData []byte
	var jsonDataErr error

	if fields != nil {
		jsonData, jsonDataErr = json.Marshal(fields)
		if jsonDataErr != nil {
			return
		}
	}

	_, file, line, _ := runtime.Caller(2)

	attrs := createTraceAttributes("Fatal", string(jsonData), file, line, h.c.ClientIP())
	h.span.AddEvent(msg, trace.WithAttributes(attrs...))
	os.Exit(1)
}
