package keeper

import (
	"context"

	"topchain/x/subscription/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"golang.org/x/crypto/sha3"
)

func (k msgServer) SubmitObfuscatedProgress(goCtx context.Context, msg *types.MsgSubmitObfuscatedProgress) (*types.MsgSubmitObfuscatedProgressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	provider := msg.Provider
	subscriptionId := msg.SubscriptionId
	subscription, found := k.GetSubscription(ctx, subscriptionId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "subscription with id "+subscriptionId+" not found")
	}
	if subscription.Provider != provider {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only the provider can submit progress")
	}

	obfuscatedHashes := msg.ObfuscatedVertexHashes
	hashesSet := types.SetFrom(obfuscatedHashes...)

	k.SetObfuscatedProgress(ctx, subscriptionId, ctx.BlockHeight(), hashesSet)

	return &types.MsgSubmitObfuscatedProgressResponse{}, nil

}

func (k msgServer) SubmitProgress(goCtx context.Context, msg *types.MsgSubmitProgress) (*types.MsgSubmitProgressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	provider := msg.Provider
	subscriptionId := msg.SubscriptionId
	subscription, found := k.GetSubscription(ctx, subscriptionId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "subscription with id "+subscriptionId+" not found")
	}
	if subscription.Provider != provider {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only the provider can submit progress")
	}

	// Validate that the obfuscated vertex hashes submitted in the previous block match the current vertex hashes
	obfuscatedProgress, found := k.GetObfuscatedProgress(ctx, subscriptionId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "obfuscated progress for subscription "+subscriptionId+" not found")
	}

	err := validateObfuscatedProgress(obfuscatedProgress, msg.VerticesHashes, provider, ctx.BlockHeight())
	if err != nil {
		return nil, errorsmod.Wrap(err, "vertex hashes don't match the obfuscated vertex hashes is invalid")
	}

	submittedHashes := msg.VerticesHashes

	progress, found := k.GetProgress(ctx, subscriptionId)
	if !found {
		hashesSet := types.SetFrom(submittedHashes...)
		for hash := range hashesSet {
			k.SetHashSubmissionBlock(ctx, provider, hash, ctx.BlockHeight())
		}
		k.SetProgress(ctx, subscriptionId, hashesSet)
		return &types.MsgSubmitProgressResponse{}, nil
	}

	initialProgressSize := len(progress)
	for _, hash := range submittedHashes {
		if !progress.Has(hash) {
			progress = progress.Add(hash)
			k.SetHashSubmissionBlock(ctx, provider, hash, ctx.BlockHeight())
		}
	}

	progressSize := len(progress) - initialProgressSize
	k.SetProgress(ctx, subscriptionId, progress)
	k.SetProgressSize(ctx, subscriptionId, ctx.BlockHeight(), progressSize)

	return &types.MsgSubmitProgressResponse{}, nil
}

func validateObfuscatedProgress(obfuscatedProgress ObfuscatedProgressData, submittedHashes []string, provider string, currentBlock int64) error {
	if obfuscatedProgress.BlockNumber != currentBlock-1 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "obfuscated progress is not for the previous block")
	}
	if len(obfuscatedProgress.Hashes) != len(submittedHashes) {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "obfuscated progress and submitted hashes have different lengths")
	}

	for _, hash := range submittedHashes {
		hasher := sha3.New256()
		hasher.Write([]byte(hash + provider))
		hashBytes := hasher.Sum(nil)
		if !obfuscatedProgress.Hashes.Has(string(hashBytes)) {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "submitted hash is not in the obfuscated progress")
		}
	}
	return nil
}
