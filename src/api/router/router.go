package router

import (
	"ProjectPractice/src/api/router/routes"
	"github.com/gin-gonic/gin"
)
//
//func New() {
//
//	r := Init()
//	return
//}

func Init() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	//engine.Use(gin.Recovery(), cors.New(cfg.CORS))
	//engine.MaxMultipartMemory = app.CFG.Service.RequestSize
	routes.SetupRoutesWithMiddleWares(engine)

	return engine
}
