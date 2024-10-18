package keeper

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"

	"topchain/x/challenge/types"
	sTypes "topchain/x/subscription/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"crypto/sha256"

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

func (k msgServer) SubmitProof(goCtx context.Context, msg *types.MsgSubmitProof) (*types.MsgSubmitProofResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	challenge, found := k.GetChallenge(ctx, msg.ChallengeId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "challenge "+msg.ChallengeId+" not found")
	}
	if challenge.Provider != msg.Provider {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "unauthorized provider")
	}

	buf := bytes.NewBuffer(challenge.ChallengedHashes)
	var challengedHashes sTypes.Set[string]
	gob.NewDecoder(buf).Decode(&challengedHashes)

	for _, vertex := range msg.Vertices {
		if challengedHashes.Has(vertex.Hash) {
			// TODO - make sure this gives the same hashing output as in ts-topology
			stringified, err := json.Marshal(vertex)
			if err != nil {
				return nil, errorsmod.Wrap(err, "failed to marshal vertex")
			}
			computedHash := sha256.Sum256(stringified)

			if !bytes.Equal(computedHash[:], []byte(vertex.Hash)) {
				return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "hashes do not match")
			}

			challengedHashes.Remove(vertex.Hash)
		}
	}

	challenge.LastActive = uint64(ctx.BlockHeight())
	buf = &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(challengedHashes)
	challenge.ChallengedHashes = buf.Bytes()

	k.SetChallenge(ctx, challenge)

	return &types.MsgSubmitProofResponse{}, nil
}

func (k msgServer) RequestDependencies(goCtx context.Context, msg *types.MsgRequestDependencies) (*types.MsgRequestDependenciesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	challenge, found := k.GetChallenge(ctx, msg.ChallengeId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "challenge "+msg.ChallengeId+" not found")
	}
	if challenge.Challenger != msg.Challenger {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "unauthorized challenger")
	}

	requester, err := sdk.AccAddressFromBech32(msg.Challenger)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid challenger address")
	}

	fee := k.PricePerVertexChallenge(ctx, msg.Challenger, challenge.Provider) * int64(len(msg.VerticesHashes))
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, requester, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin("top", fee)))
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to send coins to module account")
	}

	buf := bytes.NewBuffer(challenge.ChallengedHashes)
	var challengedHashes sTypes.Set[string]
	gob.NewDecoder(buf).Decode(&challengedHashes)

	for _, hash := range msg.VerticesHashes {
		challengedHashes.Add(hash)
	}

	challenge.LastActive = uint64(ctx.BlockHeight())
	challenge.Amount += uint64(fee)

	buf = &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(challengedHashes)
	challenge.ChallengedHashes = buf.Bytes()

	k.SetChallenge(ctx, challenge)

	return &types.MsgRequestDependenciesResponse{}, nil
}
