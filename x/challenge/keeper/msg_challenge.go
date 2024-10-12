package keeper

import (
	"context"

	"topchain/x/challenge/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) Challenge(goCtx context.Context, msg *types.MsgChallenge) (*types.MsgChallengeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	requester, err := sdk.AccAddressFromBech32(msg.Challenger)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid challenger address")
	}

	pricePerVertexChallenge := k.PricePerVertexChallenge(ctx, msg.Challenger, msg.ProviderId)
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, requester, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin("top", int64(len(msg.VerticesHashes))*pricePerVertexChallenge)))
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to send coins to module account")
	}

	// challenger := msg.Challenger
	// provider_id := msg.ProviderId
	// challenged_hashes := msg.VerticesHashes

	return &types.MsgChallengeResponse{}, nil
}
