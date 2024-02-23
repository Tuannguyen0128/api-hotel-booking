package client

import (
	"api-hotel-booking/internal/grpc/proto"
	"context"
	"log"
)

func GetAccounts(req *proto.GetAccountsRequest, c proto.AccountServiceClient) []*proto.Account {
	res, err := c.GetAccounts(context.Background(), req)
	if err != nil {
		log.Fatalln("err when getting all accounts", err)
	}
	return res.GetAccounts()
}

func CreateAccount(req *proto.Account, c proto.AccountServiceClient) *proto.Account {
	res, err := c.CreateAccount(context.Background(), req)
	if err != nil {
		log.Fatalln("err when creating account", err)
	}
	return res
}

func UpdateAccount(req *proto.Account, c proto.AccountServiceClient) *proto.Account {
	res, err := c.UpdateAccount(context.Background(), req)
	if err != nil {
		log.Fatalln("err when updating account", err)
	}
	return res
}

func DeleteAccount(req *proto.DeleteAccountRequest, c proto.AccountServiceClient) *proto.DeleteAccountResponse {
	res, err := c.DeleteAccount(context.Background(), req)
	if err != nil {
		log.Fatalln("err when updating account", err)
	}
	return res
}
