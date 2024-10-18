package keeper

import (
	"topchain/x/challenge/types"
	sTypes "topchain/x/subscription/types"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ChallengePeriod = 100
)

func (k Keeper) SetChallenge(ctx sdk.Context, challenge types.Challenge) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ChallengeKeyPrefix))

	appendedValue := k.cdc.MustMarshal(&challenge)
	store.Set([]byte(challenge.Id), appendedValue)
}

func (k Keeper) GetChallenge(ctx sdk.Context, challengeId string) (challenge types.Challenge, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ChallengeKeyPrefix))

	challengeBytes := store.Get([]byte(challengeId))
	if challengeBytes == nil {
		return challenge, false
	}

	k.cdc.MustUnmarshal(challengeBytes, &challenge)
	return challenge, true
}

func (k Keeper) GetHashSubmissionBlock(ctx sdk.Context, provider string, hash string) (block int64, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, sTypes.GetHashSubmissionBlockStoreKey(provider))

	blockBytes := store.Get([]byte(hash))
	if blockBytes == nil {
		return block, false
	}

	return int64(sdk.BigEndianToUint64(blockBytes)), true
}

func (k Keeper) PricePerVertexChallenge(ctx sdk.Context, challenger_address string, provider_id string) int64 {
	// TODO - devise a formula, set of arguments might be different
	return 1
}
