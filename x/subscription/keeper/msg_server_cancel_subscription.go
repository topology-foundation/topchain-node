package keeper

import (
	"context"

	"topchain/x/subscription/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CancelSubscription(goCtx context.Context, msg *types.MsgCancelSubscription) (*types.MsgCancelSubscriptionResponse, error) {
	_ = sdk.UnwrapSDKContext(goCtx)

	return nil, nil
}
