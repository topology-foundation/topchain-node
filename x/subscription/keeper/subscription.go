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

	providerStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.SubscriptionProviderKeyPrefix))
	var subscriptionsIds types.SubscriptionIds
	existingSubscriptionIdsByte := providerStore.Get([]byte(subscription.Provider))
	if existingSubscriptionIdsByte != nil {
		k.cdc.MustUnmarshal(existingSubscriptionIdsByte, &subscriptionsIds)
	}

	for _, id := range subscriptionsIds.Ids {
		if id == subscription.Id {
			return
		}
	}
	subscriptionsIds.Ids = append(subscriptionsIds.Ids, subscription.Id)
	providerStore.Set([]byte(subscription.Provider), k.cdc.MustMarshal(&subscriptionsIds))
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
