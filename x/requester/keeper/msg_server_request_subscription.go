package keeper

import (
	"context"

	"topchain/x/requester/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) RequestSubscription(goCtx context.Context, msg *types.MsgRequestSubscription) (*types.MsgRequestSubscriptionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgRequestSubscriptionResponse{}, nil
}
