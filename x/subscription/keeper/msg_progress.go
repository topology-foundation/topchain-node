package keeper

import (
	"context"
	"fmt"

	"topchain/utils"
	challengeKeeper "topchain/x/challenge/keeper"
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
	currentEpoch := utils.ConvertBlockToEpoch(ctx.BlockHeight())
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
		k.SetObfuscatedProgress(ctx, subscriptionId, currentEpoch, obfuscatedVerticesHash)
		return &types.MsgSubmitProgressResponse{}, nil
	}

	// Validate that the obfuscated vertex hashes submitted in the previous epoch match the current vertex hashes
	obfuscatedProgressData, _ := k.GetObfuscatedProgress(ctx, subscriptionId)
	err := validateObfuscatedProgress(obfuscatedProgressData, submittedHashes, provider, currentEpoch)
	if err != nil {
		return nil, errorsmod.Wrap(err, "vertices hashes / obfuscated data validation failed")
	}

	progress, found := k.GetProgress(ctx, subscriptionId)
	if !found {
		hashesSet := types.SetFrom(submittedHashes...)
		for hash := range hashesSet {
			k.SetHashSubmissionEpoch(ctx, provider, hash, currentEpoch)
		}
		k.SetProgress(ctx, subscriptionId, hashesSet)

		k.AddProgressDealAtEpoch(ctx, subscription.DealId, provider, currentEpoch, uint64(len(hashesSet)))
		k.AddProgressEpochsProvider(ctx, provider, subscriptionId, currentEpoch)
		return &types.MsgSubmitProgressResponse{}, nil
	}

	initialProgressSize := len(progress)
	for _, hash := range submittedHashes {
		if !progress.Has(hash) {
			progress = progress.Add(hash)
			k.SetHashSubmissionEpoch(ctx, provider, hash, currentEpoch)
		}
	}

	// Add the new obfuscated progress hash to the obfuscated progress hash set
	k.SetObfuscatedProgress(ctx, subscriptionId, currentEpoch, obfuscatedVerticesHash)

	progressSize := len(progress) - initialProgressSize

	k.SetProgress(ctx, subscriptionId, progress)
	k.AddProgressDealAtEpoch(ctx, subscription.DealId, provider, currentEpoch, uint64(progressSize))
	k.AddProgressEpochsProvider(ctx, provider, subscriptionId, currentEpoch)
	k.SetProgressSize(ctx, subscriptionId, currentEpoch, progressSize)

	return &types.MsgSubmitProgressResponse{}, nil
}

func validateObfuscatedProgress(obfuscatedProgressData ObfuscatedProgressData, submittedHashes []string, provider string, epochNumber uint64) error {
	if epochNumber != obfuscatedProgressData.EpochNumber+1 {
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

func (k msgServer) ClaimRewards(goCtx context.Context, msg *types.MsgClaimRewards) (*types.MsgClaimRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	provider := msg.Provider
	subscriptionId := msg.SubscriptionId
	currentEpoch := utils.ConvertBlockToEpoch(ctx.BlockHeight())
	challengeWindow := utils.ConvertBlockToEpoch(challengeKeeper.ChallengePeriod)
	reward := int64(0)

	subscription, found := k.GetSubscription(ctx, subscriptionId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "subscription with id "+subscriptionId+" not found")
	}

	deal, found := k.GetDeal(ctx, subscription.DealId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "this should never happen")
	}

	if subscription.Provider != provider {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only the provider can claim reward")
	}

	if currentEpoch > deal.EndEpoch+utils.DEAL_EXPIRY_CLAIM_WINDOW {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "reward claim period for deal with deal_id"+deal.Id+"ended")
	}

	lastClaimedEpoch, found := k.GetProviderLastRewardClaimedEpoch(ctx, provider, subscriptionId)
	// if the provider is claiming for the first time, `found` is false. In this case, start checking for
	// rewards from the subscription startBlock.
	if !found {
		lastClaimedEpoch = subscription.StartEpoch - 1
	}

	if currentEpoch < lastClaimedEpoch+challengeWindow {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "wait until challenge period elapses before claiming")
	}

	providerProgressEpochs, found := k.GetProgressEpochsProvider(ctx, provider, subscriptionId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "cannot claim reward without submitting progress")
	}

	// loop until the most recent block that has elapsed the challenge window.
	lastEligibleEpoch := min(currentEpoch-challengeWindow, subscription.EndEpoch)
	for epoch := lastClaimedEpoch + 1; epoch <= lastEligibleEpoch; epoch++ {
		// only compute rewards for blocks that the provider submitted progress
		if providerProgressEpochs.Has(epoch) {
			// get the progress made by all the deal subscribers at `block`
			progressDeal, found := k.GetProgressDealAtEpoch(ctx, subscription.DealId, epoch)
			if !found {
				return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "this should not happen!")
			}
			epochReward := k.CalculateEpochReward(deal)
			for _, progress := range progressDeal.Progress {
				if progress.Provider == provider {
					fmt.Println("block reward {}, progressSize: {}, progressTotal: {}", float64(epochReward), float64(progress.Size), float64(progressDeal.Total))
					reward += int64(float64(epochReward) * float64(progress.Size) / float64(progressDeal.Total))
					break
				}
			}
			providerProgressEpochs.Remove(epoch)
		}
	}
	k.SetProviderLastRewardClaimedEpoch(ctx, provider, subscriptionId, lastEligibleEpoch)
	k.SetProgressEpochsProvider(ctx, providerProgressEpochs, provider, subscriptionId)
	// send payout
	fmt.Println("reward ***************", reward)
	k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(provider), sdk.NewCoins(sdk.NewInt64Coin("top", reward)))
	deal.AvailableAmount -= uint64(reward)

	k.SetDeal(ctx, deal)
	return &types.MsgClaimRewardsResponse{}, nil
}
