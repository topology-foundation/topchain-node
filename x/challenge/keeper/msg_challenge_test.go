package keeper_test

import (
	"testing"

	"topchain/x/challenge/types"

	"github.com/stretchr/testify/require"
)

func TestSettleChallenge(t *testing.T) {
	_, ms, ctx, am := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)

	ctx = MockBlockHeight(ctx, am, 10)
	response, err := ms.Challenge(ctx, &types.MsgChallenge{Challenger: Alice, ProviderId: Bob, VerticesHashes: []string{}})
	require.NoError(t, err)

	challengeId := response.ChallengeId
	ctx = MockBlockHeight(ctx, am, 110)
	_, err = ms.SettleChallenge(ctx, &types.MsgSettleChallenge{Requester: Alice, ChallengeId: challengeId})

	// There should be an error as the challenge is not yet expired
	require.Error(t, err)
	require.Contains(t, err.Error(), "challenge is not yet expired")

	ctx = MockBlockHeight(ctx, am, 111)
	_, err = ms.SettleChallenge(ctx, &types.MsgSettleChallenge{Requester: Alice, ChallengeId: challengeId})

	// There should not be any error
	require.NoError(t, err)
}
