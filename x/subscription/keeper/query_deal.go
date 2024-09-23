package keeper

import (
	"context"
	"topchain/x/subscription/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Deal(ctx context.Context, req *types.QueryDealRequest) (*types.QueryDealResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &types.QueryDealResponse{}, nil
}

func (k Keeper) DealStatus(ctx context.Context, req *types.QueryDealStatusRequest) (*types.QueryDealStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &types.QueryDealStatusResponse{}, nil
}

func (k Keeper) Deals(ctx context.Context, req *types.QueryDealsRequest) (*types.QueryDealsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &types.QueryDealsResponse{}, nil
}
