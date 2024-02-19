package main

import (
	"api-hotel-booking/internal/app"
	"api-hotel-booking/internal/app/handler"
	"api-hotel-booking/internal/app/server"
	"api-hotel-booking/internal/app/service"
	"api-hotel-booking/internal/app/utils"
	"api-hotel-booking/internal/app/utils/header/kbheader"
	"api-hotel-booking/internal/app/utils/header/pheader"

	"api-hotel-booking/internal/app/thirdparty/config"
	"api-hotel-booking/internal/app/thirdparty/logger"
)

// LWS_MONGODB_PASSWORD=dev;LWS_MONGODB_SCHEMA_MAIN=labellingWebService;LWS_MONGODB_AUTH_SOURCE=labellingWebService

func main() {
	const configPath = "./configs"
	var logConfig logger.Config
	if err := config.LoadConfigs("log", configPath, "log", &logConfig); err != nil {
		logger.Panic("failed to load log config: %v", err)
	}
	logger.InitLogInst(logConfig)
	defer logger.Flush()

	if err := config.LoadConfigs("LWS", configPath, "app", app.CFG); err != nil {
		logger.Panic("error while reading app configs: %v", err)
	}
	logger.Info("%s", app.CFG.Environment)

	if err := config.LoadConfigs("err", configPath, "error", app.ErrorMap); err != nil {
		logger.Panic("error while mapping error file: %v", err)
	}

	a := service.NewAdapters()
	s := service.NewService(a)
	kbheader.InitKBankHeader(app.CFG.AppID, app.CFG.AppAbbr, app.CFG.IsDebug)
	pheader.Init(app.CFG.IsDebug)
	engine := handler.Init(app.CFG.Server)
	handler.RegisterRoutes(engine, s, app.CFG.IsDebug)
	svr := server.New(app.CFG.Server, engine, utils.OnClose)
	logger.Info("Start serving at port %d", app.CFG.Server.Port)
	if err := svr.ListenAndServe(app.CFG.Server.GracefulShutdownTime, app.CFG.Server.Https); err != nil {
		logger.Panic("error while starting app at port %d: %v", app.CFG.Server.Port, err)
	}

	logger.Info("Stop serving at port %d", app.CFG.Server.Port)
}
