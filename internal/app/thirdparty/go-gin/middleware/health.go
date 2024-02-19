package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func HealthCheck(endpoint string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if r := ctx.Request; r.Method == http.MethodGet && strings.EqualFold(r.URL.Path, endpoint) {
			ctx.Abort()
			ctx.String(http.StatusOK, ".")
			return
		}
	}
}
