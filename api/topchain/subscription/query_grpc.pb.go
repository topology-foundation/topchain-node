// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: topchain/subscription/query.proto

package subscription

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
	Query_Params_FullMethodName        = "/topchain.subscription.Query/Params"
	Query_Deal_FullMethodName          = "/topchain.subscription.Query/Deal"
	Query_DealStatus_FullMethodName    = "/topchain.subscription.Query/DealStatus"
	Query_Deals_FullMethodName         = "/topchain.subscription.Query/Deals"
	Query_Subscription_FullMethodName  = "/topchain.subscription.Query/Subscription"
	Query_Subscriptions_FullMethodName = "/topchain.subscription.Query/Subscriptions"
)

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueryClient interface {
	// Parameters queries the parameters of the module.
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
	Deal(ctx context.Context, in *QueryDealRequest, opts ...grpc.CallOption) (*QueryDealResponse, error)
	DealStatus(ctx context.Context, in *QueryDealStatusRequest, opts ...grpc.CallOption) (*QueryDealStatusResponse, error)
	Deals(ctx context.Context, in *QueryDealsRequest, opts ...grpc.CallOption) (*QueryDealsResponse, error)
	Subscription(ctx context.Context, in *QuerySubscriptionRequest, opts ...grpc.CallOption) (*QuerySubscriptionResponse, error)
	Subscriptions(ctx context.Context, in *QuerySubscriptionsRequest, opts ...grpc.CallOption) (*QuerySubscriptionsResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, Query_Params_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Deal(ctx context.Context, in *QueryDealRequest, opts ...grpc.CallOption) (*QueryDealResponse, error) {
	out := new(QueryDealResponse)
	err := c.cc.Invoke(ctx, Query_Deal_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) DealStatus(ctx context.Context, in *QueryDealStatusRequest, opts ...grpc.CallOption) (*QueryDealStatusResponse, error) {
	out := new(QueryDealStatusResponse)
	err := c.cc.Invoke(ctx, Query_DealStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Deals(ctx context.Context, in *QueryDealsRequest, opts ...grpc.CallOption) (*QueryDealsResponse, error) {
	out := new(QueryDealsResponse)
	err := c.cc.Invoke(ctx, Query_Deals_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Subscription(ctx context.Context, in *QuerySubscriptionRequest, opts ...grpc.CallOption) (*QuerySubscriptionResponse, error) {
	out := new(QuerySubscriptionResponse)
	err := c.cc.Invoke(ctx, Query_Subscription_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Subscriptions(ctx context.Context, in *QuerySubscriptionsRequest, opts ...grpc.CallOption) (*QuerySubscriptionsResponse, error) {
	out := new(QuerySubscriptionsResponse)
	err := c.cc.Invoke(ctx, Query_Subscriptions_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
// All implementations must embed UnimplementedQueryServer
// for forward compatibility
type QueryServer interface {
	// Parameters queries the parameters of the module.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	Deal(context.Context, *QueryDealRequest) (*QueryDealResponse, error)
	DealStatus(context.Context, *QueryDealStatusRequest) (*QueryDealStatusResponse, error)
	Deals(context.Context, *QueryDealsRequest) (*QueryDealsResponse, error)
	Subscription(context.Context, *QuerySubscriptionRequest) (*QuerySubscriptionResponse, error)
	Subscriptions(context.Context, *QuerySubscriptionsRequest) (*QuerySubscriptionsResponse, error)
	mustEmbedUnimplementedQueryServer()
}

// UnimplementedQueryServer must be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (UnimplementedQueryServer) Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (UnimplementedQueryServer) Deal(context.Context, *QueryDealRequest) (*QueryDealResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Deal not implemented")
}
func (UnimplementedQueryServer) DealStatus(context.Context, *QueryDealStatusRequest) (*QueryDealStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DealStatus not implemented")
}
func (UnimplementedQueryServer) Deals(context.Context, *QueryDealsRequest) (*QueryDealsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Deals not implemented")
}
func (UnimplementedQueryServer) Subscription(context.Context, *QuerySubscriptionRequest) (*QuerySubscriptionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Subscription not implemented")
}
func (UnimplementedQueryServer) Subscriptions(context.Context, *QuerySubscriptionsRequest) (*QuerySubscriptionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Subscriptions not implemented")
}
func (UnimplementedQueryServer) mustEmbedUnimplementedQueryServer() {}

// UnsafeQueryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueryServer will
// result in compilation errors.
type UnsafeQueryServer interface {
	mustEmbedUnimplementedQueryServer()
}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&Query_ServiceDesc, srv)
}

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Params_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Deal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryDealRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Deal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Deal_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Deal(ctx, req.(*QueryDealRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_DealStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryDealStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).DealStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_DealStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).DealStatus(ctx, req.(*QueryDealStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Deals_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryDealsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Deals(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Deals_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Deals(ctx, req.(*QueryDealsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Subscription_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuerySubscriptionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Subscription(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Subscription_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Subscription(ctx, req.(*QuerySubscriptionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Subscriptions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuerySubscriptionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Subscriptions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Subscriptions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Subscriptions(ctx, req.(*QuerySubscriptionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Query_ServiceDesc is the grpc.ServiceDesc for Query service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "topchain.subscription.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "Deal",
			Handler:    _Query_Deal_Handler,
		},
		{
			MethodName: "DealStatus",
			Handler:    _Query_DealStatus_Handler,
		},
		{
			MethodName: "Deals",
			Handler:    _Query_Deals_Handler,
		},
		{
			MethodName: "Subscription",
			Handler:    _Query_Subscription_Handler,
		},
		{
			MethodName: "Subscriptions",
			Handler:    _Query_Subscriptions_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "topchain/subscription/query.proto",
}
