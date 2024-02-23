package router

import (
	"api-hotel-booking/internal/router/routes"

	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	routes.SetupRoutesWithMiddleWares(engine)

	return engine
}
