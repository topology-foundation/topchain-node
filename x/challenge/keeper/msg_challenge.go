package keeper

import (
	"bytes"
	"context"
	"encoding/gob"

	"topchain/x/challenge/types"
	sTypes "topchain/x/subscription/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/google/uuid"
)

func (k msgServer) Challenge(goCtx context.Context, msg *types.MsgChallenge) (*types.MsgChallengeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	requester, err := sdk.AccAddressFromBech32(msg.Challenger)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid challenger address")
	}

	totalChallengePrice := k.PricePerVertexChallenge(ctx, msg.Challenger, msg.ProviderId) * int64(len(msg.VerticesHashes))
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, requester, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin("top", totalChallengePrice)))
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to send coins to module account")
	}

	challenger := msg.Challenger
	provider_id := msg.ProviderId
	challenged_hashes := msg.VerticesHashes

	currentBlock := ctx.BlockHeight()
	for _, hash := range challenged_hashes {
		block, found := k.GetHashSubmissionBlock(ctx, hash)
		if !found {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "hash "+hash+" not found")
		}
		if currentBlock-block > ChallengePeriod {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "hash "+hash+" was submitted more than "+string(ChallengePeriod)+" blocks ago")
		}
	}

	id := uuid.NewString()
	hashes := sTypes.SetFrom(challenged_hashes...)
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(hashes)

	k.SetChallenge(ctx, types.Challenge{
		Id:               id,
		Challenger:       challenger,
		Provider:         provider_id,
		Amount:           uint64(totalChallengePrice),
		LastActive:       uint64(currentBlock),
		ChallengedHashes: buf.Bytes(),
	})

	return &types.MsgChallengeResponse{ChallengeId: id}, nil
}
