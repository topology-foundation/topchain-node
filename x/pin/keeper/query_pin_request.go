package keeper

import (
	"context"

	"topchain/x/pin/types"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) PinRequestAll(ctx context.Context, req *types.QueryAllPinRequestRequest) (*types.QueryAllPinRequestResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var pinRequests []types.PinRequest

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	pinRequestStore := prefix.NewStore(store, types.KeyPrefix(types.PinRequestKeyPrefix))

	pageRes, err := query.Paginate(pinRequestStore, req.Pagination, func(key []byte, value []byte) error {
		var pinRequest types.PinRequest
		if err := k.cdc.Unmarshal(value, &pinRequest); err != nil {
			return err
		}

		pinRequests = append(pinRequests, pinRequest)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllPinRequestResponse{PinRequest: pinRequests, Pagination: pageRes}, nil
}

func (k Keeper) PinRequest(ctx context.Context, req *types.QueryGetPinRequestRequest) (*types.QueryGetPinRequestResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, found := k.GetPinRequest(
		ctx,
		req.Index,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetPinRequestResponse{PinRequest: val}, nil
}
