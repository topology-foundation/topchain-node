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
