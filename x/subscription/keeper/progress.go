package keeper

import (
	"bytes"
	"encoding/gob"
	"topchain/x/subscription/types"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ObfuscatedProgressData represents the struct we want to store
type ObfuscatedProgressData struct {
	EpochNumber uint64
	Hash        string
}

type ProgressTuple struct {
	Provider string
	Size     uint64
}

type ProgressDeal struct {
	Total    uint64
	Progress []ProgressTuple
}

const EPOCH_SIZE = 10

func (k Keeper) SetProgress(ctx sdk.Context, subscription string, hashes types.Set[string]) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ProgressKeyPrefix))

	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(hashes)
	store.Set([]byte(subscription), buf.Bytes())
}

func (k Keeper) GetProgress(ctx sdk.Context, subscription string) (hashes types.Set[string], found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ProgressKeyPrefix))

	hashesBytes := store.Get([]byte(subscription))
	if hashesBytes == nil {
		return hashes, false
	}

	buf := bytes.NewBuffer(hashesBytes)
	gob.NewDecoder(buf).Decode(&hashes)
	return hashes, true
}

func (k Keeper) SetObfuscatedProgress(ctx sdk.Context, subscription string, epochNumber uint64, hash string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ProgressObfuscatedKeyPrefix))

	data := ObfuscatedProgressData{
		EpochNumber: epochNumber,
		Hash:        hash,
	}
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(data)
	store.Set([]byte(subscription), buf.Bytes())
}

func (k Keeper) GetObfuscatedProgress(ctx sdk.Context, subscription string) (data ObfuscatedProgressData, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ProgressObfuscatedKeyPrefix))

	hashesBytes := store.Get([]byte(subscription))
	if hashesBytes == nil {
		return data, false
	}

	buf := bytes.NewBuffer(hashesBytes)
	gob.NewDecoder(buf).Decode(&data)
	return data, true
}

func (k Keeper) SetProgressSize(ctx sdk.Context, subscription string, epoch uint64, size int) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.GetProgressSizeStoreKey(subscription))

	store.Set(sdk.Uint64ToBigEndian(epoch), sdk.Uint64ToBigEndian(uint64(size)))
}

func (k Keeper) GetProgressSize(ctx sdk.Context, subscription string, epoch uint64) (size int, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.GetProgressSizeStoreKey(subscription))

	sizeBytes := store.Get(sdk.Uint64ToBigEndian(epoch))
	if sizeBytes == nil {
		return size, false
	}

	return int(sdk.BigEndianToUint64(sizeBytes)), true
}

func (k Keeper) AddProgressDealAtEpoch(ctx sdk.Context, deal string, provider string, epoch uint64, size uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.GetProgressDealStoreKey(deal))

	var progressDeal ProgressDeal
	// Fetch the progressDeal from the store or create a new if it doesn't exist.
	if progressBytes := store.Get(sdk.Uint64ToBigEndian(epoch)); progressBytes == nil {
		progressDeal = ProgressDeal{
			Total:    0,
			Progress: []ProgressTuple{},
		}
	} else {
		buf := bytes.NewBuffer(progressBytes)
		gob.NewDecoder(buf).Decode(&progressDeal)
	}

	newProgress := ProgressTuple{
		Provider: provider,
		Size:     size,
	}
	progressDeal.Progress = append(progressDeal.Progress, newProgress)
	progressDeal.Total += size

	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(progressDeal)
	store.Set(sdk.Uint64ToBigEndian(epoch), buf.Bytes())
}

func (k Keeper) GetProgressDealAtEpoch(ctx sdk.Context, deal string, epoch uint64) (ProgressDeal, bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.GetProgressDealStoreKey(deal))

	// Retrieve the progress deal data for the specified block
	progressBytes := store.Get(sdk.Uint64ToBigEndian(epoch))
	if progressBytes == nil {
		return ProgressDeal{}, false
	}

	// Decode the retrieved bytes into a ProgressDeal struct
	var progressDeal ProgressDeal
	buf := bytes.NewBuffer(progressBytes)
	gob.NewDecoder(buf).Decode(&progressDeal)

	return progressDeal, true
}

func (k Keeper) AddProgressEpochsProvider(ctx sdk.Context, provider string, subscription string, epoch uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.GetProgressBlocksProviderKey(provider))

	var epochs types.Set[uint64]
	if epochsBytes := store.Get([]byte(subscription)); epochsBytes == nil {
		epochs = types.Set[uint64]{}
	} else {
		buf := bytes.NewBuffer(epochsBytes)
		gob.NewDecoder(buf).Decode(&epochs)
	}

	epochs.Add(epoch)

	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(epochs)
	store.Set([]byte(subscription), buf.Bytes())
}

func (k Keeper) SetProgressEpochsProvider(ctx sdk.Context, epochs types.Set[uint64], provider string, subscription string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.GetProgressBlocksProviderKey(provider))

	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(epochs)
	store.Set([]byte(subscription), buf.Bytes())
}

func (k Keeper) GetProgressEpochsProvider(ctx sdk.Context, provider string, subscription string) (epochs types.Set[uint64], found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.GetProgressBlocksProviderKey(provider))

	epochsBytes := store.Get([]byte(subscription))
	if epochsBytes == nil {
		return epochs, false
	}

	buf := bytes.NewBuffer(epochsBytes)
	gob.NewDecoder(buf).Decode(&epochs)
	return epochs, true
}

func (k Keeper) SetProviderLastRewardClaimedEpoch(ctx sdk.Context, provider string, subscription string, epoch uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.GetSubscriptionProviderLastClaimedKey(provider))

	store.Set([]byte(subscription), sdk.Uint64ToBigEndian(epoch))
}

func (k Keeper) GetProviderLastRewardClaimedEpoch(ctx sdk.Context, provider, subscription string) (epoch uint64, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.GetSubscriptionProviderLastClaimedKey(provider))

	epochBytes := store.Get([]byte(subscription))
	if epochBytes == nil {
		return epoch, false
	}

	return sdk.BigEndianToUint64(epochBytes), true
}

func (k Keeper) SetHashSubmissionEpoch(ctx sdk.Context, provider string, hash string, epoch uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.GetHashSubmissionBlockStoreKey(provider))

	store.Set([]byte(hash), sdk.Uint64ToBigEndian(epoch))
}

func (k Keeper) GetHashSubmissionEpoch(ctx sdk.Context, provider string, hash string) (epoch uint64, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.GetHashSubmissionBlockStoreKey(provider))

	epochBytes := store.Get([]byte(hash))
	if epochBytes == nil {
		return epoch, false
	}

	return sdk.BigEndianToUint64(epochBytes), true
}
