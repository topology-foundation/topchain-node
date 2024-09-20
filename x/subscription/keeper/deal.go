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

func (k Keeper) AddDeal(ctx sdk.Context, deal types.Deal) string {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.DealKeyPrefix))

	appendedValue := k.cdc.MustMarshal(&deal)
	subscriptionStringified, err := json.Marshal(deal)
	if err != nil {
		panic(err)
	}

	hash := sha256.New()
	hash.Write(subscriptionStringified)
	hashed := hash.Sum(nil)
	store.Set(hashed, appendedValue)
	return hex.EncodeToString(hashed)
}

func (k Keeper) GetDeal(ctx sdk.Context, dealId string) (deal types.Deal, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.DealKeyPrefix))

	dealBytes := store.Get([]byte(dealId))
	if dealBytes == nil {
		return deal, false
	}

	k.cdc.MustUnmarshal(dealBytes, &deal)
	return deal, true
}

func (k Keeper) RemoveDeal(ctx sdk.Context, dealId string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.DealKeyPrefix))
	store.Delete([]byte(dealId))
}
