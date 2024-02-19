package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-Id"

// Key to use when setting the request ID.
type ctxKeyRequestID int

// RequestIDKey is the key that holds the unique request ID in a request context.
const RequestIDKey ctxKeyRequestID = 0

func RequestID(ctx *gin.Context) {
	reqID := ctx.GetHeader(RequestIDHeader)
	if reqID == "" {
		reqID = uuid.New().String()
	}
	ctx.Set(RequestIDHeader, reqID)

	// Attach requestID into request Context
	requestCtx := context.WithValue(ctx.Request.Context(), RequestIDKey, reqID)
	ctx.Request = ctx.Request.WithContext(requestCtx)
}

// Function to get attached requestID attached to context. Can be use for logger
// See api-hotel-booking/internal/app/thirdparty/logger with logger.ImportRequestIDFunction function
func GetRequestID(ctx context.Context) string {
	if c, ok := ctx.(*gin.Context); ok {
		return c.GetString(RequestIDHeader)
	}
	return ctx.Value(RequestIDKey).(string)
}
