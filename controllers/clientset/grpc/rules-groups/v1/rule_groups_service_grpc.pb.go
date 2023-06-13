// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.3
// source: com/coralogix/rules/v1/rule_groups_service.proto

package __

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

// RuleGroupsServiceClient is the client API for RuleGroupsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RuleGroupsServiceClient interface {
	GetRuleGroup(ctx context.Context, in *GetRuleGroupRequest, opts ...grpc.CallOption) (*GetRuleGroupResponse, error)
	CreateRuleGroup(ctx context.Context, in *CreateRuleGroupRequest, opts ...grpc.CallOption) (*CreateRuleGroupResponse, error)
	UpdateRuleGroup(ctx context.Context, in *UpdateRuleGroupRequest, opts ...grpc.CallOption) (*UpdateRuleGroupResponse, error)
	DeleteRuleGroup(ctx context.Context, in *DeleteRuleGroupRequest, opts ...grpc.CallOption) (*DeleteRuleGroupResponse, error)
}

type ruleGroupsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRuleGroupsServiceClient(cc grpc.ClientConnInterface) RuleGroupsServiceClient {
	return &ruleGroupsServiceClient{cc}
}

func (c *ruleGroupsServiceClient) GetRuleGroup(ctx context.Context, in *GetRuleGroupRequest, opts ...grpc.CallOption) (*GetRuleGroupResponse, error) {
	out := new(GetRuleGroupResponse)
	err := c.cc.Invoke(ctx, "/com.coralogix.rules.v1.RuleGroupsService/GetRuleGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ruleGroupsServiceClient) CreateRuleGroup(ctx context.Context, in *CreateRuleGroupRequest, opts ...grpc.CallOption) (*CreateRuleGroupResponse, error) {
	out := new(CreateRuleGroupResponse)
	err := c.cc.Invoke(ctx, "/com.coralogix.rules.v1.RuleGroupsService/CreateRuleGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ruleGroupsServiceClient) UpdateRuleGroup(ctx context.Context, in *UpdateRuleGroupRequest, opts ...grpc.CallOption) (*UpdateRuleGroupResponse, error) {
	out := new(UpdateRuleGroupResponse)
	err := c.cc.Invoke(ctx, "/com.coralogix.rules.v1.RuleGroupsService/UpdateRuleGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ruleGroupsServiceClient) DeleteRuleGroup(ctx context.Context, in *DeleteRuleGroupRequest, opts ...grpc.CallOption) (*DeleteRuleGroupResponse, error) {
	out := new(DeleteRuleGroupResponse)
	err := c.cc.Invoke(ctx, "/com.coralogix.rules.v1.RuleGroupsService/DeleteRuleGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RuleGroupsServiceServer is the server API for RuleGroupsService service.
// All implementations must embed UnimplementedRuleGroupsServiceServer
// for forward compatibility
type RuleGroupsServiceServer interface {
	GetRuleGroup(context.Context, *GetRuleGroupRequest) (*GetRuleGroupResponse, error)
	CreateRuleGroup(context.Context, *CreateRuleGroupRequest) (*CreateRuleGroupResponse, error)
	UpdateRuleGroup(context.Context, *UpdateRuleGroupRequest) (*UpdateRuleGroupResponse, error)
	DeleteRuleGroup(context.Context, *DeleteRuleGroupRequest) (*DeleteRuleGroupResponse, error)
	mustEmbedUnimplementedRuleGroupsServiceServer()
}

// UnimplementedRuleGroupsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRuleGroupsServiceServer struct {
}

func (UnimplementedRuleGroupsServiceServer) GetRuleGroup(context.Context, *GetRuleGroupRequest) (*GetRuleGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRuleGroup not implemented")
}
func (UnimplementedRuleGroupsServiceServer) CreateRuleGroup(context.Context, *CreateRuleGroupRequest) (*CreateRuleGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRuleGroup not implemented")
}
func (UnimplementedRuleGroupsServiceServer) UpdateRuleGroup(context.Context, *UpdateRuleGroupRequest) (*UpdateRuleGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRuleGroup not implemented")
}
func (UnimplementedRuleGroupsServiceServer) DeleteRuleGroup(context.Context, *DeleteRuleGroupRequest) (*DeleteRuleGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRuleGroup not implemented")
}
func (UnimplementedRuleGroupsServiceServer) mustEmbedUnimplementedRuleGroupsServiceServer() {}

// UnsafeRuleGroupsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RuleGroupsServiceServer will
// result in compilation errors.
type UnsafeRuleGroupsServiceServer interface {
	mustEmbedUnimplementedRuleGroupsServiceServer()
}

func RegisterRuleGroupsServiceServer(s grpc.ServiceRegistrar, srv RuleGroupsServiceServer) {
	s.RegisterService(&RuleGroupsService_ServiceDesc, srv)
}

func _RuleGroupsService_GetRuleGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRuleGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuleGroupsServiceServer).GetRuleGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.coralogix.rules.v1.RuleGroupsService/GetRuleGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuleGroupsServiceServer).GetRuleGroup(ctx, req.(*GetRuleGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuleGroupsService_CreateRuleGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRuleGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuleGroupsServiceServer).CreateRuleGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.coralogix.rules.v1.RuleGroupsService/CreateRuleGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuleGroupsServiceServer).CreateRuleGroup(ctx, req.(*CreateRuleGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuleGroupsService_UpdateRuleGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRuleGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuleGroupsServiceServer).UpdateRuleGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.coralogix.rules.v1.RuleGroupsService/UpdateRuleGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuleGroupsServiceServer).UpdateRuleGroup(ctx, req.(*UpdateRuleGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuleGroupsService_DeleteRuleGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRuleGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuleGroupsServiceServer).DeleteRuleGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.coralogix.rules.v1.RuleGroupsService/DeleteRuleGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuleGroupsServiceServer).DeleteRuleGroup(ctx, req.(*DeleteRuleGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RuleGroupsService_ServiceDesc is the grpc.ServiceDesc for RuleGroupsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RuleGroupsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "com.coralogix.rules.v1.RuleGroupsService",
	HandlerType: (*RuleGroupsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRuleGroup",
			Handler:    _RuleGroupsService_GetRuleGroup_Handler,
		},
		{
			MethodName: "CreateRuleGroup",
			Handler:    _RuleGroupsService_CreateRuleGroup_Handler,
		},
		{
			MethodName: "UpdateRuleGroup",
			Handler:    _RuleGroupsService_UpdateRuleGroup_Handler,
		},
		{
			MethodName: "DeleteRuleGroup",
			Handler:    _RuleGroupsService_DeleteRuleGroup_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "com/coralogix/rules/v1/rule_groups_service.proto",
}
