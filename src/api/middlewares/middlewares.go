package middlewares

import (
	"ProjectPractice/src/api/auth"
	"ProjectPractice/src/api/responses"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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

func SetMiddleWareAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := auth.ExtractToken(c)
		log.Println("Token:" + token)
		err := auth.TokenValidate(token)
		if err != nil {
			log.Println(err)
			responses.ERROR( http.StatusUnauthorized, err.Error())
			return
		}

		c.Next()
	}
}
