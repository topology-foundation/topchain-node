package keeper

import (
	"context"

	"topchain/x/subscription/types"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Subscription(goCtx context.Context, req *types.QuerySubscriptionRequest) (*types.QuerySubscriptionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	subscription, found := k.GetSubscription(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QuerySubscriptionResponse{Subscription: subscription}, nil
}

func (k Keeper) Subscriptions(goCtx context.Context, req *types.QuerySubscriptionsRequest) (*types.QuerySubscriptionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	req.Pagination = types.InitPagintionRequestDefaults(req.Pagination)

	ctx := sdk.UnwrapSDKContext(goCtx)
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	providerStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.SubscriptionProviderKeyPrefix))

	var subscriptionIds types.SubscriptionIds
	providerSubscriptions := providerStore.Get([]byte(req.Provider))
	if providerSubscriptions == nil {
		return &types.QuerySubscriptionsResponse{Subscriptions: []types.Subscription{}}, nil
	}
	k.cdc.MustUnmarshal(providerSubscriptions, &subscriptionIds)

	idsLen := uint64(len(subscriptionIds.Ids))
	if req.Pagination.Offset >= idsLen || idsLen == 0 {
		return &types.QuerySubscriptionsResponse{Subscriptions: []types.Subscription{}}, nil
	}
	end := req.Pagination.Offset + req.Pagination.Limit
	if end > idsLen {
		end = idsLen
	}

	var subscriptions []types.Subscription
	for _, id := range subscriptionIds.Ids[req.Pagination.Offset:end] {
		subscription, found := k.GetSubscription(ctx, id)
		if !found {
			return nil, sdkerrors.ErrKeyNotFound
		}
		subscriptions = append(subscriptions, subscription)
	}

	return &types.QuerySubscriptionsResponse{Subscriptions: subscriptions}, nil
}
