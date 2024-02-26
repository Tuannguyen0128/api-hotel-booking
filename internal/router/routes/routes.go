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

func Load() []Route {
	routes := make([]Route, 0)
	accountRouter = append(accountRouter, staffRouter...)
	accountRouter = append(accountRouter, bookingRouter...)
	return routes
}
func SetupRoutesWithMiddleWares(g *gin.Engine) {
	for _, router := range Load() {
		if router.AuthRequired == true {
			g.Handle(router.Method, "/api"+router.Uri,
				middlewares.SetMiddleWareLogger(),
				middlewares.SetMiddleWareJSON(),
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
