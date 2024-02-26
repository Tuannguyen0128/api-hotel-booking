package client

import (
	"api-hotel-booking/internal/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

var GrpcClient *Client

func InitGRPC() {
	cc, err := grpc.Dial("localhost:3001", grpc.WithIdleTimeout(5*time.Second), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("err when dial", err)
	}
	accountClient := proto.NewAccountServiceClient(cc)
	staffClient := proto.NewStaffServiceClient(cc)
	bookingClient := proto.NewBookingServiceClient(cc)

	GrpcClient = &Client{
		AccountClient: accountClient,
		StaffClient:   staffClient,
		BookingClient: bookingClient,
	}
}

type Client struct {
	AccountClient proto.AccountServiceClient
	StaffClient   proto.StaffServiceClient
	BookingClient proto.BookingServiceClient
}
