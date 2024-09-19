package keeper

import (
	"context"

	"topchain/x/requester/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) RequestSubscription(goCtx context.Context, msg *types.MsgRequestSubscription) (*types.MsgRequestSubscriptionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var subscription = types.Subscription{
		Croid:    msg.Croid,
		Ammount:  msg.Ammount,
		Duration: msg.Duration,
	}
	hash := k.AddSubscription(ctx, subscription)

	return &types.MsgRequestSubscriptionResponse{Subscriptionid: hash}, nil
}
