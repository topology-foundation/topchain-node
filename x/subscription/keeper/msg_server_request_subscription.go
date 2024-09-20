package keeper

import (
	"context"

	"topchain/x/subscription/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) RequestSubscription(goCtx context.Context, msg *types.MsgRequestSubscription) (*types.MsgRequestSubscriptionResponse, error) {
	_ = sdk.UnwrapSDKContext(goCtx)

	return nil, nil
}
