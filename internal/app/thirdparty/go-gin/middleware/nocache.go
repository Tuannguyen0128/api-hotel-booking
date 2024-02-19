package middleware

import (
	"github.com/gin-gonic/gin"
	"time"
)

// Unix epoch time
var epoch = time.Unix(0, 0).Format(time.RFC1123)

// Taken from https://github.com/mytrile/nocache
var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, no-store, no-transform, must-revalidate, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

var etagHeaders = []string{
	"ETag",
	"If-Modified-Since",
	"If-Match",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

func NoCache(ctx *gin.Context) {
	for _, v := range etagHeaders {
		if ctx.GetHeader(v) != "" {
			ctx.Request.Header.Del(v)
		}
	}

	for k, v := range noCacheHeaders {
		ctx.Header(k, v)
	}
	ctx.Next()
}
