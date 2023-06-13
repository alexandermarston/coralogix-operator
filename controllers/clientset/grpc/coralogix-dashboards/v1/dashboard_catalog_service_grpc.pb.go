// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.8
// source: com/coralogixapis/dashboards/v1/services/dashboard_catalog_service.proto

package v1

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

// DashboardCatalogServiceClient is the client API for DashboardCatalogService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DashboardCatalogServiceClient interface {
	GetDashboardCatalog(ctx context.Context, in *GetDashboardCatalogRequest, opts ...grpc.CallOption) (*GetDashboardCatalogResponse, error)
	ReplaceDefaultDashboard(ctx context.Context, in *ReplaceDefaultDashboardRequest, opts ...grpc.CallOption) (*ReplaceDefaultDashboardResponse, error)
	PinDashboard(ctx context.Context, in *PinDashboardRequest, opts ...grpc.CallOption) (*PinDashboardResponse, error)
	UnpinDashboard(ctx context.Context, in *UnpinDashboardRequest, opts ...grpc.CallOption) (*UnpinDashboardResponse, error)
}

type dashboardCatalogServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDashboardCatalogServiceClient(cc grpc.ClientConnInterface) DashboardCatalogServiceClient {
	return &dashboardCatalogServiceClient{cc}
}

func (c *dashboardCatalogServiceClient) GetDashboardCatalog(ctx context.Context, in *GetDashboardCatalogRequest, opts ...grpc.CallOption) (*GetDashboardCatalogResponse, error) {
	out := new(GetDashboardCatalogResponse)
	err := c.cc.Invoke(ctx, "/com.coralogixapis.dashboards.v1.services.DashboardCatalogService/GetDashboardCatalog", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dashboardCatalogServiceClient) ReplaceDefaultDashboard(ctx context.Context, in *ReplaceDefaultDashboardRequest, opts ...grpc.CallOption) (*ReplaceDefaultDashboardResponse, error) {
	out := new(ReplaceDefaultDashboardResponse)
	err := c.cc.Invoke(ctx, "/com.coralogixapis.dashboards.v1.services.DashboardCatalogService/ReplaceDefaultDashboard", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dashboardCatalogServiceClient) PinDashboard(ctx context.Context, in *PinDashboardRequest, opts ...grpc.CallOption) (*PinDashboardResponse, error) {
	out := new(PinDashboardResponse)
	err := c.cc.Invoke(ctx, "/com.coralogixapis.dashboards.v1.services.DashboardCatalogService/PinDashboard", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dashboardCatalogServiceClient) UnpinDashboard(ctx context.Context, in *UnpinDashboardRequest, opts ...grpc.CallOption) (*UnpinDashboardResponse, error) {
	out := new(UnpinDashboardResponse)
	err := c.cc.Invoke(ctx, "/com.coralogixapis.dashboards.v1.services.DashboardCatalogService/UnpinDashboard", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DashboardCatalogServiceServer is the server API for DashboardCatalogService service.
// All implementations must embed UnimplementedDashboardCatalogServiceServer
// for forward compatibility
type DashboardCatalogServiceServer interface {
	GetDashboardCatalog(context.Context, *GetDashboardCatalogRequest) (*GetDashboardCatalogResponse, error)
	ReplaceDefaultDashboard(context.Context, *ReplaceDefaultDashboardRequest) (*ReplaceDefaultDashboardResponse, error)
	PinDashboard(context.Context, *PinDashboardRequest) (*PinDashboardResponse, error)
	UnpinDashboard(context.Context, *UnpinDashboardRequest) (*UnpinDashboardResponse, error)
	mustEmbedUnimplementedDashboardCatalogServiceServer()
}

// UnimplementedDashboardCatalogServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDashboardCatalogServiceServer struct {
}

func (UnimplementedDashboardCatalogServiceServer) GetDashboardCatalog(context.Context, *GetDashboardCatalogRequest) (*GetDashboardCatalogResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDashboardCatalog not implemented")
}
func (UnimplementedDashboardCatalogServiceServer) ReplaceDefaultDashboard(context.Context, *ReplaceDefaultDashboardRequest) (*ReplaceDefaultDashboardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReplaceDefaultDashboard not implemented")
}
func (UnimplementedDashboardCatalogServiceServer) PinDashboard(context.Context, *PinDashboardRequest) (*PinDashboardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PinDashboard not implemented")
}
func (UnimplementedDashboardCatalogServiceServer) UnpinDashboard(context.Context, *UnpinDashboardRequest) (*UnpinDashboardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnpinDashboard not implemented")
}
func (UnimplementedDashboardCatalogServiceServer) mustEmbedUnimplementedDashboardCatalogServiceServer() {
}

// UnsafeDashboardCatalogServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DashboardCatalogServiceServer will
// result in compilation errors.
type UnsafeDashboardCatalogServiceServer interface {
	mustEmbedUnimplementedDashboardCatalogServiceServer()
}

func RegisterDashboardCatalogServiceServer(s grpc.ServiceRegistrar, srv DashboardCatalogServiceServer) {
	s.RegisterService(&DashboardCatalogService_ServiceDesc, srv)
}

func _DashboardCatalogService_GetDashboardCatalog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDashboardCatalogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DashboardCatalogServiceServer).GetDashboardCatalog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.coralogixapis.dashboards.v1.services.DashboardCatalogService/GetDashboardCatalog",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DashboardCatalogServiceServer).GetDashboardCatalog(ctx, req.(*GetDashboardCatalogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DashboardCatalogService_ReplaceDefaultDashboard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReplaceDefaultDashboardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DashboardCatalogServiceServer).ReplaceDefaultDashboard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.coralogixapis.dashboards.v1.services.DashboardCatalogService/ReplaceDefaultDashboard",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DashboardCatalogServiceServer).ReplaceDefaultDashboard(ctx, req.(*ReplaceDefaultDashboardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DashboardCatalogService_PinDashboard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PinDashboardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DashboardCatalogServiceServer).PinDashboard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.coralogixapis.dashboards.v1.services.DashboardCatalogService/PinDashboard",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DashboardCatalogServiceServer).PinDashboard(ctx, req.(*PinDashboardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DashboardCatalogService_UnpinDashboard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnpinDashboardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DashboardCatalogServiceServer).UnpinDashboard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.coralogixapis.dashboards.v1.services.DashboardCatalogService/UnpinDashboard",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DashboardCatalogServiceServer).UnpinDashboard(ctx, req.(*UnpinDashboardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DashboardCatalogService_ServiceDesc is the grpc.ServiceDesc for DashboardCatalogService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DashboardCatalogService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "com.coralogixapis.dashboards.v1.services.DashboardCatalogService",
	HandlerType: (*DashboardCatalogServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDashboardCatalog",
			Handler:    _DashboardCatalogService_GetDashboardCatalog_Handler,
		},
		{
			MethodName: "ReplaceDefaultDashboard",
			Handler:    _DashboardCatalogService_ReplaceDefaultDashboard_Handler,
		},
		{
			MethodName: "PinDashboard",
			Handler:    _DashboardCatalogService_PinDashboard_Handler,
		},
		{
			MethodName: "UnpinDashboard",
			Handler:    _DashboardCatalogService_UnpinDashboard_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "com/coralogixapis/dashboards/v1/services/dashboard_catalog_service.proto",
}
