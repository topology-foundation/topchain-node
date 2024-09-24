package keeper

import (
	"context"
	"topchain/x/subscription/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/google/uuid"
)

func (k msgServer) CreateDeal(goCtx context.Context, msg *types.MsgCreateDeal) (*types.MsgCreateDealResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	id := uuid.NewString()
	var deal = types.Deal{
		Id:              id,
		Requester:       msg.Requester,
		CroId:           msg.CroId,
		SubscriptionIds: []string{},
		Status:          types.Deal_SCHEDULED,
		InitialAmount:   msg.Amount,
		AvailableAmount: msg.Amount,
		StartBlock:      msg.StartBlock,
		EndBlock:        msg.EndBlock,
	}

	k.SetDeal(ctx, deal)

	requester, err := sdk.AccAddressFromBech32(msg.Requester)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid requester address")
	}
	sdkError := k.bankKeeper.SendCoinsFromAccountToModule(ctx, requester, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin("top", int64(msg.Amount))))
	if sdkError != nil {
		return nil, errorsmod.Wrap(sdkError, "failed to send coins to module account")
	}

	return &types.MsgCreateDealResponse{DealId: id}, nil
}

func (k msgServer) CancelDeal(goCtx context.Context, msg *types.MsgCancelDeal) (*types.MsgCancelDealResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	deal, found := k.GetDeal(ctx, msg.DealId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "deal with id"+msg.DealId+"not found")
	}
	if msg.Requester != deal.Requester {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only the requester can cancel the deal")
	}
	if deal.Status == types.Deal_SCHEDULED || deal.Status == types.Deal_INITIALIZED {
		deal.Status = types.Deal_CANCELLED
		k.SetDeal(ctx, deal)
		// return the remaining amount to the requester
		k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(deal.Requester), sdk.NewCoins(sdk.NewInt64Coin("top", int64(deal.AvailableAmount))))
		return &types.MsgCancelDealResponse{}, nil
	}
	if deal.Status == types.Deal_INACTIVE || deal.Status == types.Deal_ACTIVE {
		deal.Status = types.Deal_CANCELLED
		k.SetDeal(ctx, deal)
		for _, subscriptionId := range deal.SubscriptionIds {
			subscription, found := k.GetSubscription(ctx, subscriptionId)
			if !found {
				return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "subscription with id"+subscriptionId+"not found")
			}
			subscription.EndBlock = uint64(ctx.BlockHeight())
			k.SetSubscription(ctx, subscription)
		}
		// return the remaining amount to the requester
		k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(deal.Requester), sdk.NewCoins(sdk.NewInt64Coin("top", int64(deal.AvailableAmount))))
	}

	return &types.MsgCancelDealResponse{}, nil
}

func (k msgServer) UpdateDeal(goCtx context.Context, msg *types.MsgUpdateDeal) (*types.MsgUpdateDealResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	deal, found := k.GetDeal(ctx, msg.DealId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "deal with id"+msg.DealId+"not found")
	}
	if msg.Requester != deal.Requester {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only the requester can update the deal")
	}
	if ctx.BlockHeight() < int64(deal.StartBlock) {
		if msg.Amount != 0 {
			deal.InitialAmount = msg.Amount
			deal.AvailableAmount = msg.Amount
		}
		if msg.StartBlock != 0 {
			if int64(msg.StartBlock) < ctx.BlockHeight() {
				return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "start block must be greater than current block height")
			}
			deal.StartBlock = msg.StartBlock
		}
		if msg.EndBlock != 0 {
			if int64(msg.EndBlock) < ctx.BlockHeight() {
				return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "end block must be greater than current block height")
			}
			deal.EndBlock = msg.EndBlock
		}
	} else {
		if msg.Amount != 0 {
			if msg.Amount < deal.InitialAmount {
				return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount must be greater than initial amount")
			}
			deal.AvailableAmount = msg.Amount
		}
		if msg.StartBlock != 0 {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "cannot update start block after deal has started")
		}
		if msg.EndBlock != 0 {
			if int64(msg.EndBlock) < ctx.BlockHeight() {
				return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "end block must be greater than current block height")
			}
			deal.EndBlock = msg.EndBlock
		}
	}

	k.SetDeal(ctx, deal)

	return &types.MsgUpdateDealResponse{}, nil
}

func (k msgServer) JoinDeal(goCtx context.Context, msg *types.MsgJoinDeal) (*types.MsgJoinDealResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	deal, found := k.GetDeal(ctx, msg.DealId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "deal with id"+msg.DealId+"not found")
	}
	delegations := k.stakingKeeper.GetAllDelegations(ctx, sdk.AccAddress(msg.Provider))

	var totalStake int64 = 0
	for _, delegation := range delegations {
		totalStake += delegation.GetShares().TruncateInt64()
	}

	// need a formula to determine the necessary amount to join the deal, currently always accepted
	if totalStake < k.CalculateMinimumStake(ctx, deal) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "insufficient stake to join deal")
	}

	id := uuid.NewString()
	var subscriptionStartBlock uint64
	if deal.StartBlock < uint64(ctx.BlockHeight()) {
		subscriptionStartBlock = uint64(ctx.BlockHeight())
		deal.Status = types.Deal_ACTIVE
	} else {
		subscriptionStartBlock = deal.StartBlock
	}

	subsription := types.Subscription{
		Id:         id,
		DealId:     msg.DealId,
		Provider:   msg.Provider,
		StartBlock: subscriptionStartBlock,
		EndBlock:   deal.EndBlock,
	}
	k.SetSubscription(ctx, subsription)
	deal.SubscriptionIds = append(deal.SubscriptionIds, subsription.Id)

	k.SetDeal(ctx, deal)

	// TODO -> droak
	// rpc to TopologyNode

	return &types.MsgJoinDealResponse{}, nil
}

func (k msgServer) LeaveDeal(goCtx context.Context, msg *types.MsgLeaveDeal) (*types.MsgLeaveDealResponse, error) {
	_ = sdk.UnwrapSDKContext(goCtx)

	return &types.MsgLeaveDealResponse{}, nil
}
