package routes

import (
	"api-hotel-booking/internal/controller"
	"net/http"
)

var accountRouter = []Route{
	{
		Uri:          "/accounts",
		Method:       http.MethodGet,
		Handler:      controller.GetAccounts,
		AuthRequired: true,
	},
	{
		Uri:          "/account",
		Method:       http.MethodPost,
		Handler:      controller.CreateAccount,
		AuthRequired: true,
	},
	{
		Uri:          "/account/:id",
		Method:       http.MethodPut,
		Handler:      controller.UpdateAccount,
		AuthRequired: true,
	},
	{
		Uri:          "/account/:id",
		Method:       http.MethodDelete,
		Handler:      controller.DeleteAccount,
		AuthRequired: true,
	},
}
