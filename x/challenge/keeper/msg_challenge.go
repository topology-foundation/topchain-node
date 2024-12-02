package keeper

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"

	topTypes "topchain/types"
	utils "topchain/utils"
	"topchain/x/challenge/types"
	sKeeper "topchain/x/subscription/keeper"
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

	currentBlock := uint64(ctx.BlockHeight())
	var hashes sTypes.Set[ChallengeHash]
	deal, found := k.GetDeal(ctx, msg.DealId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("Deal %s not found", msg.DealId))
	}

	for _, hash := range msg.VerticesHashes {
		epoch, found := k.GetHashSubmissionEpoch(ctx, msg.ProviderId, hash)
		block := deal.StartBlock + deal.EpochSize*epoch
		if !found {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("hash %s not found", hash))
		} else if currentBlock-block > ChallengePeriod {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("hash %s was submitted more than %d blocks ago", hash, ChallengePeriod))
		} else {
			hashes.Add(ChallengeHash{Hash: hash, Epoch: epoch})
		}
	}

	totalChallengePrice := k.PricePerVertexChallenge(ctx, msg.Challenger, msg.ProviderId) * int64(len(msg.VerticesHashes))
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, requester, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(topTypes.TokenDenom, totalChallengePrice)))
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to send coins to module account")
	}

	id := uuid.NewString()
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(hashes)

	k.SetChallenge(ctx, types.Challenge{
		Id:               id,
		Challenger:       msg.Challenger,
		Provider:         msg.ProviderId,
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
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("challenge %s not found", msg.ChallengeId))
	}
	if challenge.Provider != msg.Provider {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "unauthorized provider")
	}
	if k.isChallengeExpired(ctx, challenge) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "challenge is expired")
	}

	buf := bytes.NewBuffer(challenge.ChallengedHashes)
	var challengedHashes sTypes.Set[string]
	gob.NewDecoder(buf).Decode(&challengedHashes)

	for _, vertex := range msg.Vertices {
		if challengedHashes.Has(vertex.Hash) {
			vertexData := map[string]interface{}{
				"operation": vertex.Operation,
				"deps":      vertex.Dependencies,
				"nodeId":    vertex.NodeId,
			}
			stringified, err := json.Marshal(vertexData)
			if err != nil {
				return nil, errorsmod.Wrap(err, fmt.Sprintf("failed to marshal vertex with hash %s", vertex.Hash))
			}
			computedHash := sha256.Sum256(stringified)

			if !bytes.Equal(computedHash[:], []byte(vertex.Hash)) {
				return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("hash %s does not match the computed hash", vertex.Hash))
			}

			k.SetProof(ctx, msg.ChallengeId, *vertex)
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
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("challenge %s not found", msg.ChallengeId))
	}
	if challenge.Challenger != msg.Challenger {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "unauthorized challenger")
	}
	if k.isChallengeExpired(ctx, challenge) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "challenge is expired")
	}

	requester, err := sdk.AccAddressFromBech32(msg.Challenger)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid challenger address")
	}

	fee := k.PricePerVertexChallenge(ctx, msg.Challenger, challenge.Provider) * int64(len(msg.VerticesHashes))
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, requester, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(topTypes.TokenDenom, fee)))
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to send coins to module account")
	}

	deal, found := k.GetDeal(ctx, challenge.DealId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("Deal %s not found", challenge.DealId))
	}

	buf := bytes.NewBuffer(challenge.ChallengedHashes)
	var challengedHashes sTypes.Set[string]
	gob.NewDecoder(buf).Decode(&challengedHashes)

	currentBlock := uint64(ctx.BlockHeight())
	for _, hash := range msg.VerticesHashes {
		epoch, found := k.GetHashSubmissionEpoch(ctx, challenge.Provider, hash)
		block := deal.StartBlock + deal.EpochSize*epoch
		if !found {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("hash %s not found", hash))
		}
		if currentBlock-block > ChallengePeriod {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("hash %s was submitted more than %d blocks ago", hash, ChallengePeriod))
		} else {
			challengedHashes.Add(hash)
		}
	}

	challenge.LastActive = uint64(currentBlock)
	challenge.Amount += uint64(fee)

	buf = &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(challengedHashes)
	challenge.ChallengedHashes = buf.Bytes()

	k.SetChallenge(ctx, challenge)

	return &types.MsgRequestDependenciesResponse{}, nil
}

func (k msgServer) SettleChallenge(goCtx context.Context, msg *types.MsgSettleChallenge) (*types.MsgSettleChallengeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	challenge, found := k.GetChallenge(ctx, msg.ChallengeId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("challenge %s not found", msg.ChallengeId))
	}
	if msg.Requester != challenge.Challenger && msg.Requester != challenge.Provider {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "not the challenger or the provider")
	}
	if k.isChallengeExpired(ctx, challenge) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "challenge is not yet expired")
	}

	buf := bytes.NewBuffer(challenge.ChallengedHashes)
	var challengedHashes sTypes.Set[ChallengeHash]
	gob.NewDecoder(buf).Decode(&challengedHashes)

	deal, _ := k.GetDeal(ctx, challenge.DealId)

	coins := sdk.NewCoins(sdk.NewInt64Coin("top", int64(challenge.Amount)))
	if len(challengedHashes) == 0 {
		// all hashes were verified - send coins to provider, remove challenge
		k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(challenge.Provider), coins)
	} else {
		// remove the provider's progress for each of the epochs where challenge failed
		rewardConfiscated := uint64(0)
		var failedEpochs sTypes.Set[uint64]
		for challengedHash := range challengedHashes {

			if failedEpochs.Has(challengedHash.Epoch) {
				// this epoch was already processed
				continue
			}
			progressDeal, found := k.GetProgressDealAtEpoch(ctx, challenge.DealId, challengedHash.Epoch)
			if !found {
				return nil, errorsmod.Wrap(sdkerrors.ErrNotFound, "this should not happen")
			}
			isEpochRugged := updateProgressDeal(&progressDeal, challenge.Provider)
			if isEpochRugged {
				rewardConfiscated += deal.RewardPerEpoch
			}
			failedEpochs.Add(challengedHash.Epoch)
			// update the progressDeal
			k.SetProgressDealAtEpoch(ctx, challenge.DealId, challengedHash.Epoch, progressDeal)
		}
		// the confiscated reward goes back to the available reward pool
		deal.AvailableAmount += rewardConfiscated
		// reward per epoch is incremented for all the remaining epochs
		currentEpoch := utils.ConvertBlockToEpoch(uint64(ctx.BlockHeight())-deal.StartBlock, deal.EpochSize)
		deal.RewardPerEpoch += rewardConfiscated / (deal.NumEpochs - currentEpoch + 1)
		k.SetDeal(ctx, deal)
	}
	k.RemoveChallenge(ctx, challenge.Id)

	return &types.MsgSettleChallengeResponse{}, nil
}

func updateProgressDeal(progressDeal *sKeeper.ProgressDeal, provider string) bool {
	amount := uint64(0)

	for _, progress := range progressDeal.Progress {
		if progress.Provider == provider {
			amount = progress.Size
			progress.Size = 0
			break
		}
	}
	progressDeal.Total -= amount

	return progressDeal.Total == 0
}
