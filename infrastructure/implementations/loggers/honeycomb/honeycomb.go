package honeycomb

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"log"
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
func (h *HoneycombRepository) Start(c *gin.Context, info string) trace.Span {
	h.c = c
	_, file, line, _ := runtime.Caller(2)
	// Retrieve otel_context from Gin context if it exists
	ctxValue, exists := c.Get("otel_context")
	var ctx context.Context
	if exists {
		ctx, _ = ctxValue.(context.Context)
	} else {
		log.Println("Getting this from" + info)
		ctx = c.Request.Context()
	}
	tracer := otel.Tracer("")
	_, span := tracer.Start(ctx, info, trace.WithAttributes(
		attribute.String("file", file),
		attribute.String("client_ip", h.c.ClientIP()),
		attribute.Int("line", line),
	))

	h.span = span

	return span
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

func (h *HoneycombRepository) GetSpan() trace.Span {
	return h.span
}

func (h *HoneycombRepository) GetContext() *gin.Context {
	return h.c
}

func (h *HoneycombRepository) UseGivenSpan(span trace.Span) {
	h.span = span
}