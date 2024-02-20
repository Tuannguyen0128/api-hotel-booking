package routes

import (
	"api-hotel-booking/internal/controller"
	"net/http"
)

var userRouter = []Route{
	{
		Uri:          "/accounts",
		Method:       http.MethodGet,
		Handler:      controller.GetAccount,
		AuthRequired: true,
	},
	{
		Uri:          "/account/{id}",
		Method:       http.MethodGet,
		Handler:      controller.GetAccount,
		AuthRequired: true,
	},
	{
		Uri:          "/account",
		Method:       http.MethodPost,
		Handler:      controller.CreateAccount,
		AuthRequired: false,
	},
	{
		Uri:          "/account/{id}",
		Method:       http.MethodPut,
		Handler:      controller.UpdateAccount,
		AuthRequired: true,
	},
	{
		Uri:          "/account/{id}",
		Method:       http.MethodDelete,
		Handler:      controller.DeleteAccount,
		AuthRequired: true,
	},
	{
		Uri:          "/get-account-by-email",
		Method:       http.MethodGet,
		Handler:      controller.FindByEmail,
		AuthRequired: true,
	},
	{
		Uri:          "/get-account-by-merchantcode",
		Method:       http.MethodGet,
		Handler:      controller.FindByMerchantCode,
		AuthRequired: true,
	},
}
