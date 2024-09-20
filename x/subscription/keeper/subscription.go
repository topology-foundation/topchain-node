package keeper

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"topchain/x/subscription/types"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) AddSubscription(ctx sdk.Context, subscription types.Subscription) string {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.SubscriptionKeyPrefix))
	appendedValue := k.cdc.MustMarshal(&subscription)
	subscriptionStringified, err := json.Marshal(subscription)
	if err != nil {
		panic(err)
	}
	hash := sha256.New()
	hash.Write(subscriptionStringified)
	hashed := hash.Sum(nil)
	store.Set(hashed, appendedValue)
	return hex.EncodeToString(hashed)
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

func (k Keeper) RemoveSubscription(ctx sdk.Context, subscriptionId string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.SubscriptionKeyPrefix))
	store.Delete([]byte(subscriptionId))
}
