package keeper

import (
	"context"

	"topchain/x/subscription/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) RequestSubscription(goCtx context.Context, msg *types.MsgRequestSubscription) (*types.MsgRequestSubscriptionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var subscription = types.Subscription{
		Creator:  msg.Creator,
		CroId:    msg.CroId,
		Amount:   msg.Amount,
		Duration: msg.Duration,
	}
	hash := k.AddSubscription(ctx, subscription)

	return &types.MsgRequestSubscriptionResponse{SubscriptionId: hash}, nil
}
