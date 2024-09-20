package keeper

import (
	"context"
	"topchain/x/subscription/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) LeaveDeal(goCtx context.Context, msg *types.MsgLeaveDeal) (*types.MsgLeaveDealResponse, error) {
	_ = sdk.UnwrapSDKContext(goCtx)

	return &types.MsgLeaveDealResponse{}, nil
}
