package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

const timeFormatRFC3339Milli = "2006-01-02T15:04:05.000Z07:00"

type (
	logger interface {
		Print(level string, fields map[string]interface{}, format string, args ...interface{})
		GetRequestID(ctx context.Context) string
	}
)

func RequestLogger(log logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		fields := make(map[string]interface{})
		if reqID := log.GetRequestID(ctx); reqID != "" {
			fields["request_id"] = reqID
		}
		scheme := "http"
		if ctx.Request.TLS != nil {
			scheme = "https"
		}
		fields["http_scheme"] = scheme
		fields["http_proto"] = ctx.Request.Proto
		fields["http_method"] = ctx.Request.Method
		fields["user_agent"] = ctx.Request.UserAgent()
		fields["request_time"] = start.Format(timeFormatRFC3339Milli)
		fields["uri"] = fmt.Sprintf("%s://%s%s", scheme, ctx.Request.Host, ctx.Request.RequestURI)
		log.Print("info", fields, "request started")
		defer func() {
			elapsedTime := time.Since(start)
			fields["response_status"] = ctx.Writer.Status()
			fields["response_elapsed_ms"] = elapsedTime.Milliseconds()
			log.Print("info", fields, "request completed with status code [%d] after %vms", ctx.Writer.Status(), fields["response_elapsed_ms"])
		}()
		ctx.Next()
	}
}
