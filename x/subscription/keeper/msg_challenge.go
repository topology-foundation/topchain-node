package keeper

import (
	"context"

	"topchain/x/subscription/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) Challenge(goCtx context.Context, msg *types.MsgChallenge) (*types.MsgChallengeResponse, error) {
	_ = sdk.UnwrapSDKContext(goCtx)

	return &types.MsgChallengeResponse{}, nil
}
