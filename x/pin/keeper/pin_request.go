package keeper

import (
	"context"

	"topchain/x/pin/types"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
)

// SetPinRequest set a specific pinRequest in the store from its index
func (k Keeper) SetPinRequest(ctx context.Context, pinRequest types.PinRequest) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PinRequestKeyPrefix))
	b := k.cdc.MustMarshal(&pinRequest)
	store.Set(types.PinRequestKey(
		pinRequest.Index,
	), b)
}

// GetPinRequest returns a pinRequest from its index
func (k Keeper) GetPinRequest(
	ctx context.Context,
	index string,

) (val types.PinRequest, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PinRequestKeyPrefix))

	b := store.Get(types.PinRequestKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemovePinRequest removes a pinRequest from the store
func (k Keeper) RemovePinRequest(
	ctx context.Context,
	index string,

) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PinRequestKeyPrefix))
	store.Delete(types.PinRequestKey(
		index,
	))
}

// GetAllPinRequest returns all pinRequest
func (k Keeper) GetAllPinRequest(ctx context.Context) (list []types.PinRequest) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PinRequestKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PinRequest
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
