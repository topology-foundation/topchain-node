package keeper

import (
	"context"
	"fmt"
	"math"

	topTypes "topchain/types"
	"topchain/utils"
	"topchain/utils/validation"
	"topchain/x/subscription/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/google/uuid"
)

func validateMsgCreateDeal(msg *types.MsgCreateDeal) error {
	if err := validation.ValidateEpochSize(msg.EpochSize); err != nil {
		return err
	}
	if err := validation.ValidateNumEpochs(msg.NumEpochs); err != nil {
		return err
	}
	if err := validation.ValidatePositiveAmount(msg.Amount); err != nil {
		return err
	}
	if err := validation.ValidateNonEmptyString(msg.CroId); err != nil {
		return err
	}
	if err := validation.ValidateAddress(msg.Requester); err != nil {
		return err
	}
	return nil
}

func (k msgServer) CreateDeal(goCtx context.Context, msg *types.MsgCreateDeal) (*types.MsgCreateDealResponse, error) {
	err := validateMsgCreateDeal(msg)
	if err != nil {
		fmt.Println("error in validateMsgCreateDeal", err)
		return nil, err
	}

	rewardPerEpoch := msg.Amount / msg.NumEpochs

	ctx := sdk.UnwrapSDKContext(goCtx)
	id := uuid.NewString()
	deal := types.Deal{
		Id:              id,
		Requester:       msg.Requester,
		CroId:           msg.CroId,
		SubscriptionIds: []string{},
		Status:          types.Deal_SCHEDULED,
		TotalAmount:     msg.Amount,
		AvailableAmount: msg.Amount,
		StartBlock:      msg.StartBlock,
		EpochSize:       msg.EpochSize,
		NumEpochs:       msg.NumEpochs,
		RewardPerEpoch:  rewardPerEpoch,
		InitialFrontier: msg.InitialFrontier,
	}

	requester, err := sdk.AccAddressFromBech32(msg.Requester)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid requester address")
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, requester, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(topTypes.TokenDenom, int64(msg.Amount))))
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to send coins to module account")
	}

	k.SetDeal(ctx, deal)

	return &types.MsgCreateDealResponse{DealId: id}, nil
}

func validateMsgCancelDeal(msg *types.MsgCancelDeal) error {
	if err := validation.ValidateNonEmptyString(msg.DealId); err != nil {
		return err
	}
	if err := validation.ValidateAddress(msg.Requester); err != nil {
		return err
	}
	return nil
}

func (k msgServer) CancelDeal(goCtx context.Context, msg *types.MsgCancelDeal) (*types.MsgCancelDealResponse, error) {
	err := validateMsgCancelDeal(msg)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	deal, found := k.GetDeal(ctx, msg.DealId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "deal with id "+msg.DealId+" not found")
	}
	if msg.Requester != deal.Requester {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only the requester can cancel the deal")
	}
	if deal.Status == types.Deal_SCHEDULED || deal.Status == types.Deal_INITIALIZED {
		deal.Status = types.Deal_CANCELLED
		k.SetDeal(ctx, deal)
		// return the remaining amount to the requester
		k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(deal.Requester), sdk.NewCoins(sdk.NewInt64Coin(topTypes.TokenDenom, int64(deal.AvailableAmount))))
		return &types.MsgCancelDealResponse{}, nil
	}
	if deal.Status == types.Deal_INACTIVE || deal.Status == types.Deal_ACTIVE {
		deal.Status = types.Deal_CANCELLED
		k.SetDeal(ctx, deal)
		for _, subscriptionId := range deal.SubscriptionIds {
			subscription, found := k.GetSubscription(ctx, subscriptionId)
			if !found {
				return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "SHOULD NOT HAPPEN: subscription with id "+subscriptionId+" not found")
			}
			subscription.EndEpoch = utils.ConvertBlockToEpoch(uint64(ctx.BlockHeight())-deal.StartBlock, deal.EpochSize)
			k.SetSubscription(ctx, subscription)
		}
		// return the remaining amount to the requester
		k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(deal.Requester), sdk.NewCoins(sdk.NewInt64Coin(topTypes.TokenDenom, int64(deal.AvailableAmount))))
	}

	return &types.MsgCancelDealResponse{}, nil
}

func (k msgServer) UpdateDeal(goCtx context.Context, msg *types.MsgUpdateDeal) (*types.MsgUpdateDealResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	deal, found := k.GetDeal(ctx, msg.DealId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "deal with id "+msg.DealId+" not found")
	}
	requester, err := sdk.AccAddressFromBech32(msg.Requester)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid requester address")
	}
	if msg.Requester != deal.Requester {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only the requester can update the deal")
	}
	// Amount, StartBlock, and NumEpochs are optional arguments an default to 0 if not provided
	currentBlock := uint64(ctx.BlockHeight())
	currentEpoch := utils.ConvertBlockToEpoch(currentBlock-deal.StartBlock, deal.EpochSize)
	if currentBlock < deal.StartBlock {
		if msg.Amount != 0 {
			if msg.Amount < deal.TotalAmount {
				amountToReturn := deal.TotalAmount - msg.Amount
				k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, requester, sdk.NewCoins(sdk.NewInt64Coin(topTypes.TokenDenom, int64(amountToReturn))))
			} else if msg.Amount > deal.TotalAmount {
				amountToDeposit := msg.Amount - deal.TotalAmount
				sdkError := k.bankKeeper.SendCoinsFromAccountToModule(ctx, requester, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(topTypes.TokenDenom, int64(amountToDeposit))))
				if sdkError != nil {
					return nil, errorsmod.Wrap(sdkError, "failed to send coins to module account")
				}
			}
			deal.TotalAmount = msg.Amount
			deal.AvailableAmount = msg.Amount
		}
		if msg.StartBlock != 0 {
			if msg.StartBlock < currentBlock {
				return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "start block must be greater than current block")
			}
			deal.StartBlock = msg.StartBlock
		}
		if msg.NumEpochs != 0 {
			deal.NumEpochs = msg.NumEpochs
		}
		deal.RewardPerEpoch = deal.TotalAmount / deal.NumEpochs
	} else {
		amountToDeposit := uint64(0)
		if msg.Amount != 0 {
			if msg.Amount < deal.TotalAmount {
				return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount must be greater than initial amount")
			}
			amountToDeposit = msg.Amount - deal.TotalAmount
			requester, err := sdk.AccAddressFromBech32(msg.Requester)
			if err != nil {
				return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid requester address")
			}
			sdkError := k.bankKeeper.SendCoinsFromAccountToModule(ctx, requester, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(topTypes.TokenDenom, int64(amountToDeposit))))
			if sdkError != nil {
				return nil, errorsmod.Wrap(sdkError, "failed to send coins to module account")
			}
			deal.TotalAmount = msg.Amount
			deal.AvailableAmount += amountToDeposit
		}
		if msg.StartBlock != 0 {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "cannot update start block after deal has started")
		}
		if msg.NumEpochs != 0 {
			if msg.NumEpochs <= currentEpoch {
				return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "cannot update num epochs to a value less than or equal to the current epoch")
			}
			deal.NumEpochs = msg.NumEpochs
		}
		// the reward amount topup only applies to epochs >= currentEpoch
		deal.RewardPerEpoch = deal.RewardPerEpoch + amountToDeposit/(deal.NumEpochs-currentEpoch+1)
	}

	k.SetDeal(ctx, deal)

	return &types.MsgUpdateDealResponse{}, nil
}

func validateMsgIncrementDealAmount(msg *types.MsgIncrementDealAmount) error {
	if err := validation.ValidatePositiveAmount(msg.Amount); err != nil {
		return err
	}
	if err := validation.ValidateNonEmptyString(msg.DealId); err != nil {
		return err
	}
	if err := validation.ValidateAddress(msg.Requester); err != nil {
		return err
	}
	return nil
}

func (k msgServer) IncrementDealAmount(goCtx context.Context, msg *types.MsgIncrementDealAmount) (*types.MsgIncrementDealAmountResponse, error) {
	if err := validateMsgIncrementDealAmount(msg); err != nil {
		return nil, err
	}

	requester, err := sdk.AccAddressFromBech32(msg.Requester)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid requester address")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	deal, found := k.GetDeal(ctx, msg.DealId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "deal with id "+msg.DealId+" not found")
	}
	currentEpoch := utils.ConvertBlockToEpoch(uint64(ctx.BlockHeight())-deal.StartBlock, deal.EpochSize)
	if msg.Requester != deal.Requester {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only the requester can increment the deal amount")
	}

	if k.IsDealUnavailable(deal.Status) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "cannot topup the expired deal with id "+msg.DealId)
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, requester, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(topTypes.TokenDenom, int64(msg.Amount))))
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to send coins to module account")
	}
	deal.TotalAmount += msg.Amount
	deal.AvailableAmount += msg.Amount

	// the reward amount topup only applies to epochs >= currentEpoch
	deal.RewardPerEpoch = deal.RewardPerEpoch + msg.Amount/(deal.NumEpochs-currentEpoch+1)

	k.SetDeal(ctx, deal)
	return &types.MsgIncrementDealAmountResponse{}, nil
}

func validateMsgJoinDeal(msg *types.MsgJoinDeal) error {
	if err := validation.ValidateNonEmptyString(msg.DealId); err != nil {
		return err
	}
	if err := validation.ValidateAddress(msg.Provider); err != nil {
		return err
	}
	return nil
}

func (k msgServer) JoinDeal(goCtx context.Context, msg *types.MsgJoinDeal) (*types.MsgJoinDealResponse, error) {
	err := validateMsgJoinDeal(msg)
	if err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	deal, found := k.GetDeal(ctx, msg.DealId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "deal with id "+msg.DealId+" not found")
	}

	if k.IsDealUnavailable(deal.Status) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidVersion, "deal with id "+msg.DealId+"is not available to join")
	}

	if k.DealHasProvider(ctx, deal, msg.Provider) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidVersion, "provider is already subscribed to the deal with id "+msg.DealId)
	}

	delegations, err := k.stakingKeeper.GetDelegatorDelegations(ctx, sdk.AccAddress(msg.Provider), math.MaxUint16)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "provider "+msg.Provider+" not found")
	}

	var totalStake int64 = 0
	for _, delegation := range delegations {
		totalStake += delegation.GetShares().TruncateInt64()
	}

	// need a formula to determine the necessary amount to join the deal, currently always accepted
	if totalStake < k.CalculateMinimumStake(ctx, deal) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "insufficient stake to join deal")
	}

	id := uuid.NewString()
	var subscriptionStartEpoch uint64
	currentBlock := uint64(ctx.BlockHeight())

	currentEpoch := (currentBlock - deal.StartBlock + 1) / deal.EpochSize

	if deal.StartBlock < currentBlock {
		subscriptionStartEpoch = currentEpoch
		deal.Status = types.Deal_ACTIVE
	} else {
		// start from the first epoch
		subscriptionStartEpoch = 0
	}

	subscription := types.Subscription{
		Id:         id,
		DealId:     msg.DealId,
		Provider:   msg.Provider,
		StartEpoch: subscriptionStartEpoch,
		EndEpoch:   (deal.EpochSize*deal.NumEpochs + 1) / deal.EpochSize,
	}
	k.SetSubscription(ctx, subscription)
	deal.SubscriptionIds = append(deal.SubscriptionIds, subscription.Id)

	k.SetDeal(ctx, deal)

	return &types.MsgJoinDealResponse{SubscriptionId: id}, nil
}

func validateMsgLeaveDeal(msg *types.MsgLeaveDeal) error {
	if err := validation.ValidateNonEmptyString(msg.DealId); err != nil {
		return err
	}
	if err := validation.ValidateAddress(msg.Provider); err != nil {
		return err
	}
	return nil
}

func (k msgServer) LeaveDeal(goCtx context.Context, msg *types.MsgLeaveDeal) (*types.MsgLeaveDealResponse, error) {
	err := validateMsgLeaveDeal(msg)
	if err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	deal, found := k.GetDeal(ctx, msg.DealId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "deal with id "+msg.DealId+" not found")
	}

	isSubscribed := false
	for _, subscriptionId := range deal.SubscriptionIds {
		subscription, found := k.GetSubscription(ctx, subscriptionId)
		if !found {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "SHOULD NOT HAPPEN: subscription with id "+subscriptionId+" not found")
		}
		if subscription.Provider == msg.Provider {
			isSubscribed = true
			subscription.EndEpoch = utils.ConvertBlockToEpoch(uint64(ctx.BlockHeight())-deal.StartBlock, deal.EpochSize)
			k.SetSubscription(ctx, subscription)
		}
	}

	if !isSubscribed {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "provider must be subscribed to the deal with id "+msg.DealId+" to leave it")
	}

	if !k.IsDealActive(ctx, deal) {
		deal.Status = types.Deal_INACTIVE
		k.SetDeal(ctx, deal)
	}

	return &types.MsgLeaveDealResponse{}, nil
}
