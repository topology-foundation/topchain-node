package keeper

import (
	"context"
	"fmt"

	"topchain/x/requester/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CancelSubscription(goCtx context.Context, msg *types.MsgCancelSubscription) (*types.MsgCancelSubscriptionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	subscription, found := k.GetSubscription(ctx, msg.Subscriptionid)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("Key %s doesn't exist", msg.Subscriptionid))
	}
	if msg.Creator != subscription.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Incorrect owner")
	}
	k.RemoveSubscription(ctx, msg.Subscriptionid)

	return &types.MsgCancelSubscriptionResponse{}, nil
}
