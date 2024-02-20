package routes

import (
	"ProjectPractice/src/api/controller"
	"net/http"
)

var loginRoutes = []Route{
	{
		Uri:          "/login",
		Method:       http.MethodPost,
		Handler:      controller.Login,
		AuthRequired: false,
	},
}
