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
	EpochNumber int64
	Hash        string
}

type ProgressTuple struct {
	Provider string
	Size     int64
}

type ProgressDeal struct {
	Total    int64
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

func (k Keeper) SetObfuscatedProgress(ctx sdk.Context, subscription string, epochNumber int64, hash string) {
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

func (k Keeper) SetProgressSize(ctx sdk.Context, subscription string, block int64, size int) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.GetProgressSizeStoreKey(subscription))

	store.Set(sdk.Uint64ToBigEndian(uint64(block)), sdk.Uint64ToBigEndian(uint64(size)))
}

func (k Keeper) GetProgressSize(ctx sdk.Context, subscription string, block int64) (size int, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.GetProgressSizeStoreKey(subscription))

	sizeBytes := store.Get(sdk.Uint64ToBigEndian(uint64(block)))
	if sizeBytes == nil {
		return size, false
	}

	return int(sdk.BigEndianToUint64(sizeBytes)), true
}

func (k Keeper) AddProgressDealAtBlock(ctx sdk.Context, deal string, provider string, block int64, size int64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.GetProgressDealStoreKey(deal))

	var progressDeal ProgressDeal
	// Fetch the progressDeal from the store or create a new if it doesn't exist.
	if progressBytes := store.Get(sdk.Uint64ToBigEndian(uint64(block))); progressBytes == nil {
		progressDeal = ProgressDeal{
			Total:    int64(0),
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
	store.Set(sdk.Uint64ToBigEndian(uint64(block)), buf.Bytes())
}

func (k Keeper) GetProgressDealAtBlock(ctx sdk.Context, deal string, block int64) (ProgressDeal, bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.GetProgressDealStoreKey(deal))

	// Retrieve the progress deal data for the specified block
	progressBytes := store.Get(sdk.Uint64ToBigEndian(uint64(block)))
	if progressBytes == nil {
		return ProgressDeal{}, false
	}

	// Decode the retrieved bytes into a ProgressDeal struct
	var progressDeal ProgressDeal
	buf := bytes.NewBuffer(progressBytes)
	gob.NewDecoder(buf).Decode(&progressDeal)

	return progressDeal, true
}

func (k Keeper) AddProgressBlocksProvider(ctx sdk.Context, provider string, subscription string, block int64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.GetProgressBlocksProviderKey(provider))

	var blocks types.Set[int64]
	if blocksBytes := store.Get([]byte(subscription)); blocksBytes == nil {
		blocks = types.Set[int64]{}
	} else {
		buf := bytes.NewBuffer(blocksBytes)
		gob.NewDecoder(buf).Decode(&blocks)
	}

	blocks.Add(block)

	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(blocks)
	store.Set([]byte(subscription), buf.Bytes())
}

func (k Keeper) GetProgressBlocksProvider(ctx sdk.Context, provider string, subscription string) (blocks types.Set[int64], found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.GetProgressBlocksProviderKey(provider))

	blocksBytes := store.Get([]byte(subscription))
	if blocksBytes == nil {
		return blocks, false
	}

	buf := bytes.NewBuffer(blocksBytes)
	gob.NewDecoder(buf).Decode(&blocks)
	return blocks, true
}

func (k Keeper) SetProviderLastRewardClaimedBlock(ctx sdk.Context, provider string, subscription string, block int64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.GetSubscriptionProviderLastClaimedKey(provider))

	store.Set([]byte(subscription), sdk.Uint64ToBigEndian(uint64(block)))
}

func (k Keeper) GetProviderLastRewardClaimedBlock(ctx sdk.Context, provider, subscription string) (block int64, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.GetSubscriptionProviderLastClaimedKey(provider))

	blockBytes := store.Get([]byte(subscription))
	if blockBytes == nil {
		return block, false
	}

	return int64(sdk.BigEndianToUint64(blockBytes)), true
}

func (k Keeper) SetHashSubmissionBlock(ctx sdk.Context, provider string, hash string, block int64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.GetHashSubmissionBlockStoreKey(provider))

	store.Set([]byte(hash), sdk.Uint64ToBigEndian(uint64(block)))
}

func (k Keeper) GetHashSubmissionBlock(ctx sdk.Context, provider string, hash string) (block int64, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.GetHashSubmissionBlockStoreKey(provider))

	blockBytes := store.Get([]byte(hash))
	if blockBytes == nil {
		return block, false
	}

	return int64(sdk.BigEndianToUint64(blockBytes)), true
}
