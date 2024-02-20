package client

import (
	"api-hotel-booking/internal/grpc/proto"
	"context"
	"log"
)

func GetAllAccount(req proto.AccountRequest, c proto.AccountServiceClient) []*proto.Account {
	res, err := c.GetAllAccount(context.Background(), &req)
	if err != nil {
		log.Fatalln("err when call sum", err)
	}
	//log.Println("calSum response", res.Accounts())
	return res.GetAccounts()
}
