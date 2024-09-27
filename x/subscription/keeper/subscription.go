package keeper

import (
	"topchain/x/subscription/types"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetSubscription(ctx sdk.Context, subscription types.Subscription) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.SubscriptionKeyPrefix))

	appendedValue := k.cdc.MustMarshal(&subscription)
	store.Set([]byte(subscription.Id), appendedValue)

	providerStore := prefix.NewStore(storeAdapter, types.GetProviderStoreKey(subscription.Provider))
	providerStore.Set([]byte(subscription.Id), appendedValue)
}

func (k Keeper) GetSubscription(ctx sdk.Context, subscriptionId string) (subscription types.Subscription, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.SubscriptionKeyPrefix))

	subscriptionBytes := store.Get([]byte(subscriptionId))
	if subscriptionBytes == nil {
		return subscription, false
	}

	k.cdc.MustUnmarshal(subscriptionBytes, &subscription)
	return subscription, true
}
