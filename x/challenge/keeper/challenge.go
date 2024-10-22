package keeper

import (
	"topchain/x/challenge/types"
	sTypes "topchain/x/subscription/types"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ChallengePeriod  = 100
	InactivityPeriod = 100
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

func (k Keeper) RemoveChallenge(ctx sdk.Context, challengeId string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ChallengeKeyPrefix))

	store.Delete([]byte(challengeId))
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

// Iterate over all challenges and apply the given callback function
func (k Keeper) IterateChallenges(ctx sdk.Context, shouldBreak func(challenge types.Challenge) bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ChallengeKeyPrefix)).Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var challenge types.Challenge
		k.cdc.MustUnmarshal(iterator.Value(), &challenge)
		if shouldBreak(challenge) {
			break
		}
	}
}

func (k Keeper) PricePerVertexChallenge(ctx sdk.Context, challenger_address string, provider_id string) int64 {
	return 1
}
