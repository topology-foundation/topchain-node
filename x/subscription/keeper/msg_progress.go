package keeper

import (
	"context"

	"topchain/x/subscription/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"golang.org/x/crypto/sha3"
)

func (k msgServer) SubmitProgress(goCtx context.Context, msg *types.MsgSubmitProgress) (*types.MsgSubmitProgressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	provider := msg.Provider
	subscriptionId := msg.SubscriptionId
	obfuscatedVerticesHash := msg.ObfuscatedVerticesHash
	blockHeight := ctx.BlockHeight()
	epochNumber := blockHeight / EPOCH_SIZE
	submittedHashes := msg.PreviousVerticesHashes

	subscription, found := k.GetSubscription(ctx, subscriptionId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "subscription with id "+subscriptionId+" not found")
	}
	if subscription.Provider != provider {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only the provider can submit progress")
	}

	// this is the first obfuscated progress batch submission
	if len(submittedHashes) == 0 {
		k.SetObfuscatedProgress(ctx, subscriptionId, epochNumber, obfuscatedVerticesHash)
		return &types.MsgSubmitProgressResponse{}, nil
	}

	// Validate that the obfuscated vertex hashes submitted in the previous epoch match the current vertex hashes
	obfuscatedProgressData, _ := k.GetObfuscatedProgress(ctx, subscriptionId)
	err := validateObfuscatedProgress(obfuscatedProgressData, submittedHashes, provider, epochNumber)
	if err != nil {
		return nil, errorsmod.Wrap(err, "vertices hashes / obfuscated data validation failed")
	}

	progress, found := k.GetProgress(ctx, subscriptionId)
	if !found {
		hashesSet := types.SetFrom(submittedHashes...)
		for hash := range hashesSet {
			k.SetHashSubmissionBlock(ctx, provider, hash, blockHeight)
		}
		k.SetProgress(ctx, subscriptionId, hashesSet)
		k.SetProgressSize(ctx, subscriptionId, blockHeight, len(hashesSet))
		return &types.MsgSubmitProgressResponse{}, nil
	}

	initialProgressSize := len(progress)
	for _, hash := range submittedHashes {
		if !progress.Has(hash) {
			progress = progress.Add(hash)
			k.SetHashSubmissionBlock(ctx, provider, hash, blockHeight)
		}
	}

	// Add the new obfuscated progress hash to the obfuscated progress hash set
	k.SetObfuscatedProgress(ctx, subscriptionId, epochNumber, obfuscatedVerticesHash)

	progressSize := len(progress) - initialProgressSize

	k.SetProgress(ctx, subscriptionId, progress)
	k.SetProgressSize(ctx, subscriptionId, blockHeight, progressSize)

	return &types.MsgSubmitProgressResponse{}, nil
}

func validateObfuscatedProgress(obfuscatedProgressData ObfuscatedProgressData, submittedHashes []string, provider string, epochNumber int64) error {
	if epochNumber != obfuscatedProgressData.EpochNumber {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Revealing vertices hashes exceeded epoch duration")
	}
	hasher := sha3.New256()
	for _, hash := range submittedHashes {
		hasher.Write([]byte(hash))
	}
	hasher.Write([]byte(provider))
	hashBytes := hasher.Sum(nil)
	obfuscatedHash := string(hashBytes)
	if obfuscatedProgressData.Hash != obfuscatedHash {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "submitted vertices hashes don't match previous epoch obfuscated hash")
	}
	return nil
}
