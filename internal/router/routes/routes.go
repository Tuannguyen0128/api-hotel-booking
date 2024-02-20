package routes

import (
	"api-hotel-booking/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type Route struct {
	Uri          string
	Method       string
	Handler      func(c *gin.Context)
	AuthRequired bool
}

//	func SetupRoutes(c *gin.Engine) {
//		for _, router := range Load() {
//			c.Handle(router.Method,router.Uri, router.Handler)
//		}
//	}
func Load() []Route {
	routes := append(userRouter, merchantroutes...)
	routes = append(routes, loginRoutes...)
	return routes
}
func SetupRoutesWithMiddleWares(g *gin.Engine) {
	for _, router := range Load() {
		if router.AuthRequired == true {
			g.Handle(router.Method, "/api"+router.Uri,
				middlewares.SetMiddleWareLogger(),
				middlewares.SetMiddleWareJSON(),
				middlewares.SetMiddleWareAuthentication(),
				router.Handler,
			)
		} else {
			g.Handle(router.Method, "/api"+router.Uri,
				middlewares.SetMiddleWareLogger(),
				middlewares.SetMiddleWareJSON(),
				router.Handler)
		}
	}
}
