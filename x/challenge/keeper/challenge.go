package keeper

import (
	"bytes"
	"encoding/gob"
	"topchain/x/challenge/types"
	sKeeper "topchain/x/subscription/keeper"
	sTypes "topchain/x/subscription/types"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ChallengePeriod  = 100
	InactivityPeriod = 100
)

type ChallengeHash struct {
	Hash  string
	Epoch uint64
}

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

func (k Keeper) GetHashSubmissionEpoch(ctx sdk.Context, provider string, hash string) (block uint64, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, sTypes.GetHashSubmissionBlockStoreKey(provider))

	epochBytes := store.Get([]byte(hash))
	if epochBytes == nil {
		return block, false
	}
	return sdk.BigEndianToUint64(epochBytes), true
}
func (k Keeper) SetDeal(ctx sdk.Context, deal sTypes.Deal) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(sTypes.DealKeyPrefix))

	appendedValue := k.cdc.MustMarshal(&deal)
	store.Set([]byte(deal.Id), appendedValue)

	providerStore := prefix.NewStore(storeAdapter, sTypes.GetRequesterStoreKey(deal.Requester))
	providerStore.Set([]byte(deal.Id), []byte{})
}

func (k Keeper) GetDeal(ctx sdk.Context, dealId string) (deal sTypes.Deal, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(sTypes.DealKeyPrefix))
	dealBytes := store.Get([]byte(dealId))
	if dealBytes == nil {
		return deal, false
	}
	k.cdc.MustUnmarshal(dealBytes, &deal)
	return deal, true
}

func (k Keeper) GetProgressDealAtEpoch(ctx sdk.Context, deal string, epoch uint64) (sKeeper.ProgressDeal, bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, sTypes.GetProgressDealStoreKey(deal))
	// Retrieve the progress deal data for the specified block
	progressBytes := store.Get(sdk.Uint64ToBigEndian(epoch))
	if progressBytes == nil {
		return sKeeper.ProgressDeal{}, false
	}
	// Decode the retrieved bytes into a ProgressDeal struct
	var progressDeal sKeeper.ProgressDeal
	buf := bytes.NewBuffer(progressBytes)
	gob.NewDecoder(buf).Decode(&progressDeal)
	return progressDeal, true
}

func (k Keeper) SetProgressDealAtEpoch(ctx sdk.Context, deal string, epoch uint64, progressDeal sKeeper.ProgressDeal) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, sTypes.GetProgressDealStoreKey(deal))

	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(progressDeal)
	store.Set(sdk.Uint64ToBigEndian(epoch), buf.Bytes())
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

func (k Keeper) isChallengeExpired(ctx sdk.Context, challenge types.Challenge) bool {
	return challenge.LastActive+InactivityPeriod >= uint64(ctx.BlockHeight())
}
