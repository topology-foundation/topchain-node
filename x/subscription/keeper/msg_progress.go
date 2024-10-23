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
	obfuscatedHash := msg.obfuscatedHash
	subscriptionId := msg.SubscriptionId

	// store the obfuscatedHash in the obfuscatedProgressStore
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
	obfuscatedProgressHashSet, found := k.GetObfuscatedProgress(ctx, subscriptionId, provider)
	if !found {
		hashesSet := types.SetFrom(msg.NewBatchObfuscatedProgressHash)
		k.SetObfuscatedProgress(ctx, subscriptionId, provider, hashesSet)
		// early return as this is the first obfuscated progress batch submission
		return &types.MsgSubmitProgressResponse{}, nil
	}

	submittedHashes := msg.PreviousBatchVerticesHashes

	obfuscatedProgressHashSet, err := validateAndUpdateObfuscatedProgress(obfuscatedProgressHashSet, submittedHashes, provider)
	if err != nil {
		return nil, errorsmod.Wrap(err, "obfuscated progress hash for the submitted vertices hashes set does not exist")
	}

	// Add the new obfuscated progress hash to the obfuscated progress hash set
	obfuscatedProgressHashSet = obfuscatedProgressHashSet.Add(msg.NewBatchObfuscatedProgressHash)
	k.SetObfuscatedProgress(ctx, subscriptionId, provider, obfuscatedProgressHashSet)

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

func validateAndUpdateObfuscatedProgress(obfuscatedProgressHashSet types.Set[string], submittedHashes []string, provider string) (types.Set[string], error) {

	hasher := sha3.New256()
	for _, hash := range submittedHashes {
		hasher.Write([]byte(hash))
	}
	hasher.Write([]byte(provider))
	hashBytes := hasher.Sum(nil)
	obfuscatedHash := string(hashBytes)
	if !obfuscatedProgressHashSet.Has(obfuscatedHash) {
		return obfuscatedProgressHashSet, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "submitted hash is not in the obfuscated progress")
	}
	// Remove the obfuscated hash from the set
	obfuscatedProgressHashSet = obfuscatedProgressHashSet.Remove(obfuscatedHash)
	return obfuscatedProgressHashSet, nil
}
