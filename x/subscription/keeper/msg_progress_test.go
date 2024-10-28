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

	_, err = ms.SubmitProgress(ctx, &types.MsgSubmitProgress{Provider: providerId, SubscriptionId: subscriptionId, VerticesHashes: verticesHashes1, ObfuscatedVerticesHash: obfuscatedHash2})

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

	_, err = ms.SubmitProgress(ctx, &types.MsgSubmitProgress{Provider: providerId, SubscriptionId: subscriptionId, VerticesHashes: verticesHashes1, ObfuscatedVerticesHash: obfuscatedHash2})

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

	_, err = ms.SubmitProgress(ctx, &types.MsgSubmitProgress{Provider: providerId, SubscriptionId: subscriptionId, VerticesHashes: verticesHashes1, ObfuscatedVerticesHash: obfuscatedHash2})

	// There should be an error beacuse you submit after the epoch deadline
	require.Error(t, err)

}

func TestCheckPayoutAfterSubmitProgress(t *testing.T) {
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
	obfuscatedHash1 := ObfuscatedDataHashHelper(verticesHashes1, providerId)

	// submit obfuscated hash
	_, err = ms.SubmitProgress(ctx, &types.MsgSubmitProgress{Provider: providerId, SubscriptionId: subscriptionId, ObfuscatedVerticesHash: obfuscatedHash1})
	// There should not be any error
	require.NoError(t, err)
	// Jump to block 15
	ctx = MockBlockHeight(ctx, am, 15)

	// create mock vertices hashes and the corresponding obfuscated hash
	verticesHashes2 := []string{"666", "777", "888", "999", "1010"}
	obfuscatedHash2 := ObfuscatedDataHashHelper(verticesHashes2, providerId)

	// Get the "top" coin balance for Bob before submitting the real vertices hashes
	coin, err := CheckBankBalance(ctx, Bob)
	require.NoError(t, err)
	topBalanceBefore := coin.Amount.Uint64()

	// submit real vertices hash
	_, err = ms.SubmitProgress(ctx, &types.MsgSubmitProgress{Provider: providerId, SubscriptionId: subscriptionId, VerticesHashes: verticesHashes1, ObfuscatedVerticesHash: obfuscatedHash2})

	require.NoError(t, err)

	// progress, _ := k.GetProgressSize(ctx, subscriptionId, 15)
	// fmt.Println("Bob's progress", progress)

	// Get the "top" coin balance for Bob before submitting the real vertices hashes

	ctx = MockBlockHeight(ctx, am, 15)

	coin, err = CheckBankBalance(ctx, Bob)
	require.NoError(t, err)
	topBalanceAfter := coin.Amount.Uint64()

	require.Greater(t, topBalanceAfter, topBalanceBefore)
}
