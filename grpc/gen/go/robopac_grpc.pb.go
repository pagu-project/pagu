// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: robopac.proto

package robopac

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

// RoboPacClient is the client API for RoboPac service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RoboPacClient interface {
	Execute(ctx context.Context, in *ExecuteRequest, opts ...grpc.CallOption) (*ExecuteResponse, error)
}

type roboPacClient struct {
	cc grpc.ClientConnInterface
}

func NewRoboPacClient(cc grpc.ClientConnInterface) RoboPacClient {
	return &roboPacClient{cc}
}

func (c *roboPacClient) Execute(ctx context.Context, in *ExecuteRequest, opts ...grpc.CallOption) (*ExecuteResponse, error) {
	out := new(ExecuteResponse)
	err := c.cc.Invoke(ctx, "/robopac.RoboPac/Execute", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RoboPacServer is the server API for RoboPac service.
// All implementations should embed UnimplementedRoboPacServer
// for forward compatibility.
type RoboPacServer interface {
	Execute(context.Context, *ExecuteRequest) (*ExecuteResponse, error)
}

// UnimplementedRoboPacServer should be embedded to have forward compatible implementations.
type UnimplementedRoboPacServer struct {
}

func (UnimplementedRoboPacServer) Execute(context.Context, *ExecuteRequest) (*ExecuteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Execute not implemented")
}

// UnsafeRoboPacServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RoboPacServer will
// result in compilation errors.
type UnsafeRoboPacServer interface {
	mustEmbedUnimplementedRoboPacServer()
}

func RegisterRoboPacServer(s grpc.ServiceRegistrar, srv RoboPacServer) {
	s.RegisterService(&RoboPac_ServiceDesc, srv)
}

func _RoboPac_Execute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecuteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoboPacServer).Execute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/robopac.RoboPac/Execute",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoboPacServer).Execute(ctx, req.(*ExecuteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RoboPac_ServiceDesc is the grpc.ServiceDesc for RoboPac service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy).
var RoboPac_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "robopac.RoboPac",
	HandlerType: (*RoboPacServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Execute",
			Handler:    _RoboPac_Execute_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "robopac.proto",
}
