package routes

import (
	"api-hotel-booking/internal/controller"
	"net/http"
)

var userRouter = []Route{
	{
		Uri:          "/teammembers",
		Method:       http.MethodGet,
		Handler:      controller.GetTeamMembers,
		AuthRequired: true,
	},
	{
		Uri:          "/teammember/{id}",
		Method:       http.MethodGet,
		Handler:      controller.GetTeamMembers,
		AuthRequired: true,
	},
	{
		Uri:          "/teammember",
		Method:       http.MethodPost,
		Handler:      controller.CreateTeamMember,
		AuthRequired: false,
	},
	{
		Uri:          "/teammember/{id}",
		Method:       http.MethodPut,
		Handler:      controller.UpdateTeamMember,
		AuthRequired: true,
	},
	{
		Uri:          "/teammember/{id}",
		Method:       http.MethodDelete,
		Handler:      controller.DeleteTeamMember,
		AuthRequired: true,
	},
	{
		Uri:          "/getteammemberbyemail",
		Method:       http.MethodGet,
		Handler:      controller.FindByEmail,
		AuthRequired: true,
	},
	{
		Uri:          "/gettemmemberbymerchantcode",
		Method:       http.MethodGet,
		Handler:      controller.FindByMerchantCode,
		AuthRequired: true,
	},
}
