package routes

import (
	"api-hotel-booking/internal/controller"
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
