package controller

import (
	"api-hotel-booking/internal/grpc/client"
	"api-hotel-booking/internal/grpc/proto"
	"api-hotel-booking/internal/models"
	"api-hotel-booking/internal/responses"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetBookings(c *gin.Context) {
	q := c.Request.URL.Query()
	limitS := q.Get("limit")
	limit, err := strconv.Atoi(limitS)
	if err != nil {
		log.Println(err)
	}

	pageS := q.Get("page")
	page, err := strconv.Atoi(pageS)
	if err != nil {
		log.Println(err)
	}

	id := q.Get("id")
	position := q.Get("position")
	staffID := q.Get("staff_id")

	bookings := client.GetBookings(&proto.GetBookingsRequest{Page: int32(page), Offset: int32(limit), Id: id, Position: position, StaffId: staffID}, client.GrpcClient.BookingClient)
	if bookings == nil {
		c.JSON(http.StatusOK, &proto.Booking{})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

func CreateBooking(c *gin.Context) {
	booking := models.Booking{}
	err := c.ShouldBindJSON(&booking)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}

	bookingRequest := &proto.Booking{
		Id:            "",
		GuestId:       booking.GuestID,
		AccountId:     booking.AccountID,
		CheckinTime:   booking.CheckinTime.String(),
		CheckoutTime:  booking.CheckoutTime.String(),
		TotalPrice:    fmt.Sprintf("%f", booking.TotalPrice),
		Status:        booking.Status,
		PaymentStatus: booking.PaymentStatus,
	}

	createdBooking := client.CreateBooking(bookingRequest, client.GrpcClient.BookingClient)
	c.JSON(http.StatusCreated, createdBooking)
}

func UpdateBooking(c *gin.Context) {
	booking := models.Booking{}
	err := c.ShouldBindJSON(&booking)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.ERROR(http.StatusUnprocessableEntity, err.Error()))
		return
	}

	bookingRequest := &proto.Booking{
		Id:            c.Param("id"),
		GuestId:       booking.GuestID,
		AccountId:     booking.AccountID,
		CheckinTime:   booking.CheckinTime.String(),
		CheckoutTime:  booking.CheckoutTime.String(),
		TotalPrice:    fmt.Sprintf("%f", booking.TotalPrice),
		Status:        booking.Status,
		PaymentStatus: booking.PaymentStatus,
	}
	fmt.Println("bookingRequest", bookingRequest)
	updatedBooking := client.UpdateBooking(bookingRequest, client.GrpcClient.BookingClient)
	c.JSON(http.StatusOK, updatedBooking)
}

func DeleteBooking(c *gin.Context) {
	id := c.Param("id")
	bookingId := &proto.DeleteBookingRequest{Id: id}

	deletedBooking := client.DeleteBooking(bookingId, client.GrpcClient.BookingClient)
	c.JSON(http.StatusOK, deletedBooking)
}
