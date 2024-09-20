package keeper

import (
	"context"
	"topchain/x/subscription/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Subscription(ctx context.Context, req *types.QuerySubscriptionRequest) (*types.QuerySubscriptionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &types.QuerySubscriptionResponse{}, nil
}

func (k Keeper) Subscriptions(ctx context.Context, req *types.QuerySubscriptionsRequest) (*types.QuerySubscriptionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &types.QuerySubscriptionsResponse{}, nil
}
