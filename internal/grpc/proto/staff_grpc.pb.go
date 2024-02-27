// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.3
// source: staff.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	StaffService_GetStaffs_FullMethodName   = "/protobufs.StaffService/GetStaffs"
	StaffService_CreateStaff_FullMethodName = "/protobufs.StaffService/CreateStaff"
	StaffService_UpdateStaff_FullMethodName = "/protobufs.StaffService/UpdateStaff"
	StaffService_DeleteStaff_FullMethodName = "/protobufs.StaffService/DeleteStaff"
)

// StaffServiceClient is the client API for StaffService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StaffServiceClient interface {
	GetStaffs(ctx context.Context, in *GetStaffsRequest, opts ...grpc.CallOption) (*GetStaffsResponse, error)
	CreateStaff(ctx context.Context, in *Staff, opts ...grpc.CallOption) (*Staff, error)
	UpdateStaff(ctx context.Context, in *Staff, opts ...grpc.CallOption) (*Staff, error)
	DeleteStaff(ctx context.Context, in *DeleteStaffRequest, opts ...grpc.CallOption) (*DeleteStaffResponse, error)
}

type staffServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStaffServiceClient(cc grpc.ClientConnInterface) StaffServiceClient {
	return &staffServiceClient{cc}
}

func (c *staffServiceClient) GetStaffs(ctx context.Context, in *GetStaffsRequest, opts ...grpc.CallOption) (*GetStaffsResponse, error) {
	out := new(GetStaffsResponse)
	err := c.cc.Invoke(ctx, StaffService_GetStaffs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *staffServiceClient) CreateStaff(ctx context.Context, in *Staff, opts ...grpc.CallOption) (*Staff, error) {
	out := new(Staff)
	err := c.cc.Invoke(ctx, StaffService_CreateStaff_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *staffServiceClient) UpdateStaff(ctx context.Context, in *Staff, opts ...grpc.CallOption) (*Staff, error) {
	out := new(Staff)
	err := c.cc.Invoke(ctx, StaffService_UpdateStaff_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *staffServiceClient) DeleteStaff(ctx context.Context, in *DeleteStaffRequest, opts ...grpc.CallOption) (*DeleteStaffResponse, error) {
	out := new(DeleteStaffResponse)
	err := c.cc.Invoke(ctx, StaffService_DeleteStaff_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StaffServiceServer is the server API for StaffService service.
// All implementations must embed UnimplementedStaffServiceServer
// for forward compatibility
type StaffServiceServer interface {
	GetStaffs(context.Context, *GetStaffsRequest) (*GetStaffsResponse, error)
	CreateStaff(context.Context, *Staff) (*Staff, error)
	UpdateStaff(context.Context, *Staff) (*Staff, error)
	DeleteStaff(context.Context, *DeleteStaffRequest) (*DeleteStaffResponse, error)
	mustEmbedUnimplementedStaffServiceServer()
}

// UnimplementedStaffServiceServer must be embedded to have forward compatible implementations.
type UnimplementedStaffServiceServer struct {
}

func (UnimplementedStaffServiceServer) GetStaffs(context.Context, *GetStaffsRequest) (*GetStaffsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStaffs not implemented")
}
func (UnimplementedStaffServiceServer) CreateStaff(context.Context, *Staff) (*Staff, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateStaff not implemented")
}
func (UnimplementedStaffServiceServer) UpdateStaff(context.Context, *Staff) (*Staff, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateStaff not implemented")
}
func (UnimplementedStaffServiceServer) DeleteStaff(context.Context, *DeleteStaffRequest) (*DeleteStaffResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteStaff not implemented")
}
func (UnimplementedStaffServiceServer) mustEmbedUnimplementedStaffServiceServer() {}

// UnsafeStaffServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StaffServiceServer will
// result in compilation errors.
type UnsafeStaffServiceServer interface {
	mustEmbedUnimplementedStaffServiceServer()
}

func RegisterStaffServiceServer(s grpc.ServiceRegistrar, srv StaffServiceServer) {
	s.RegisterService(&StaffService_ServiceDesc, srv)
}

func _StaffService_GetStaffs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStaffsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StaffServiceServer).GetStaffs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StaffService_GetStaffs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StaffServiceServer).GetStaffs(ctx, req.(*GetStaffsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StaffService_CreateStaff_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Staff)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StaffServiceServer).CreateStaff(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StaffService_CreateStaff_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StaffServiceServer).CreateStaff(ctx, req.(*Staff))
	}
	return interceptor(ctx, in, info, handler)
}

func _StaffService_UpdateStaff_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Staff)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StaffServiceServer).UpdateStaff(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StaffService_UpdateStaff_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StaffServiceServer).UpdateStaff(ctx, req.(*Staff))
	}
	return interceptor(ctx, in, info, handler)
}

func _StaffService_DeleteStaff_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteStaffRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StaffServiceServer).DeleteStaff(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StaffService_DeleteStaff_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StaffServiceServer).DeleteStaff(ctx, req.(*DeleteStaffRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// StaffService_ServiceDesc is the grpc.ServiceDesc for StaffService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StaffService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protobufs.StaffService",
	HandlerType: (*StaffServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetStaffs",
			Handler:    _StaffService_GetStaffs_Handler,
		},
		{
			MethodName: "CreateStaff",
			Handler:    _StaffService_CreateStaff_Handler,
		},
		{
			MethodName: "UpdateStaff",
			Handler:    _StaffService_UpdateStaff_Handler,
		},
		{
			MethodName: "DeleteStaff",
			Handler:    _StaffService_DeleteStaff_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "staff.proto",
}