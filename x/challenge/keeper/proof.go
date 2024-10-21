package keeper

import (
	"topchain/x/challenge/types"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetProof(ctx sdk.Context, challengeId string, vertex types.Vertex) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.GetProofStoreKey(challengeId))

	appendedValue := k.cdc.MustMarshal(&vertex)
	store.Set([]byte(vertex.Hash), appendedValue)
}

func (k Keeper) GetProof(ctx sdk.Context, challengeId string, vertexHash string) (vertex types.Vertex, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.GetProofStoreKey(challengeId))

	vertexBytes := store.Get([]byte(vertexHash))
	if vertexBytes == nil {
		return vertex, false
	}

	k.cdc.MustUnmarshal(vertexBytes, &vertex)
	return vertex, true
}
