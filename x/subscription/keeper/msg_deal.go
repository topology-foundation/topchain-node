package keeper

import (
	"context"
	"topchain/x/subscription/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateDeal(goCtx context.Context, msg *types.MsgCreateDeal) (*types.MsgCreateDealResponse, error) {
	_ = sdk.UnwrapSDKContext(goCtx)

	return &types.MsgCreateDealResponse{}, nil
}

func (k msgServer) CancelDeal(goCtx context.Context, msg *types.MsgCancelDeal) (*types.MsgCancelDealResponse, error) {
	_ = sdk.UnwrapSDKContext(goCtx)

	return &types.MsgCancelDealResponse{}, nil
}

func (k msgServer) UpdateDeal(goCtx context.Context, msg *types.MsgUpdateDeal) (*types.MsgUpdateDealResponse, error) {
	_ = sdk.UnwrapSDKContext(goCtx)

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
