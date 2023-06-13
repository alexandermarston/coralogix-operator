// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.8
// source: com/coralogixapis/logs2metrics/v2/logs2metrics_service.proto

package __

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// Logs2MetricServiceClient is the client API for Logs2MetricService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type Logs2MetricServiceClient interface {
	CreateL2M(ctx context.Context, in *CreateL2MRequest, opts ...grpc.CallOption) (*L2M, error)
	ListL2M(ctx context.Context, in *ListL2MRequest, opts ...grpc.CallOption) (*ListL2MResponse, error)
	ReplaceL2M(ctx context.Context, in *ReplaceL2MRequest, opts ...grpc.CallOption) (*L2M, error)
	GetL2M(ctx context.Context, in *GetL2MRequest, opts ...grpc.CallOption) (*L2M, error)
	DeleteL2M(ctx context.Context, in *DeleteL2MRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	AtomicBatchExecuteL2M(ctx context.Context, in *AtomicBatchExecuteL2MRequest, opts ...grpc.CallOption) (*AtomicBatchExecuteL2MResponse, error)
}

type logs2MetricServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLogs2MetricServiceClient(cc grpc.ClientConnInterface) Logs2MetricServiceClient {
	return &logs2MetricServiceClient{cc}
}

func (c *logs2MetricServiceClient) CreateL2M(ctx context.Context, in *CreateL2MRequest, opts ...grpc.CallOption) (*L2M, error) {
	out := new(L2M)
	err := c.cc.Invoke(ctx, "/com.coralogixapis.logs2metrics.v2.Logs2MetricService/CreateL2M", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logs2MetricServiceClient) ListL2M(ctx context.Context, in *ListL2MRequest, opts ...grpc.CallOption) (*ListL2MResponse, error) {
	out := new(ListL2MResponse)
	err := c.cc.Invoke(ctx, "/com.coralogixapis.logs2metrics.v2.Logs2MetricService/ListL2M", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logs2MetricServiceClient) ReplaceL2M(ctx context.Context, in *ReplaceL2MRequest, opts ...grpc.CallOption) (*L2M, error) {
	out := new(L2M)
	err := c.cc.Invoke(ctx, "/com.coralogixapis.logs2metrics.v2.Logs2MetricService/ReplaceL2M", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logs2MetricServiceClient) GetL2M(ctx context.Context, in *GetL2MRequest, opts ...grpc.CallOption) (*L2M, error) {
	out := new(L2M)
	err := c.cc.Invoke(ctx, "/com.coralogixapis.logs2metrics.v2.Logs2MetricService/GetL2M", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logs2MetricServiceClient) DeleteL2M(ctx context.Context, in *DeleteL2MRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/com.coralogixapis.logs2metrics.v2.Logs2MetricService/DeleteL2M", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logs2MetricServiceClient) AtomicBatchExecuteL2M(ctx context.Context, in *AtomicBatchExecuteL2MRequest, opts ...grpc.CallOption) (*AtomicBatchExecuteL2MResponse, error) {
	out := new(AtomicBatchExecuteL2MResponse)
	err := c.cc.Invoke(ctx, "/com.coralogixapis.logs2metrics.v2.Logs2MetricService/AtomicBatchExecuteL2M", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Logs2MetricServiceServer is the server API for Logs2MetricService service.
// All implementations must embed UnimplementedLogs2MetricServiceServer
// for forward compatibility
type Logs2MetricServiceServer interface {
	CreateL2M(context.Context, *CreateL2MRequest) (*L2M, error)
	ListL2M(context.Context, *ListL2MRequest) (*ListL2MResponse, error)
	ReplaceL2M(context.Context, *ReplaceL2MRequest) (*L2M, error)
	GetL2M(context.Context, *GetL2MRequest) (*L2M, error)
	DeleteL2M(context.Context, *DeleteL2MRequest) (*emptypb.Empty, error)
	AtomicBatchExecuteL2M(context.Context, *AtomicBatchExecuteL2MRequest) (*AtomicBatchExecuteL2MResponse, error)
	mustEmbedUnimplementedLogs2MetricServiceServer()
}

// UnimplementedLogs2MetricServiceServer must be embedded to have forward compatible implementations.
type UnimplementedLogs2MetricServiceServer struct {
}

func (UnimplementedLogs2MetricServiceServer) CreateL2M(context.Context, *CreateL2MRequest) (*L2M, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateL2M not implemented")
}
func (UnimplementedLogs2MetricServiceServer) ListL2M(context.Context, *ListL2MRequest) (*ListL2MResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListL2M not implemented")
}
func (UnimplementedLogs2MetricServiceServer) ReplaceL2M(context.Context, *ReplaceL2MRequest) (*L2M, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReplaceL2M not implemented")
}
func (UnimplementedLogs2MetricServiceServer) GetL2M(context.Context, *GetL2MRequest) (*L2M, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetL2M not implemented")
}
func (UnimplementedLogs2MetricServiceServer) DeleteL2M(context.Context, *DeleteL2MRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteL2M not implemented")
}
func (UnimplementedLogs2MetricServiceServer) AtomicBatchExecuteL2M(context.Context, *AtomicBatchExecuteL2MRequest) (*AtomicBatchExecuteL2MResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AtomicBatchExecuteL2M not implemented")
}
func (UnimplementedLogs2MetricServiceServer) mustEmbedUnimplementedLogs2MetricServiceServer() {}

// UnsafeLogs2MetricServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to Logs2MetricServiceServer will
// result in compilation errors.
type UnsafeLogs2MetricServiceServer interface {
	mustEmbedUnimplementedLogs2MetricServiceServer()
}

func RegisterLogs2MetricServiceServer(s grpc.ServiceRegistrar, srv Logs2MetricServiceServer) {
	s.RegisterService(&Logs2MetricService_ServiceDesc, srv)
}

func _Logs2MetricService_CreateL2M_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateL2MRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Logs2MetricServiceServer).CreateL2M(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.coralogixapis.logs2metrics.v2.Logs2MetricService/CreateL2M",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Logs2MetricServiceServer).CreateL2M(ctx, req.(*CreateL2MRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Logs2MetricService_ListL2M_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListL2MRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Logs2MetricServiceServer).ListL2M(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.coralogixapis.logs2metrics.v2.Logs2MetricService/ListL2M",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Logs2MetricServiceServer).ListL2M(ctx, req.(*ListL2MRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Logs2MetricService_ReplaceL2M_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReplaceL2MRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Logs2MetricServiceServer).ReplaceL2M(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.coralogixapis.logs2metrics.v2.Logs2MetricService/ReplaceL2M",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Logs2MetricServiceServer).ReplaceL2M(ctx, req.(*ReplaceL2MRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Logs2MetricService_GetL2M_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetL2MRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Logs2MetricServiceServer).GetL2M(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.coralogixapis.logs2metrics.v2.Logs2MetricService/GetL2M",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Logs2MetricServiceServer).GetL2M(ctx, req.(*GetL2MRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Logs2MetricService_DeleteL2M_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteL2MRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Logs2MetricServiceServer).DeleteL2M(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.coralogixapis.logs2metrics.v2.Logs2MetricService/DeleteL2M",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Logs2MetricServiceServer).DeleteL2M(ctx, req.(*DeleteL2MRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Logs2MetricService_AtomicBatchExecuteL2M_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AtomicBatchExecuteL2MRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Logs2MetricServiceServer).AtomicBatchExecuteL2M(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.coralogixapis.logs2metrics.v2.Logs2MetricService/AtomicBatchExecuteL2M",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Logs2MetricServiceServer).AtomicBatchExecuteL2M(ctx, req.(*AtomicBatchExecuteL2MRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Logs2MetricService_ServiceDesc is the grpc.ServiceDesc for Logs2MetricService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Logs2MetricService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "com.coralogixapis.logs2metrics.v2.Logs2MetricService",
	HandlerType: (*Logs2MetricServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateL2M",
			Handler:    _Logs2MetricService_CreateL2M_Handler,
		},
		{
			MethodName: "ListL2M",
			Handler:    _Logs2MetricService_ListL2M_Handler,
		},
		{
			MethodName: "ReplaceL2M",
			Handler:    _Logs2MetricService_ReplaceL2M_Handler,
		},
		{
			MethodName: "GetL2M",
			Handler:    _Logs2MetricService_GetL2M_Handler,
		},
		{
			MethodName: "DeleteL2M",
			Handler:    _Logs2MetricService_DeleteL2M_Handler,
		},
		{
			MethodName: "AtomicBatchExecuteL2M",
			Handler:    _Logs2MetricService_AtomicBatchExecuteL2M_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "com/coralogixapis/logs2metrics/v2/logs2metrics_service.proto",
}
