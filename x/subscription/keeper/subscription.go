package keeper

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"topchain/x/subscription/types"
)

func (k Keeper) AddSubscription(ctx sdk.Context, subscription types.Subscription) string {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.SubscriptionKey))
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

func (k Keeper) RemoveSubscription(ctx sdk.Context, subscriptionId string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.SubscriptionKey))
	store.Delete([]byte(subscriptionId))
}

func (k Keeper) GetSubscription(ctx sdk.Context, subscriptionId string) (val types.Subscription, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.SubscriptionKey))
	b := store.Get([]byte(subscriptionId))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
