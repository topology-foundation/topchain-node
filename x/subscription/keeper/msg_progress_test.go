package keeper_test

import (
	"testing"

	"topchain/x/subscription/types"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/sha3"
)

func ObfuscatedDataHashHelper(verticesHashes []string, provider string) string {
	hasher := sha3.New256()
	for _, hash := range verticesHashes {
		hasher.Write([]byte(hash))
	}
	hasher.Write([]byte(provider))
	hashBytes := hasher.Sum(nil)
	return string(hashBytes)
}

func TestSubmitProgress(t *testing.T) {
	_, ms, ctx, am := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)

	response, err := ms.CreateDeal(ctx, &types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 10000, StartBlock: 10, EndBlock: 20})
	require.NoError(t, err)

	dealId := response.DealId

	// Jump to block 12 to initiate the deal
	ctx = MockBlockHeight(ctx, am, 12)
	// Provider joins the deal after it is initiated
	joinDeal := types.MsgJoinDeal{Provider: Bob, DealId: dealId}
	joinResponse, err := ms.JoinDeal(ctx, &joinDeal)
	require.NoError(t, err)

	subscriptionId := joinResponse.SubscriptionId
	providerId := Bob

	// create mock vertices hashes and the corresponding obfuscated hash
	verticesHashes1 := []string{"000", "111", "222", "333", "444", "555"}
	obfuscatedHash1 := ObfuscatedDataHashHelper(verticesHashes1, providerId)

	// submit progress
	_, err = ms.SubmitProgress(ctx, &types.MsgSubmitProgress{Provider: providerId, SubscriptionId: subscriptionId, ObfuscatedVerticesHash: obfuscatedHash1})
	// There should not be any error
	require.NoError(t, err)
	// Jump to block 13
	ctx = MockBlockHeight(ctx, am, 13)

	// create mock vertices hashes and the corresponding obfuscated hash
	verticesHashes2 := []string{"666", "777", "888", "999", "1010"}
	obfuscatedHash2 := ObfuscatedDataHashHelper(verticesHashes2, providerId)

	_, err = ms.SubmitProgress(ctx, &types.MsgSubmitProgress{Provider: providerId, SubscriptionId: subscriptionId, PreviousVerticesHashes: verticesHashes1, ObfuscatedVerticesHash: obfuscatedHash2})

	// There should not be any error
	require.NoError(t, err)

}

func TestSubmitProgressWithIncorrectObfuscatedHash(t *testing.T) {
	_, ms, ctx, am := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)

	response, err := ms.CreateDeal(ctx, &types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 10000, StartBlock: 10, EndBlock: 20})
	require.NoError(t, err)

	dealId := response.DealId

	// Jump to block 12 to initiate the deal
	ctx = MockBlockHeight(ctx, am, 12)
	// Provider joins the deal after it is initiated
	joinDeal := types.MsgJoinDeal{Provider: Bob, DealId: dealId}
	joinResponse, err := ms.JoinDeal(ctx, &joinDeal)
	require.NoError(t, err)

	subscriptionId := joinResponse.SubscriptionId
	providerId := Bob

	// create mock vertices hashes and the corresponding obfuscated hash
	verticesHashes1 := []string{"000", "111", "222", "333", "444", "555"}
	// obfuscatedHash1 := MockObfuscatedDataHash(verticesHashes1, providerId)
	obfuscatedHash1 := "oogabooga"

	// submit progress
	_, err = ms.SubmitProgress(ctx, &types.MsgSubmitProgress{Provider: providerId, SubscriptionId: subscriptionId, ObfuscatedVerticesHash: obfuscatedHash1})
	// There should not be any error
	require.NoError(t, err)
	// Jump to block 13
	ctx = MockBlockHeight(ctx, am, 13)

	// create mock vertices hashes and the corresponding obfuscated hash
	verticesHashes2 := []string{"666", "777", "888", "999", "1010"}
	obfuscatedHash2 := ObfuscatedDataHashHelper(verticesHashes2, providerId)

	_, err = ms.SubmitProgress(ctx, &types.MsgSubmitProgress{Provider: providerId, SubscriptionId: subscriptionId, PreviousVerticesHashes: verticesHashes1, ObfuscatedVerticesHash: obfuscatedHash2})

	// There should be an error because you submitted the wrong vertices hashes
	require.Error(t, err)

}

func TestSubmitProgressAfterEpochDeadline(t *testing.T) {
	_, ms, ctx, am := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)

	response, err := ms.CreateDeal(ctx, &types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 10000, StartBlock: 10, EndBlock: 25})
	require.NoError(t, err)

	dealId := response.DealId

	// Jump to block 12 to initiate the deal
	ctx = MockBlockHeight(ctx, am, 12)
	// Provider joins the deal after it is initiated
	joinDeal := types.MsgJoinDeal{Provider: Bob, DealId: dealId}
	joinResponse, err := ms.JoinDeal(ctx, &joinDeal)
	require.NoError(t, err)

	subscriptionId := joinResponse.SubscriptionId
	providerId := Bob

	// create mock vertices hashes and the corresponding obfuscated hash
	verticesHashes1 := []string{"000", "111", "222", "333", "444", "555"}
	// obfuscatedHash1 := MockObfuscatedDataHash(verticesHashes1, providerId)
	obfuscatedHash1 := ObfuscatedDataHashHelper(verticesHashes1, providerId)

	// submit progress
	_, err = ms.SubmitProgress(ctx, &types.MsgSubmitProgress{Provider: providerId, SubscriptionId: subscriptionId, ObfuscatedVerticesHash: obfuscatedHash1})
	// There should not be any error
	require.NoError(t, err)
	// Jump to block 23
	ctx = MockBlockHeight(ctx, am, 23)

	// create mock vertices hashes and the corresponding obfuscated hash
	verticesHashes2 := []string{"666", "777", "888", "999", "1010"}
	obfuscatedHash2 := ObfuscatedDataHashHelper(verticesHashes2, providerId)

	_, err = ms.SubmitProgress(ctx, &types.MsgSubmitProgress{Provider: providerId, SubscriptionId: subscriptionId, PreviousVerticesHashes: verticesHashes1, ObfuscatedVerticesHash: obfuscatedHash2})

	// There should be an error beacuse you submit after the epoch deadline
	require.Error(t, err)
}

func TestClaimRewardsSingleProvider(t *testing.T) {
	k, ms, ctx, am := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)

	response, err := ms.CreateDeal(ctx, &types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 10000, StartBlock: 10, EndBlock: 20})
	require.NoError(t, err)

	dealId := response.DealId

	// Jump to block 12 to initiate the deal
	ctx = MockBlockHeight(ctx, am, 12)
	// Provider joins the deal after it is initiated
	joinDeal := types.MsgJoinDeal{Provider: Bob, DealId: dealId}
	joinResponse, err := ms.JoinDeal(ctx, &joinDeal)
	require.NoError(t, err)

	subscriptionId := joinResponse.SubscriptionId
	providerId := Bob

	// create mock vertices hashes and the corresponding obfuscated hash
	verticesHashes1 := []string{"000", "111", "222", "333", "444", "555"}
	obfuscatedHash1 := ObfuscatedDataHashHelper(verticesHashes1, providerId)

	// submit progress
	_, err = ms.SubmitProgress(ctx, &types.MsgSubmitProgress{Provider: providerId, SubscriptionId: subscriptionId, ObfuscatedVerticesHash: obfuscatedHash1})
	// There should not be any error
	require.NoError(t, err)
	// Jump to block 13
	ctx = MockBlockHeight(ctx, am, 13)

	// create mock vertices hashes and the corresponding obfuscated hash
	verticesHashes2 := []string{"666", "777", "888", "999", "1010"}
	obfuscatedHash2 := ObfuscatedDataHashHelper(verticesHashes2, providerId)

	_, err = ms.SubmitProgress(ctx, &types.MsgSubmitProgress{Provider: providerId, SubscriptionId: subscriptionId, PreviousVerticesHashes: verticesHashes1, ObfuscatedVerticesHash: obfuscatedHash2})

	// There should not be any error
	require.NoError(t, err)

	// Jump 100 (ChallengePeriod) blocks forward to elapse the challenge window.
	ctx = MockBlockHeight(ctx, am, 13+100)

	// deal amount available before claiming rewards
	deal, _ := k.GetDeal(ctx, dealId)
	dealAmountBefore := deal.AvailableAmount

	_, err = ms.ClaimRewards(ctx, &types.MsgClaimRewards{Provider: providerId, SubscriptionId: subscriptionId})
	require.NoError(t, err)

	deal, _ = k.GetDeal(ctx, dealId)
	dealAmountAfter := deal.AvailableAmount

	require.Greater(t, dealAmountBefore, dealAmountAfter)
}

func TestClaimRewardsDoubleProviders(t *testing.T) {
	k, ms, ctx, am := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)

	response, err := ms.CreateDeal(ctx, &types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 10000, StartBlock: 10, EndBlock: 20})
	require.NoError(t, err)

	dealId := response.DealId
	providerId1 := Bob
	providerId2 := Carol

	// Jump to block 12 to initiate the deal
	ctx = MockBlockHeight(ctx, am, 12)
	// Provider Bob joins the deal after it is initiated
	joinDeal := types.MsgJoinDeal{Provider: providerId1, DealId: dealId}
	joinResponse, err := ms.JoinDeal(ctx, &joinDeal)
	require.NoError(t, err)

	subscriptionId1 := joinResponse.SubscriptionId

	// Provider Carol joins the deal after it is initiated
	joinDeal = types.MsgJoinDeal{Provider: providerId2, DealId: dealId}
	joinResponse, err = ms.JoinDeal(ctx, &joinDeal)
	require.NoError(t, err)

	subscriptionId2 := joinResponse.SubscriptionId

	// create mock vertices hashes and the corresponding obfuscated hash
	verticesHashes1 := []string{"000", "111", "222", "333", "444", "555"}
	obfuscatedHash1 := ObfuscatedDataHashHelper(verticesHashes1, providerId1)
	obfuscatedHash2 := ObfuscatedDataHashHelper(verticesHashes1, providerId2)

	// submit progress provider1
	_, err = ms.SubmitProgress(ctx, &types.MsgSubmitProgress{Provider: providerId1, SubscriptionId: subscriptionId1, ObfuscatedVerticesHash: obfuscatedHash1})
	// There should not be any error
	require.NoError(t, err)

	// submit progress provider2
	_, err = ms.SubmitProgress(ctx, &types.MsgSubmitProgress{Provider: providerId2, SubscriptionId: subscriptionId2, ObfuscatedVerticesHash: obfuscatedHash2})
	// There should not be any error
	require.NoError(t, err)

	// Jump to block 13
	ctx = MockBlockHeight(ctx, am, 13)

	// create mock vertices hashes and the corresponding obfuscated hash
	verticesHashes2 := []string{"666", "777", "888", "999", "1010"}
	obfuscatedHash12 := ObfuscatedDataHashHelper(verticesHashes2, providerId1)
	obfuscatedHash22 := ObfuscatedDataHashHelper(verticesHashes2, providerId2)

	_, err = ms.SubmitProgress(ctx, &types.MsgSubmitProgress{Provider: providerId1, SubscriptionId: subscriptionId1, PreviousVerticesHashes: verticesHashes1, ObfuscatedVerticesHash: obfuscatedHash12})
	// There should not be any error
	require.NoError(t, err)

	_, err = ms.SubmitProgress(ctx, &types.MsgSubmitProgress{Provider: providerId2, SubscriptionId: subscriptionId2, PreviousVerticesHashes: verticesHashes1, ObfuscatedVerticesHash: obfuscatedHash22})
	// There should not be any error
	require.NoError(t, err)

	// Jump 100 (ChallengePeriod) blocks forward to elapse the challenge window.
	ctx = MockBlockHeight(ctx, am, 13+100)

	// deal amount available before claiming rewards
	deal, _ := k.GetDeal(ctx, dealId)
	dealAmountBeforeClaims := deal.AvailableAmount

	// Provider1 claims rewards
	_, err = ms.ClaimRewards(ctx, &types.MsgClaimRewards{Provider: providerId1, SubscriptionId: subscriptionId1})
	require.NoError(t, err)

	deal, _ = k.GetDeal(ctx, dealId)
	dealAmountAfter1Claims := deal.AvailableAmount

	require.Greater(t, dealAmountBeforeClaims, dealAmountAfter1Claims)

	// Provider1 claims rewards
	_, err = ms.ClaimRewards(ctx, &types.MsgClaimRewards{Provider: providerId2, SubscriptionId: subscriptionId2})
	require.NoError(t, err)

	deal, _ = k.GetDeal(ctx, dealId)
	dealAmountAfter2Claims := deal.AvailableAmount

	require.Greater(t, dealAmountAfter1Claims, dealAmountAfter2Claims)

	// the reward claimed by both providers must be equal since they submit the same vertices hashes
	require.Equal(t, dealAmountBeforeClaims-dealAmountAfter1Claims, dealAmountAfter1Claims-dealAmountAfter2Claims)
}
