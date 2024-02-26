package routes

import (
	"api-hotel-booking/internal/controller"
	"net/http"
)

var staffRouter = []Route{
	{
		Uri:          "/staffs",
		Method:       http.MethodGet,
		Handler:      controller.GetStaffs,
		AuthRequired: true,
	},
	{
		Uri:          "/staff",
		Method:       http.MethodPost,
		Handler:      controller.CreateStaff,
		AuthRequired: true,
	},
	{
		Uri:          "/staff/:id",
		Method:       http.MethodPut,
		Handler:      controller.UpdateStaff,
		AuthRequired: true,
	},
	{
		Uri:          "/staff/:id",
		Method:       http.MethodDelete,
		Handler:      controller.DeleteStaff,
		AuthRequired: true,
	},
}
