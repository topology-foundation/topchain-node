package keeper_test

import (
	"testing"

	"mandu/x/subscription/types"

	"github.com/stretchr/testify/require"
)

func TestMsgServerCreateDealMsg(t *testing.T) {
	k, ms, ctx, _ := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	response, err := ms.CreateDeal(ctx, &types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 10000, StartBlock: 10, EndBlock: 20})
	require.NoError(t, err)

	require.NotEmpty(t, response)
	require.NotEmpty(t, response.DealId)
}

func TestMsgServerCreateDealScheduled(t *testing.T) {
	k, ms, ctx, _ := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	createDeal := types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 10000, StartBlock: 10, EndBlock: 20}

	response, err := ms.CreateDeal(ctx, &createDeal)
	require.NoError(t, err)

	deal, found := k.GetDeal(ctx, response.DealId)
	require.True(t, found)

	require.EqualValues(t, response.DealId, deal.Id)

	require.EqualValues(t, createDeal, types.MsgCreateDeal{Requester: deal.Requester, CroId: deal.CroId, Amount: deal.TotalAmount, StartBlock: deal.StartBlock, EndBlock: deal.EndBlock})

	require.Equal(t, deal.Status, types.Deal_SCHEDULED)
}

func TestMsgServerCreateDealInitializedStatus(t *testing.T) {
	k, ms, ctx, am := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	createDeal := types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 10000, StartBlock: 10, EndBlock: 20}
	response, err := ms.CreateDeal(ctx, &createDeal)
	require.NoError(t, err)

	// Get the deal from the storage
	deal, found := k.GetDeal(ctx, response.DealId)
	require.True(t, found)

	// Jump to block number 11
	ctx = MockBlockHeight(ctx, am, 10)

	// The deal must be initialized after entering block 10
	deal, _ = k.GetDeal(ctx, response.DealId)

	require.Equal(t, deal.Status, types.Deal_INITIALIZED)
}

func TestMsgServerCancelDealCorrectRequester(t *testing.T) {
	k, ms, ctx, _ := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	createDeal := types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 10000, StartBlock: 10, EndBlock: 20}
	createResponse, err := ms.CreateDeal(ctx, &createDeal)
	require.NoError(t, err)

	// Get the deal from the storage
	_, found := k.GetDeal(ctx, createResponse.DealId)
	require.True(t, found)

	// Now send a cancel message
	cancelDeal := types.MsgCancelDeal{Requester: Alice, DealId: createResponse.DealId}
	_, err = ms.CancelDeal(ctx, &cancelDeal)
	require.NoError(t, err)

	// Get the deal from the storage
	deal, found := k.GetDeal(ctx, createResponse.DealId)
	if !found {
		t.Fatalf("Deal not found")
	}
	require.EqualValues(t, deal.Status, types.Deal_CANCELLED)
}

func TestMsgServerCancelDealIncorrectRequester(t *testing.T) {
	k, ms, ctx, _ := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	createDeal := types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 10000, StartBlock: 10, EndBlock: 20}
	createResponse, err := ms.CreateDeal(ctx, &createDeal)
	require.NoError(t, err)

	// Get the deal from the storage
	_, found := k.GetDeal(ctx, createResponse.DealId)
	require.True(t, found)

	// Now send a cancel message
	cancelDeal := types.MsgCancelDeal{Requester: Bob, DealId: createResponse.DealId}
	_, err = ms.CancelDeal(ctx, &cancelDeal)

	// The error should not be nil because the incorrect requester sends the CancelDeal message
	require.NotNil(t, err)
}

func TestMsgServerUpdateDealIncorrectRequesterMsg(t *testing.T) {}
func TestMsgServerUpdateScheduledDealCorrectStartBlockMsg(t *testing.T) {
	k, ms, ctx, _ := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	createResponse, err := ms.CreateDeal(ctx, &types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 10000, StartBlock: 10, EndBlock: 20})
	require.NoError(t, err)

	updateDeal := types.MsgUpdateDeal{Requester: Alice, DealId: createResponse.DealId, StartBlock: 11}
	_, err = ms.UpdateDeal(ctx, &updateDeal)
	require.NoError(t, err)

	deal, found := k.GetDeal(ctx, createResponse.DealId)

	require.True(t, found)
	require.EqualValues(t, deal.StartBlock, 11)
}

func TestMsgServerUpdateScheduledDealIncorrectStartBlockMsg(t *testing.T) {
	k, ms, ctx, am := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	createResponse, err := ms.CreateDeal(ctx, &types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 10000, StartBlock: 10, EndBlock: 20})
	require.NoError(t, err)

	// Jump to block 9
	ctx = MockBlockHeight(ctx, am, 9)

	updateDeal := types.MsgUpdateDeal{Requester: Alice, DealId: createResponse.DealId, StartBlock: 7}
	_, err = ms.UpdateDeal(ctx, &updateDeal)

	require.NotNil(t, err)
}

func TestMsgServerUpdateInitiatedDealIncorrectStartBlockMsg(t *testing.T) {
	k, ms, ctx, am := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	createResponse, err := ms.CreateDeal(ctx, &types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 10000, StartBlock: 10, EndBlock: 20})
	require.NoError(t, err)

	// Jump to block 12
	ctx = MockBlockHeight(ctx, am, 12)

	updateDeal := types.MsgUpdateDeal{Requester: Alice, DealId: createResponse.DealId, StartBlock: 7}
	_, err = ms.UpdateDeal(ctx, &updateDeal)

	// It should return an error because the StartBlock can't be updated once the deal is initiated.
	require.NotNil(t, err)
}

func TestMsgServerUpdateScheduledDealIncrementAmountMsg(t *testing.T) {
	k, ms, ctx, _ := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	createResponse, err := ms.CreateDeal(ctx, &types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 10000, StartBlock: 10, EndBlock: 20})
	require.NoError(t, err)

	updateDeal := types.MsgUpdateDeal{Requester: Alice, DealId: createResponse.DealId, Amount: 12000}
	_, err = ms.UpdateDeal(ctx, &updateDeal)
	require.NoError(t, err)

	deal, found := k.GetDeal(ctx, createResponse.DealId)
	require.True(t, found)
	require.EqualValues(t, deal.TotalAmount, 12000)
}

func TestMsgServerUpdateScheduledDealDecrementAmountMsg(t *testing.T) {
	k, ms, ctx, _ := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	createResponse, err := ms.CreateDeal(ctx, &types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 10000, StartBlock: 10, EndBlock: 20})
	require.NoError(t, err)

	_, err = ms.UpdateDeal(ctx, &types.MsgUpdateDeal{Requester: Alice, DealId: createResponse.DealId, Amount: 5000})
	require.NoError(t, err)

	deal, found := k.GetDeal(ctx, createResponse.DealId)
	require.True(t, found)
	require.EqualValues(t, deal.TotalAmount, 5000)
}

func TestMsgServerUpdateScheduledDealDecrementTotalAmountMsg(t *testing.T) {
	k, ms, ctx, _ := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	createResponse, err := ms.CreateDeal(ctx, &types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 10000, StartBlock: 10, EndBlock: 20})
	require.NoError(t, err)

	_, err = ms.UpdateDeal(ctx, &types.MsgUpdateDeal{Requester: Alice, DealId: createResponse.DealId, Amount: 0})
	require.NoError(t, err)

	deal, found := k.GetDeal(ctx, createResponse.DealId)
	require.True(t, found)
	// Amount should be unchanged because you cannot withdraw full amount while the deal is still active.
	require.EqualValues(t, deal.TotalAmount, 10000)
}

func TestMsgServerUpdateInitiatedDealIncrementAmountMsg(t *testing.T) {
	k, ms, ctx, am := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	createResponse, err := ms.CreateDeal(ctx, &types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 10000, StartBlock: 10, EndBlock: 20})
	require.NoError(t, err)

	// Jump to block 12 to initiate the deal
	ctx = MockBlockHeight(ctx, am, 12)

	_, err = ms.UpdateDeal(ctx, &types.MsgUpdateDeal{Requester: Alice, DealId: createResponse.DealId, Amount: 15000})
	require.NoError(t, err)

	deal, found := k.GetDeal(ctx, createResponse.DealId)
	require.True(t, found)
	require.EqualValues(t, deal.TotalAmount, 15000)
}

func TestMsgServerUpdateInitiatedDealDecrementAmountMsg(t *testing.T) {
	k, ms, ctx, am := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	createResponse, err := ms.CreateDeal(ctx, &types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 10000, StartBlock: 10, EndBlock: 20})
	require.NoError(t, err)

	// Jump to block 12 to initiate the deal
	ctx = MockBlockHeight(ctx, am, 12)

	_, err = ms.UpdateDeal(ctx, &types.MsgUpdateDeal{Requester: Alice, DealId: createResponse.DealId, Amount: 9000})
	// It should return an error because you're not allowed to decrease the amount after deal initiation
	require.NotNil(t, err)
}

func TestMsgServerJoinDealBeforeInitiationMsg(t *testing.T) {
	k, ms, ctx, _ := setupMsgServer(t)

	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	// Create a new deal
	createDeal := types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 1000, StartBlock: 10, EndBlock: 20}
	createResponse, err := ms.CreateDeal(ctx, &createDeal)
	require.NoError(t, err)

	dealId := createResponse.DealId

	// Get the deal from the storage
	deal, found := k.GetDeal(ctx, dealId)

	require.True(t, found)
	// Assert the status of the deal to be "SCHEDULED"
	require.EqualValues(t, deal.Status, types.Deal_SCHEDULED)

	// Provider joins the deal before it is initiated
	joinDeal := types.MsgJoinDeal{Provider: Bob, DealId: dealId}
	joinResponse, err := ms.JoinDeal(ctx, &joinDeal)
	require.NoError(t, err)

	// Check if the subscription exists
	sub, found := k.GetSubscription(ctx, joinResponse.SubscriptionId)

	require.True(t, found)
	require.EqualValues(t, sub.Provider, Bob)

	// Check if the subscription exists in the deal's subscriptionIds
	deal, _ = k.GetDeal(ctx, dealId)

	// Assert that the last id in deal's subscriptionIds' is sub's id
	require.EqualValues(t, deal.SubscriptionIds[len(deal.SubscriptionIds)-1], sub.Id)
	require.EqualValues(t, dealId, sub.DealId)
}

func TestMsgServerJoinInitiatedDealMsg(t *testing.T) {
	k, ms, ctx, am := setupMsgServer(t)

	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	// Create a new deal
	createDeal := types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 1000, StartBlock: 10, EndBlock: 20}
	createResponse, err := ms.CreateDeal(ctx, &createDeal)
	require.NoError(t, err)

	dealId := createResponse.DealId

	// Get the deal from the storage
	deal, found := k.GetDeal(ctx, dealId)

	require.True(t, found)
	// Assert the status of the deal to be "SCHEDULED"
	require.EqualValues(t, deal.Status, types.Deal_SCHEDULED)

	// Jump to block 12
	ctx = MockBlockHeight(ctx, am, 12)
	// Provider joins the deal after it is initiated
	joinDeal := types.MsgJoinDeal{Provider: Bob, DealId: dealId}
	joinResponse, err := ms.JoinDeal(ctx, &joinDeal)
	require.NoError(t, err)

	// Check if the subscription exists
	sub, found := k.GetSubscription(ctx, joinResponse.SubscriptionId)

	require.True(t, found)
	require.EqualValues(t, sub.Provider, Bob)

	// Check if the subscription exists in the deal's subscriptionIds
	deal, _ = k.GetDeal(ctx, dealId)

	// Assert that the last id in deal's subscriptionIds' is sub's id
	require.EqualValues(t, deal.SubscriptionIds[len(deal.SubscriptionIds)-1], sub.Id)
	require.EqualValues(t, dealId, sub.DealId)

	// Check if the deal's status has changed to ACTIVE
	require.EqualValues(t, deal.Status, types.Deal_ACTIVE)
}

func TestMsgServerJoinCancelledDealMsg(t *testing.T) {
	k, ms, ctx, _ := setupMsgServer(t)

	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	// Create a new deal
	createDeal := types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 1000, StartBlock: 10, EndBlock: 20}
	createResponse, err := ms.CreateDeal(ctx, &createDeal)
	require.NoError(t, err)

	dealId := createResponse.DealId

	// Get the deal from the storage
	_, found := k.GetDeal(ctx, dealId)
	require.True(t, found)

	// Cancel the deal
	cancelDeal := types.MsgCancelDeal{Requester: Alice, DealId: dealId}
	_, err = ms.CancelDeal(ctx, &cancelDeal)
	require.NoError(t, err)

	// Check if the status is changed to CANCELLED
	deal, _ := k.GetDeal(ctx, dealId)
	require.EqualValues(t, deal.Status, types.Deal_CANCELLED)

	// Provider joins the deal before it is initiated
	joinDeal := types.MsgJoinDeal{Provider: Bob, DealId: dealId}
	_, err = ms.JoinDeal(ctx, &joinDeal)

	require.NotNil(t, err)
}

func TestMsgServerJoinSameDealMoreThanOnceMsg(t *testing.T) {
	k, ms, ctx, _ := setupMsgServer(t)

	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	// Create a new deal
	createDeal := types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 1000, StartBlock: 10, EndBlock: 20}
	createResponse, err := ms.CreateDeal(ctx, &createDeal)
	require.NoError(t, err)

	dealId := createResponse.DealId

	// Provider joins the deal
	joinDeal := types.MsgJoinDeal{Provider: Bob, DealId: dealId}
	_, err = ms.JoinDeal(ctx, &joinDeal)
	require.NoError(t, err)

	// Provider tries to join the same deal again
	_, err = ms.JoinDeal(ctx, &joinDeal)

	// It is disallowed to join a deal already subscribed to
	require.NotNil(t, err)
}

func TestMsgServerIncrementDealAmount(t *testing.T) {
	k, ms, ctx, _ := setupMsgServer(t)

	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	// Create a new deal
	createDeal := types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 1000, StartBlock: 10, EndBlock: 20}
	createResponse, err := ms.CreateDeal(ctx, &createDeal)
	require.NoError(t, err)

	dealId := createResponse.DealId

	// topup the deal amount
	incrementDeal := types.MsgIncrementDealAmount{Requester: Alice, DealId: dealId, Amount: 1000}
	_, err = ms.IncrementDealAmount(ctx, &incrementDeal)

	require.NoError(t, err)
}

func TestMsgServerIncrementDealAmountIncorrectRequester(t *testing.T) {
	k, ms, ctx, _ := setupMsgServer(t)

	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	// Create a new deal
	createDeal := types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 1000, StartBlock: 10, EndBlock: 20}
	createResponse, err := ms.CreateDeal(ctx, &createDeal)
	require.NoError(t, err)

	dealId := createResponse.DealId

	// topup the deal amount
	incrementDeal := types.MsgIncrementDealAmount{Requester: Bob, DealId: dealId, Amount: 1000}
	_, err = ms.IncrementDealAmount(ctx, &incrementDeal)

	require.NotNil(t, err)
}

func TestMsgServerIncrementCancelledDealAmount(t *testing.T) {
	k, ms, ctx, _ := setupMsgServer(t)

	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	// Create a new deal
	createDeal := types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 1000, StartBlock: 10, EndBlock: 20}
	createResponse, err := ms.CreateDeal(ctx, &createDeal)
	require.NoError(t, err)

	dealId := createResponse.DealId

	// cancel the deal
	cancelDeal := types.MsgCancelDeal{Requester: Alice, DealId: dealId}
	_, err = ms.CancelDeal(ctx, &cancelDeal)
	require.NoError(t, err)

	// topup the deal amount
	incrementDeal := types.MsgIncrementDealAmount{Requester: Alice, DealId: dealId, Amount: 1000}
	_, err = ms.IncrementDealAmount(ctx, &incrementDeal)

	require.NotNil(t, err)
}

func TestMsgServerIncrementExpiredDealAmount(t *testing.T) {
	k, ms, ctx, am := setupMsgServer(t)

	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	// Create a new deal
	createDeal := types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 1000, StartBlock: 10, EndBlock: 20}
	createResponse, err := ms.CreateDeal(ctx, &createDeal)
	require.NoError(t, err)

	dealId := createResponse.DealId

	// Jump to block 12 to initialize the deal
	ctx = MockBlockHeight(ctx, am, 12)
	// Jump to block 25 to expire the deal
	ctx = MockBlockHeight(ctx, am, 25)

	// topup the deal amount
	incrementDeal := types.MsgIncrementDealAmount{Requester: Alice, DealId: dealId, Amount: 1000}
	_, err = ms.IncrementDealAmount(ctx, &incrementDeal)

	require.NotNil(t, err)
}

func TestMsgServerLeaveJoinedDealMsg(t *testing.T) {
	k, ms, ctx, _ := setupMsgServer(t)

	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	// Create a new deal
	createDeal := types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 1000, StartBlock: 10, EndBlock: 20}
	createResponse, err := ms.CreateDeal(ctx, &createDeal)
	require.NoError(t, err)

	dealId := createResponse.DealId

	// Provider joins the deal
	joinDeal := types.MsgJoinDeal{Provider: Bob, DealId: dealId}
	_, err = ms.JoinDeal(ctx, &joinDeal)

	require.NoError(t, err)

	leaveDeal := types.MsgLeaveDeal{Provider: Bob, DealId: dealId}
	// Provider tries to leave the deal
	_, err = ms.LeaveDeal(ctx, &leaveDeal)

	require.NoError(t, err)
}

func TestMsgServerLeaveNotJoinedDealMsg(t *testing.T) {
	k, ms, ctx, _ := setupMsgServer(t)

	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	// Create a new deal
	createDeal := types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 1000, StartBlock: 10, EndBlock: 20}
	createResponse, err := ms.CreateDeal(ctx, &createDeal)
	require.NoError(t, err)

	dealId := createResponse.DealId

	leaveDeal := types.MsgLeaveDeal{Provider: Bob, DealId: dealId}
	// Provider tries to leave the deal it has not joined
	_, err = ms.LeaveDeal(ctx, &leaveDeal)

	// It should error because you can't leave a deal you did not join
	require.NotNil(t, err)
}

func TestMsgServerJoinLeaveJoinDeallMsg(t *testing.T) {
	k, ms, ctx, am := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	// Create a new deal
	createDeal := types.MsgCreateDeal{Requester: Alice, CroId: "alicecro", Amount: 1000, StartBlock: 10, EndBlock: 20}
	createResponse, err := ms.CreateDeal(ctx, &createDeal)
	require.NoError(t, err)

	dealId := createResponse.DealId

	// Provider joins the deal
	joinDeal := types.MsgJoinDeal{Provider: Bob, DealId: dealId}
	_, err = ms.JoinDeal(ctx, &joinDeal)
	require.NoError(t, err)

	leaveDeal := types.MsgLeaveDeal{Provider: Bob, DealId: dealId}
	// Provider tries to leave the deal it has not joined
	_, err = ms.LeaveDeal(ctx, &leaveDeal)
	require.NoError(t, err)

	// Jump one block forward
	ctx = MockBlockHeight(ctx, am, 1)
	// Provider joins the deal again
	_, err = ms.JoinDeal(ctx, &joinDeal)

	require.NoError(t, err)
}
