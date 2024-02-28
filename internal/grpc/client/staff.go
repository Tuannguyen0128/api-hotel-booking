package client

import (
	"api-hotel-booking/internal/grpc/proto"
	"context"
	"log"
)

func GetStaffs(req *proto.GetStaffsRequest, c proto.StaffServiceClient) []*proto.Staff {
	res, err := c.GetStaffs(context.Background(), req)
	if err != nil {
		log.Fatalln("err when getting all staffs", err)
	}
	return res.GetStaffs()
}

func CreateStaff(req *proto.Staff, c proto.StaffServiceClient) *proto.CreateStaffResponse {
	res, err := c.CreateStaff(context.Background(), req)
	if err != nil {
		log.Fatalln("err when creating staff", err)
	}
	return res
}

func UpdateStaff(req *proto.Staff, c proto.StaffServiceClient) *proto.Staff {
	res, err := c.UpdateStaff(context.Background(), req)
	if err != nil {
		log.Fatalln("err when updating staff", err)
	}
	return res
}

func DeleteStaff(req *proto.DeleteStaffRequest, c proto.StaffServiceClient) *proto.DeleteStaffResponse {
	res, err := c.DeleteStaff(context.Background(), req)
	if err != nil {
		log.Fatalln("err when deleting staff", err)
	}
	return res
}
