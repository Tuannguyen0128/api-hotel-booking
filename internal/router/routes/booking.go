package routes

import (
	"api-hotel-booking/internal/controller"
	"net/http"
)

var bookingRouter = []Route{
	{
		Uri:          "/bookings",
		Method:       http.MethodGet,
		Handler:      controller.GetBookings,
		AuthRequired: true,
	},
	{
		Uri:          "/booking",
		Method:       http.MethodPost,
		Handler:      controller.CreateBooking,
		AuthRequired: true,
	},
	{
		Uri:          "/booking/:id",
		Method:       http.MethodPut,
		Handler:      controller.UpdateBooking,
		AuthRequired: true,
	},
	{
		Uri:          "/booking/:id",
		Method:       http.MethodDelete,
		Handler:      controller.DeleteBooking,
		AuthRequired: true,
	},
}
