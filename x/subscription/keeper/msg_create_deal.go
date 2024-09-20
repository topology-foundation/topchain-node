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
