package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"api-hotel-booking/internal/app/service/user"

	"api-hotel-booking/internal/app"
	"api-hotel-booking/internal/app/server"
	"api-hotel-booking/internal/app/service/session"
	"api-hotel-booking/internal/app/thirdparty/go-gin/middleware"
	"api-hotel-booking/internal/app/thirdparty/go-gin/middleware/logstats"
	"api-hotel-booking/internal/app/thirdparty/logger"
	"api-hotel-booking/internal/app/utils"
	"api-hotel-booking/internal/app/utils/header/pheader"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var log = logger.WithModule("handler")

type (
	route struct {
		URL     string `json:"url"`
		Method  string `json:"method"`
		Enabled bool   `json:"enabled"`
		Handler gin.HandlerFunc
	}

	Service interface {
		Login(ctx context.Context, req session.LoginRequest) (session.LoginResponse, error)
		ExtendSession(ctx context.Context, req session.ExtendRequest) (session.ExtendResponse, error)
		Logout(req session.LogoutRequest) error
		CreateUser(ctx context.Context, req user.CreateRequest) (user.CreateResponse, error)
		EditUser(ctx context.Context, req user.EditRequest) (user.EditResponse, error)
		GetUser(ctx context.Context, req user.GetRequest) (user.GetResponse, error)
		DeleteUser(ctx context.Context, req user.DeleteRequest) (user.DeleteResponse, error)
		ChangePassword(ctx context.Context, req user.ChangePassRequest) (user.ChangePassResponse, error)
		ListUser(ctx context.Context, req user.ListRequest) (user.ListResponse, error)
		ForgetPassword(ctx context.Context, req user.ForgetPassRequest) error
		ResetPassword(ctx context.Context, req user.ResetPassRequest) (user.ResetPassResponse, error)
	}
)

func Init(cfg server.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	engine.Use(gin.Recovery(), cors.New(cfg.CORS))
	engine.MaxMultipartMemory = app.CFG.Service.RequestSize

	return engine
}

func RegisterRoutes(e *gin.Engine, s Service, isDebug bool) {
	healthCheckPath := app.CFG.Service.Prefix + "/healthz"
	route := e.Group(
		"",
		middleware.HealthCheck(healthCheckPath),
		middleware.NoCache,
		middleware.RequestID,
		logstats.RequestLoggingMiddleware(log, "PARTNER_LICENSE", []string{http.MethodGet + " " + healthCheckPath}),
	)
	logger.ImportRequestIDFunction(middleware.GetRequestID)
	log.Info("register url Get %s", healthCheckPath)
	route.GET(healthCheckPath)

	doRegisterSubGroup(route, app.CFG.Service.Prefix+"/web", WebRoutesV1(app.CFG.Service.API, s, isDebug))
}

func doRegisterSubGroup(g *gin.RouterGroup, subGroup string, routes map[string]route) {
	sg := g.Group(subGroup)
	for _, r := range routes {
		if r.Enabled {
			log.Info("register url %s %s%s", r.Method, sg.BasePath(), r.URL)
			sg.Handle(
				r.Method,
				r.URL,
				r.Handler,
			)
		} else {
			log.Info("url %s %s%s is disabled", r.Method, sg.BasePath(), r.URL)
			sg.Handle(
				r.Method,
				r.URL,
				func(ctx *gin.Context) {
					utils.AbortWithError(ctx, app.ErrorMap.FeatureDisabled)
				},
			)
		}
	}
}

func WebRoutesV1(cfg app.EndPoints, s Service, isDebug bool) map[string]route {
	return map[string]route{
		"Login": {
			URL:     cfg.Web.V1.Login.URL,
			Method:  cfg.Web.V1.Login.Method,
			Enabled: cfg.Web.V1.Login.Enabled,
			Handler: func(c *gin.Context) {
				var req session.LoginRequest
				if e := c.ShouldBindJSON(&req); e != nil {
					bJson, err := json.Marshal(req)
					if err != nil {
						log.KWarn(c, "cannot marshal empty req to log out err:"+err.Error())
					}
					re := app.ErrorMap.InvalidRequest.AddDebug("need structure req:" + string(bJson))
					log.KError(c, "response error: %s", re.SLog())
					c.JSON(re.HttpCode, pheader.GetResponseError(pheader.Request{}, re))
				} else if res, err := s.Login(c, req); err != nil {
					log.KError(c, "response error: %s", app.GetResponseError(err, true).SLog())
					re := app.GetResponseError(err, isDebug)
					c.JSON(re.HttpCode, pheader.GetResponseError(req.Header, re))

				} else {
					res.Header = pheader.GetResponseHeader(req.Header, nil)
					c.JSON(http.StatusOK, res)
				}

				log.KInfo(c, "Client request id:%s Client Corr id:not have", req.Header.RequestId)
				logstats.SetPartnerId(c, req.Header.PartnerId)
				logstats.SetAPIName(c, "Login")
			},
		},
		"ExtendSession": {
			URL:     cfg.Web.V1.ExtendSession.URL,
			Method:  cfg.Web.V1.ExtendSession.Method,
			Enabled: cfg.Web.V1.ExtendSession.Enabled,
			Handler: func(c *gin.Context) {
				var req session.ExtendRequest
				if e := c.ShouldBindJSON(&req); e != nil {
					bJson, err := json.Marshal(req)
					if err != nil {
						log.KWarn(c, "cannot marshal empty req to log out err:"+err.Error())
					}
					re := app.ErrorMap.InvalidRequest.AddDebug("need structure req:" + string(bJson))
					log.KError(c, "response error: %s", re.SLog())
					c.JSON(re.HttpCode, pheader.GetResponseError(pheader.Request{}, re))
				} else if res, err := s.ExtendSession(c, req); err != nil {
					log.KError(c, "response error: %s", app.GetResponseError(err, true).SLog())
					re := app.GetResponseError(err, isDebug)
					c.JSON(re.HttpCode, pheader.GetResponseError(req.Header, re))

				} else {
					res.Header = pheader.GetResponseHeader(req.Header, nil)
					c.JSON(http.StatusOK, res)
				}

				log.KInfo(c, "Client request id:%s Client Corr id:not have", req.Header.RequestId)
				logstats.SetPartnerId(c, req.Header.PartnerId)
				logstats.SetAPIName(c, "ExtendSession")
			},
		},
		"Logout": {
			URL:     cfg.Web.V1.Logout.URL,
			Method:  cfg.Web.V1.Logout.Method,
			Enabled: cfg.Web.V1.Logout.Enabled,
			Handler: func(c *gin.Context) {
				var req session.LogoutRequest
				if e := c.ShouldBindJSON(&req); e != nil {
					bJson, err := json.Marshal(req)
					if err != nil {
						log.KWarn(c, "cannot marshal empty req to log out err:"+err.Error())
					}
					re := app.ErrorMap.InvalidRequest.AddDebug("need structure req:" + string(bJson))
					log.KError(c, "response error: %s", re.SLog())
					c.JSON(re.HttpCode, pheader.GetResponseError(pheader.Request{}, re))
				} else if err := s.Logout(req); err != nil {
					log.KError(c, "response error: %s", app.GetResponseError(err, true).SLog())
					re := app.GetResponseError(err, isDebug)
					c.JSON(re.HttpCode, pheader.GetResponseError(req.Header, re))

				} else {
					res := pheader.GetResponseError(req.Header, nil)
					c.JSON(http.StatusOK, res)
				}

				log.KInfo(c, "Client request id:%s Client Corr id:not have", req.Header.RequestId)
				logstats.SetPartnerId(c, req.Header.PartnerId)
				logstats.SetAPIName(c, "Logout")
			},
		},
		"CreateUser": {
			URL:     cfg.Web.V1.CreateUser.URL,
			Method:  cfg.Web.V1.CreateUser.Method,
			Enabled: cfg.Web.V1.CreateUser.Enabled,
			Handler: func(c *gin.Context) {
				var req user.CreateRequest
				if e := c.ShouldBindJSON(&req); e != nil {
					bJson, err := json.Marshal(req)
					if err != nil {
						log.KWarn(c, "cannot marshal empty req to log out err:"+err.Error())
					}
					re := app.ErrorMap.InvalidRequest.AddDebug("need structure req:" + string(bJson))
					log.KError(c, "response error: %s", re.SLog())
					c.JSON(re.HttpCode, pheader.GetResponseError(pheader.Request{}, re))
				} else if res, err := s.CreateUser(c, req); err != nil {
					log.KError(c, "response error: %s", app.GetResponseError(err, true).SLog())
					re := app.GetResponseError(err, isDebug)
					c.JSON(re.HttpCode, pheader.GetResponseError(req.Header, re))

				} else {
					res.Header = pheader.GetResponseHeader(req.Header, nil)
					c.JSON(http.StatusOK, res)
				}

				log.KInfo(c, "Client request id:%s Client Corr id:not have", req.Header.RequestId)
				logstats.SetPartnerId(c, req.Header.PartnerId)
				logstats.SetAPIName(c, "CreateUser")
			},
		},
		"EditUser": {
			URL:     cfg.Web.V1.EditUser.URL,
			Method:  cfg.Web.V1.EditUser.Method,
			Enabled: cfg.Web.V1.EditUser.Enabled,
			Handler: func(c *gin.Context) {
				var req user.EditRequest
				if e := c.ShouldBindJSON(&req); e != nil {
					bJson, err := json.Marshal(req)
					if err != nil {
						log.KWarn(c, "cannot marshal empty req to log out err:"+err.Error())
					}
					re := app.ErrorMap.InvalidRequest.AddDebug("need structure req:" + string(bJson))
					log.KError(c, "response error: %s", re.SLog())
					c.JSON(re.HttpCode, pheader.GetResponseError(pheader.Request{}, re))
				} else if res, err := s.EditUser(c, req); err != nil {
					log.KError(c, "response error: %s", app.GetResponseError(err, true).SLog())
					re := app.GetResponseError(err, isDebug)
					c.JSON(re.HttpCode, pheader.GetResponseError(req.Header, re))

				} else {
					res.Header = pheader.GetResponseHeader(req.Header, nil)
					c.JSON(http.StatusOK, res)
				}

				log.KInfo(c, "Client request id:%s Client Corr id:not have", req.Header.RequestId)
				logstats.SetPartnerId(c, req.Header.PartnerId)
				logstats.SetAPIName(c, "EditUser")
			},
		},
		"GetUser": {
			URL:     cfg.Web.V1.GetUser.URL,
			Method:  cfg.Web.V1.GetUser.Method,
			Enabled: cfg.Web.V1.GetUser.Enabled,
			Handler: func(c *gin.Context) {
				var req user.GetRequest
				if e := c.ShouldBindJSON(&req); e != nil {
					bJson, err := json.Marshal(req)
					if err != nil {
						log.KWarn(c, "cannot marshal empty req to log out err:"+err.Error())
					}
					re := app.ErrorMap.InvalidRequest.AddDebug("need structure req:" + string(bJson))
					log.KError(c, "response error: %s", re.SLog())
					c.JSON(re.HttpCode, pheader.GetResponseError(pheader.Request{}, re))
				} else if res, err := s.GetUser(c, req); err != nil {
					log.KError(c, "response error: %s", app.GetResponseError(err, true).SLog())
					re := app.GetResponseError(err, isDebug)
					c.JSON(re.HttpCode, pheader.GetResponseError(req.Header, re))

				} else {
					res.Header = pheader.GetResponseHeader(req.Header, nil)
					c.JSON(http.StatusOK, res)
				}

				log.KInfo(c, "Client request id:%s Client Corr id:not have", req.Header.RequestId)
				logstats.SetPartnerId(c, req.Header.PartnerId)
				logstats.SetAPIName(c, "GetUser")
			},
		},
		"ChangePassword": {
			URL:     cfg.Web.V1.ChangePassword.URL,
			Method:  cfg.Web.V1.ChangePassword.Method,
			Enabled: cfg.Web.V1.ChangePassword.Enabled,
			Handler: func(c *gin.Context) {
				var req user.ChangePassRequest
				if e := c.ShouldBindJSON(&req); e != nil {
					bJson, err := json.Marshal(req)
					if err != nil {
						log.KWarn(c, "cannot marshal empty req to log out err:"+err.Error())
					}
					re := app.ErrorMap.InvalidRequest.AddDebug("need structure req:" + string(bJson))
					log.KError(c, "response error: %s", re.SLog())
					c.JSON(re.HttpCode, pheader.GetResponseError(pheader.Request{}, re))
				} else if res, err := s.ChangePassword(c, req); err != nil {
					log.KError(c, "response error: %s", app.GetResponseError(err, true).SLog())
					re := app.GetResponseError(err, isDebug)
					c.JSON(re.HttpCode, pheader.GetResponseError(req.Header, re))

				} else {
					res.Header = pheader.GetResponseHeader(req.Header, nil)
					c.JSON(http.StatusOK, res)
				}

				log.KInfo(c, "Client request id:%s Client Corr id:not have", req.Header.RequestId)
				logstats.SetPartnerId(c, req.Header.PartnerId)
				logstats.SetAPIName(c, "ChangePassword")
			},
		},
		"ListUser": {
			URL:     cfg.Web.V1.ListUser.URL,
			Method:  cfg.Web.V1.ListUser.Method,
			Enabled: cfg.Web.V1.ListUser.Enabled,
			Handler: func(c *gin.Context) {
				var req user.ListRequest
				if e := c.ShouldBindJSON(&req); e != nil {
					bJson, err := json.Marshal(req)
					if err != nil {
						log.KWarn(c, "cannot marshal empty req to log out err:"+err.Error())
					}
					re := app.ErrorMap.InvalidRequest.AddDebug("need structure req:" + string(bJson))
					log.KError(c, "response error: %s", re.SLog())
					c.JSON(re.HttpCode, pheader.GetResponseError(pheader.Request{}, re))
				} else if res, err := s.ListUser(c, req); err != nil {
					log.KError(c, "response error: %s", app.GetResponseError(err, true).SLog())
					re := app.GetResponseError(err, isDebug)
					c.JSON(re.HttpCode, pheader.GetResponseError(req.Header, re))

				} else {
					res.Header = pheader.GetResponseHeader(req.Header, nil)
					c.JSON(http.StatusOK, res)
				}

				log.KInfo(c, "Client request id:%s Client Corr id:not have", req.Header.RequestId)
				logstats.SetPartnerId(c, req.Header.PartnerId)
				logstats.SetAPIName(c, "ListUser")
			},
		},
		"DeleteUser": {
			URL:     cfg.Web.V1.DeleteUser.URL,
			Method:  cfg.Web.V1.DeleteUser.Method,
			Enabled: cfg.Web.V1.DeleteUser.Enabled,
			Handler: func(c *gin.Context) {
				var req user.DeleteRequest
				if e := c.ShouldBindJSON(&req); e != nil {
					bJson, err := json.Marshal(req)
					if err != nil {
						log.KWarn(c, "cannot marshal empty req to log out err:"+err.Error())
					}
					re := app.ErrorMap.InvalidRequest.AddDebug("need structure req:" + string(bJson))
					log.KError(c, "response error: %s", re.SLog())
					c.JSON(re.HttpCode, pheader.GetResponseError(pheader.Request{}, re))
				} else if res, err := s.DeleteUser(c, req); err != nil {
					log.KError(c, "response error: %s", app.GetResponseError(err, true).SLog())
					re := app.GetResponseError(err, isDebug)
					c.JSON(re.HttpCode, pheader.GetResponseError(req.Header, re))

				} else {
					res.Header = pheader.GetResponseHeader(req.Header, nil)
					c.JSON(http.StatusOK, res)
				}

				log.KInfo(c, "Client request id:%s Client Corr id:not have", req.Header.RequestId)
				logstats.SetPartnerId(c, req.Header.PartnerId)
				logstats.SetAPIName(c, "DeleteUser")
			},
		},
		"ForgetPassword": {
			URL:     cfg.Web.V1.ForgetPassword.URL,
			Method:  cfg.Web.V1.ForgetPassword.Method,
			Enabled: cfg.Web.V1.ForgetPassword.Enabled,
			Handler: func(c *gin.Context) {
				var req user.ForgetPassRequest
				if e := c.ShouldBindJSON(&req); e != nil {
					bJson, err := json.Marshal(req)
					if err != nil {
						log.KWarn(c, "cannot marshal empty req to log out err:"+err.Error())
					}
					re := app.ErrorMap.InvalidRequest.AddDebug("need structure req:" + string(bJson))
					log.KError(c, "response error: %s", re.SLog())
					c.JSON(re.HttpCode, pheader.GetResponseError(pheader.Request{}, re))
				} else if err := s.ForgetPassword(c, req); err != nil {
					log.KError(c, "response error: %s", app.GetResponseError(err, true).SLog())

					// this api never return error
					res := pheader.GetResponseError(req.Header, nil)
					c.JSON(http.StatusOK, res)

				} else {
					res := pheader.GetResponseError(req.Header, nil)
					c.JSON(http.StatusOK, res)
				}

				log.KInfo(c, "Client request id:%s Client Corr id:not have", req.Header.RequestId)
				logstats.SetPartnerId(c, req.Header.PartnerId)
				logstats.SetAPIName(c, "ForgetPassword")
			},
		},
		"ResetPassword": {
			URL:     cfg.Web.V1.ResetPassword.URL,
			Method:  cfg.Web.V1.ResetPassword.Method,
			Enabled: cfg.Web.V1.ResetPassword.Enabled,
			Handler: func(c *gin.Context) {
				var req user.ResetPassRequest
				if e := c.ShouldBindJSON(&req); e != nil {
					bJson, err := json.Marshal(req)
					if err != nil {
						log.KWarn(c, "cannot marshal empty req to log out err:"+err.Error())
					}
					re := app.ErrorMap.InvalidRequest.AddDebug("need structure req:" + string(bJson))
					log.KError(c, "response error: %s", re.SLog())
					c.JSON(re.HttpCode, pheader.GetResponseError(pheader.Request{}, re))
				} else if res, err := s.ResetPassword(c, req); err != nil {
					log.KError(c, "response error: %s", app.GetResponseError(err, true).SLog())
					re := app.GetResponseError(err, isDebug)
					c.JSON(re.HttpCode, pheader.GetResponseError(req.Header, re))

				} else {
					res.Header = pheader.GetResponseHeader(req.Header, nil)
					c.JSON(http.StatusOK, res)
				}

				log.KInfo(c, "Client request id:%s Client Corr id:not have", req.Header.RequestId)
				logstats.SetPartnerId(c, req.Header.PartnerId)
				logstats.SetAPIName(c, "ForgetPassword")
			},
		},
	}
}
