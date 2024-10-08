// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.3
// source: counter.proto

package counter

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
	CounterService_GetUnreadCountTotalV1_FullMethodName = "/counter.CounterService/GetUnreadCountTotalV1"
	CounterService_GetUnreadCountV1_FullMethodName      = "/counter.CounterService/GetUnreadCountV1"
	CounterService_MarkMessagesAsReadV1_FullMethodName  = "/counter.CounterService/MarkMessagesAsReadV1"
)

// CounterServiceClient is the client API for CounterService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CounterServiceClient interface {
	GetUnreadCountTotalV1(ctx context.Context, in *GetUnreadCountTotalV1Request, opts ...grpc.CallOption) (*GetUnreadCountTotalV1Response, error)
	GetUnreadCountV1(ctx context.Context, in *GetUnreadCountV1Request, opts ...grpc.CallOption) (*GetUnreadCountV1Response, error)
	MarkMessagesAsReadV1(ctx context.Context, in *MarkMessagesAsReadV1Request, opts ...grpc.CallOption) (*MarkMessagesAsReadV1Response, error)
}

type counterServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCounterServiceClient(cc grpc.ClientConnInterface) CounterServiceClient {
	return &counterServiceClient{cc}
}

func (c *counterServiceClient) GetUnreadCountTotalV1(ctx context.Context, in *GetUnreadCountTotalV1Request, opts ...grpc.CallOption) (*GetUnreadCountTotalV1Response, error) {
	out := new(GetUnreadCountTotalV1Response)
	err := c.cc.Invoke(ctx, CounterService_GetUnreadCountTotalV1_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *counterServiceClient) GetUnreadCountV1(ctx context.Context, in *GetUnreadCountV1Request, opts ...grpc.CallOption) (*GetUnreadCountV1Response, error) {
	out := new(GetUnreadCountV1Response)
	err := c.cc.Invoke(ctx, CounterService_GetUnreadCountV1_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *counterServiceClient) MarkMessagesAsReadV1(ctx context.Context, in *MarkMessagesAsReadV1Request, opts ...grpc.CallOption) (*MarkMessagesAsReadV1Response, error) {
	out := new(MarkMessagesAsReadV1Response)
	err := c.cc.Invoke(ctx, CounterService_MarkMessagesAsReadV1_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CounterServiceServer is the server API for CounterService service.
// All implementations must embed UnimplementedCounterServiceServer
// for forward compatibility
type CounterServiceServer interface {
	GetUnreadCountTotalV1(context.Context, *GetUnreadCountTotalV1Request) (*GetUnreadCountTotalV1Response, error)
	GetUnreadCountV1(context.Context, *GetUnreadCountV1Request) (*GetUnreadCountV1Response, error)
	MarkMessagesAsReadV1(context.Context, *MarkMessagesAsReadV1Request) (*MarkMessagesAsReadV1Response, error)
	mustEmbedUnimplementedCounterServiceServer()
}

// UnimplementedCounterServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCounterServiceServer struct {
}

func (UnimplementedCounterServiceServer) GetUnreadCountTotalV1(context.Context, *GetUnreadCountTotalV1Request) (*GetUnreadCountTotalV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUnreadCountTotalV1 not implemented")
}
func (UnimplementedCounterServiceServer) GetUnreadCountV1(context.Context, *GetUnreadCountV1Request) (*GetUnreadCountV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUnreadCountV1 not implemented")
}
func (UnimplementedCounterServiceServer) MarkMessagesAsReadV1(context.Context, *MarkMessagesAsReadV1Request) (*MarkMessagesAsReadV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MarkMessagesAsReadV1 not implemented")
}
func (UnimplementedCounterServiceServer) mustEmbedUnimplementedCounterServiceServer() {}

// UnsafeCounterServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CounterServiceServer will
// result in compilation errors.
type UnsafeCounterServiceServer interface {
	mustEmbedUnimplementedCounterServiceServer()
}

func RegisterCounterServiceServer(s grpc.ServiceRegistrar, srv CounterServiceServer) {
	s.RegisterService(&CounterService_ServiceDesc, srv)
}

func _CounterService_GetUnreadCountTotalV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUnreadCountTotalV1Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CounterServiceServer).GetUnreadCountTotalV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CounterService_GetUnreadCountTotalV1_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CounterServiceServer).GetUnreadCountTotalV1(ctx, req.(*GetUnreadCountTotalV1Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _CounterService_GetUnreadCountV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUnreadCountV1Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CounterServiceServer).GetUnreadCountV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CounterService_GetUnreadCountV1_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CounterServiceServer).GetUnreadCountV1(ctx, req.(*GetUnreadCountV1Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _CounterService_MarkMessagesAsReadV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MarkMessagesAsReadV1Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CounterServiceServer).MarkMessagesAsReadV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CounterService_MarkMessagesAsReadV1_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CounterServiceServer).MarkMessagesAsReadV1(ctx, req.(*MarkMessagesAsReadV1Request))
	}
	return interceptor(ctx, in, info, handler)
}

// CounterService_ServiceDesc is the grpc.ServiceDesc for CounterService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CounterService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "counter.CounterService",
	HandlerType: (*CounterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUnreadCountTotalV1",
			Handler:    _CounterService_GetUnreadCountTotalV1_Handler,
		},
		{
			MethodName: "GetUnreadCountV1",
			Handler:    _CounterService_GetUnreadCountV1_Handler,
		},
		{
			MethodName: "MarkMessagesAsReadV1",
			Handler:    _CounterService_MarkMessagesAsReadV1_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "counter.proto",
}
