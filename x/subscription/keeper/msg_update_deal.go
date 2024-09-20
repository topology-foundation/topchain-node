package keeper

import (
	"context"
	"topchain/x/subscription/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateDeal(goCtx context.Context, msg *types.MsgUpdateDeal) (*types.MsgUpdateDealResponse, error) {
	_ = sdk.UnwrapSDKContext(goCtx)

	return &types.MsgUpdateDealResponse{}, nil
}
