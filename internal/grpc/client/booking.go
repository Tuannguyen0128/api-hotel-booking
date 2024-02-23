package client

import (
	"api-hotel-booking/internal/grpc/proto"
	"context"
	"log"
)

func GetBookings(req *proto.GetBookingsRequest, c proto.BookingServiceClient) []*proto.Booking {
	res, err := c.GetBookings(context.Background(), req)
	if err != nil {
		log.Fatalln("err when getting all bookings", err)
	}
	return res.GetBookings()
}

func CreateBooking(req *proto.Booking, c proto.BookingServiceClient) *proto.Booking {
	res, err := c.CreateBooking(context.Background(), req)
	if err != nil {
		log.Fatalln("err when creating booking", err)
	}
	return res
}

func UpdateBooking(req *proto.Booking, c proto.BookingServiceClient) *proto.Booking {
	res, err := c.UpdateBooking(context.Background(), req)
	if err != nil {
		log.Fatalln("err when updating booking", err)
	}
	return res
}

func DeleteBooking(req *proto.DeleteBookingRequest, c proto.BookingServiceClient) *proto.DeleteBookingResponse {
	res, err := c.DeleteBooking(context.Background(), req)
	if err != nil {
		log.Fatalln("err when deleting booking", err)
	}
	return res
}
