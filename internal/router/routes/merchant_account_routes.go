package routes

import (
	"api-hotel-booking/internal/controller"
	"net/http"
)

var merchantroutes = []Route{
	{
		Uri:          "/merchantaccounts",
		Method:       http.MethodGet,
		Handler:      controller.GetMerchantAccounts,
		AuthRequired: true,
	},
	{
		Uri:          "/merchantaccount",
		Method:       http.MethodPost,
		Handler:      controller.CreateMerchantAccount,
		AuthRequired: true,
	},
	{
		Uri:          "/merchantaccount/{id}",
		Method:       http.MethodGet,
		Handler:      controller.FindMerchantAccountByID,
		AuthRequired: true,
	},
	{
		Uri:          "/merchantaccount/{id}",
		Method:       http.MethodPut,
		Handler:      controller.UpdateMerchantAccount,
		AuthRequired: true,
	},
	{
		Uri:          "/merchantaccount/{id}",
		Method:       http.MethodDelete,
		Handler:      controller.DeleteMerchantAccount,
		AuthRequired: true,
	},
}
