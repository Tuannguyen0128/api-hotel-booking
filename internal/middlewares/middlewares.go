package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
)

func SetMiddleWareLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("%s %s%s %s", c.Request.Method, c.Request.Host, c.Request.RequestURI, c.Request.Proto)
		c.Next()
	}
}
func SetMiddleWareJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Next()
	}
}
