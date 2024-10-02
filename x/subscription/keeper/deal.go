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

	providerStore := prefix.NewStore(storeAdapter, types.GetRequesterStoreKey(deal.Requester))
	providerStore.Set([]byte(deal.Id), []byte{})
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

// check if at least one subscription is active
func (k Keeper) IsDealActive(ctx sdk.Context, deal types.Deal) bool {
	for _, subscriptionId := range deal.SubscriptionIds {
		subscription, found := k.GetSubscription(ctx, subscriptionId)
		if !found {
			continue
		}
		if subscription.StartBlock <= uint64(ctx.BlockHeight()) && subscription.EndBlock >= uint64(ctx.BlockHeight()) {
			return true
		}
	}
	return false
}

func (k Keeper) GetAllActiveProviders(ctx sdk.Context, deal types.Deal) []string {
	providers := []string{}
	for _, subscriptionId := range deal.SubscriptionIds {
		subscription, found := k.GetSubscription(ctx, subscriptionId)
		if !found {
			continue
		}
		if subscription.StartBlock <= uint64(ctx.BlockHeight()) && subscription.EndBlock >= uint64(ctx.BlockHeight()) {
			providers = append(providers, subscription.Provider)
		}
	}
	return providers
}

// Need a formula
func (k Keeper) CalculateMinimumStake(ctx sdk.Context, deal types.Deal) int64 {
	return 0
}

func (k Keeper) CalculateBlockReward(ctx sdk.Context, deal types.Deal) int64 {
	remainingBlocks := deal.EndBlock - uint64(ctx.BlockHeight())
	return int64(deal.AvailableAmount) / int64(remainingBlocks)
}

// Iterate over all deals and apply the given callback function
func (k Keeper) IterateDeals(ctx sdk.Context, shouldBreak func(deal types.Deal) bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := prefix.NewStore(storeAdapter, types.KeyPrefix(types.DealKeyPrefix)).Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var deal types.Deal
		k.cdc.MustUnmarshal(iterator.Value(), &deal)
		if shouldBreak(deal) {
			break
		}
	}
}
