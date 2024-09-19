package keeper

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"topchain/x/requester/types"
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
