package keeper

import (
	"context"
	"topchain/x/subscription/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) DealStatus(ctx context.Context, req *types.QueryDealStatusRequest) (*types.QueryDealStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &types.QueryDealStatusResponse{}, nil
}
