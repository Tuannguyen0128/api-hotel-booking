package logstats

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type (
	logger interface {
		Print(level string, fields map[string]interface{}, format string, args ...interface{})
		GetRequestID(ctx context.Context) string
	}
)

// RequestLoggingMiddleware skipAPIPath in format 'GET /health' (upper case method (space) url
func RequestLoggingMiddleware(log logger, serviceID string, skipAPIPath []string) gin.HandlerFunc {
	skipAPIPathMap := make(map[string]bool)
	for i := range skipAPIPath {
		skipAPIPathMap[skipAPIPath[i]] = true
	}
	return func(ctx *gin.Context) {
		if skipAPIPathMap[fmt.Sprintf("%s %s", ctx.Request.Method, ctx.Request.RequestURI)] {
			ctx.Next()
			return
		}
		start := time.Now()
		fields := make(map[string]interface{})
		fields["request_id"] = getRequestId(ctx, log)
		scheme := "http"
		if ctx.Request.TLS != nil {
			scheme = "https"
		}
		fields["service_id"] = serviceID
		fields["http_scheme"] = scheme
		fields["http_proto"] = ctx.Request.Proto
		fields["http_method"] = ctx.Request.Method
		fields["api_path"] = ctx.Request.RequestURI
		fields["user_agent"] = ctx.Request.UserAgent()
		fields["request_time"] = start.Format(timeFormatRFC3339Milli)
		fields["uri"] = fmt.Sprintf("%s://%s%s", scheme, ctx.Request.Host, ctx.Request.RequestURI)
		log.Print("info", fields, "request started")
		defer func() {
			fields["log_type"] = "log-stats"
			fields["caller_service_id"] = getCallerServiceId(ctx)
			fields["partner_id"] = getPartnerId(ctx)
			fields["api_name"] = getAPIName(ctx)
			fields["return_http_status"] = ctx.Writer.Status()
			fields["remark"] = getRemark(ctx)
			fields["addition_info"] = getAdditionInfo(ctx)
			responseTime := time.Now()
			fields["response_time"] = responseTime.Format(timeFormatRFC3339Milli)
			fields["response_elapsed_ms"] = responseTime.Sub(start).Milliseconds()
			log.Print("info", fields, "request completed with status code [%d] after %vms", ctx.Writer.Status(), fields["response_elapsed_ms"])
		}()
		ctx.Next()
	}
}
