package keeper

import (
	"context"

	"topchain/x/subscription/types"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ListSubscription(ctx context.Context, req *types.QueryListSubscriptionRequest) (*types.QueryListSubscriptionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.SubscriptionKey))

	var subscriptions []types.Subscription
	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		var subscription types.Subscription
		if err := k.cdc.Unmarshal(value, &subscription); err != nil {
			return err
		}

		subscriptions = append(subscriptions, subscription)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryListSubscriptionResponse{Subscription: subscriptions, Pagination: pageRes}, nil
}
