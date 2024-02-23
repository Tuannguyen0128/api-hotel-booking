package client

import (
	"api-hotel-booking/internal/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var GrpcClient *Client

func InitGRPC() {
	cc, err := grpc.Dial("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("err when dial", err)
	}
	accountClient := proto.NewAccountServiceClient(cc)
	GrpcClient = &Client{AccountClient: accountClient}
}

type Client struct {
	AccountClient proto.AccountServiceClient
	StaffClient   proto.StaffServiceClient
}
