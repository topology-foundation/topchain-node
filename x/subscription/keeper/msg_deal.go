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

	k.AddDeal(ctx, deal)

	return &types.MsgCreateDealResponse{DealId: id}, nil
}

func (k msgServer) CancelDeal(goCtx context.Context, msg *types.MsgCancelDeal) (*types.MsgCancelDealResponse, error) {
	_ = sdk.UnwrapSDKContext(goCtx)

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

	k.AddDeal(ctx, deal)

	return &types.MsgUpdateDealResponse{}, nil
}

func (k msgServer) JoinDeal(goCtx context.Context, msg *types.MsgJoinDeal) (*types.MsgJoinDealResponse, error) {
	_ = sdk.UnwrapSDKContext(goCtx)

	return &types.MsgJoinDealResponse{}, nil
}

func (k msgServer) LeaveDeal(goCtx context.Context, msg *types.MsgLeaveDeal) (*types.MsgLeaveDealResponse, error) {
	_ = sdk.UnwrapSDKContext(goCtx)

	return &types.MsgLeaveDealResponse{}, nil
}
