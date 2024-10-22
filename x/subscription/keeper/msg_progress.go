package keeper

import (
	"context"

	"topchain/x/subscription/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

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
