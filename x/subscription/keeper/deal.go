package keeper

import (
	"topchain/x/subscription/types"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetDeal(ctx sdk.Context, deal types.Deal) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.DealKeyPrefix))

	appendedValue := k.cdc.MustMarshal(&deal)
	store.Set([]byte(deal.Id), appendedValue)
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

// Need a formula
func (k Keeper) CalculateMinimumStake(ctx sdk.Context, deal types.Deal) int64 {
	return 0
}
