package logger

import (
	"context"
	"fmt"
)

type (
	// Log with context input, support to use with go-chi and go-gin
	// It can print log line in json with attached request_id. See ImportRequestIDFunction
	ContextLogger interface {
		KDebug(ctx context.Context, format string, args ...interface{})
		KInfo(ctx context.Context, format string, args ...interface{})
		KWarn(ctx context.Context, format string, args ...interface{})
		KError(ctx context.Context, format string, args ...interface{})

		DDebug(ctx context.Context, data interface{}, format string, args ...interface{})
		DInfo(ctx context.Context, data interface{}, format string, args ...interface{})
		DWarn(ctx context.Context, data interface{}, format string, args ...interface{})
		DError(ctx context.Context, data interface{}, format string, args ...interface{})

		GetRequestID(ctx context.Context) string
	}

	// Help to customize the output log of object of DDebug/DInfo/DWarn/DError function in ContextLogger
	// Can be used to strip out unwanted/sensitive data from the log
	// See logJson function
	JsonCustomizer interface {
		Customize() interface{}
	}
)

var requestIDFn func(ctx context.Context) string

// Add the logic function to get request id from input context. Since the context may different from other web framework. See api-hotel-booking/internal/app/thirdparty/chi_middleware or api-hotel-booking/internal/app/thirdparty/gin_middleware
func ImportRequestIDFunction(fn func(ctx context.Context) string) {
	if fn != nil {
		requestIDFn = fn
	}
}

func (e *entry) KDebug(ctx context.Context, format string, args ...interface{}) {
	e.contextLog(e.inst.Debugw, ctx, nil, format, args...)
}

func (e *entry) KInfo(ctx context.Context, format string, args ...interface{}) {
	e.contextLog(e.inst.Infow, ctx, nil, format, args...)
}

func (e *entry) KWarn(ctx context.Context, format string, args ...interface{}) {
	e.contextLog(e.inst.Warnw, ctx, nil, format, args...)
}

func (e *entry) KError(ctx context.Context, format string, args ...interface{}) {
	e.contextLog(e.inst.Errorw, ctx, nil, format, args...)
}

func (e *entry) DDebug(ctx context.Context, data interface{}, format string, args ...interface{}) {
	e.contextLog(e.inst.Debugw, ctx, data, format, args...)
}

func (e *entry) DInfo(ctx context.Context, data interface{}, format string, args ...interface{}) {
	e.contextLog(e.inst.Infow, ctx, data, format, args...)
}

func (e *entry) DWarn(ctx context.Context, data interface{}, format string, args ...interface{}) {
	e.contextLog(e.inst.Warnw, ctx, data, format, args...)
}

func (e *entry) DError(ctx context.Context, data interface{}, format string, args ...interface{}) {
	e.contextLog(e.inst.Errorw, ctx, data, format, args...)
}

// Get RequestID from the input context. ImportRequestIDFunction need to call first with proper fn
func (e *entry) GetRequestID(ctx context.Context) string {
	if requestIDFn != nil {
		return requestIDFn(ctx)
	}
	return ""
}

func (e *entry) contextLog(fn LogFieldsFn, ctx context.Context, data interface{}, template string, args ...interface{}) {
	var keyAndValues []interface{}
	if reqID := e.GetRequestID(ctx); reqID != "" {
		keyAndValues = append(keyAndValues, "request_id", reqID)
	}
	if data != nil {
		keyAndValues = append(keyAndValues, "object_data", logJson(data))
	}
	fn(fmt.Sprintf(template, args...), keyAndValues...)
}

// Custom output of object that implements JsonCustomizer interface
func logJson(v interface{}) interface{} {
	if data, ok := v.(JsonCustomizer); ok {
		return data.Customize()
	}
	return v
}
