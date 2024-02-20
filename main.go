package main

import (
	"api-hotel-booking/internal/grpc/client"
	"api-hotel-booking/internal/server"
)

func main() {
	client.InitGRPC()
	server.Run()
}
