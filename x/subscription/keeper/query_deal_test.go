package keeper_test

import (
	"testing"

	testutil "topchain/testutil/keeper"

	"topchain/x/subscription/types"

	qtypes "github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
)

func TestQueryDeal(t *testing.T) {
	k, ms, ctx, _ := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	// Create a deal
	createDeal := types.MsgCreateDeal{Requester: testutil.Alice, CroId: "alicecro", Amount: 1000, StartBlock: 10, EndBlock: 20}
	createResponse, err := ms.CreateDeal(ctx, &createDeal)
	require.Nil(t, err)

	// Query the deal
	queryResponse, err := k.Deal(ctx, &types.QueryDealRequest{Id: createResponse.DealId})
	require.Nil(t, err)

	require.EqualValues(t, createDeal, types.MsgCreateDeal{Requester: queryResponse.Deal.Requester, CroId: queryResponse.Deal.CroId, Amount: queryResponse.Deal.TotalAmount, StartBlock: queryResponse.Deal.StartBlock, EndBlock: queryResponse.Deal.EndBlock})
}

func TestQueryDealStatus(t *testing.T) {
	k, ms, ctx, _ := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	// Create a deal
	createDeal := types.MsgCreateDeal{Requester: testutil.Alice, CroId: "alicecro", Amount: 1000, StartBlock: 10, EndBlock: 20}
	createResponse, err := ms.CreateDeal(ctx, &createDeal)
	require.Nil(t, err)

	// Query the deal
	queryStatusResponse, err := k.DealStatus(ctx, &types.QueryDealStatusRequest{Id: createResponse.DealId})
	require.Nil(t, err)

	require.EqualValues(t, queryStatusResponse.Status, types.Deal_SCHEDULED)

}

func TestQueryDeals(t *testing.T) {
	k, ms, ctx, _ := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	// Create a first deal
	createDeal1 := types.MsgCreateDeal{Requester: testutil.Alice, CroId: "alicecro1", Amount: 1000, StartBlock: 10, EndBlock: 20}
	createResponse1, err := ms.CreateDeal(ctx, &createDeal1)
	require.Nil(t, err)

	deal1, found := k.GetDeal(ctx, createResponse1.DealId)
	require.True(t, found)

	// Create a second deal
	createDeal2 := types.MsgCreateDeal{Requester: testutil.Alice, CroId: "alicecro2", Amount: 1000, StartBlock: 20, EndBlock: 30}
	createResponse2, err := ms.CreateDeal(ctx, &createDeal2)
	require.Nil(t, err)

	deal2, found := k.GetDeal(ctx, createResponse2.DealId)
	require.True(t, found)

	// Create a third deal
	createDeal3 := types.MsgCreateDeal{Requester: testutil.Alice, CroId: "alicecro", Amount: 1000, StartBlock: 10, EndBlock: 20}
	createResponse3, err := ms.CreateDeal(ctx, &createDeal3)
	require.Nil(t, err)

	deal3, found := k.GetDeal(ctx, createResponse3.DealId)
	require.True(t, found)

	// Query for all the deals by alice
	queryDealsResponse, err := k.Deals(ctx, &types.QueryDealsRequest{Requester: testutil.Alice})
	require.Nil(t, err)

	require.Contains(t, queryDealsResponse.Deals, deal1)
	require.Contains(t, queryDealsResponse.Deals, deal2)
	require.Contains(t, queryDealsResponse.Deals, deal3)
}

func TestQueryDealsWithPagination(t *testing.T) {
	k, ms, ctx, _ := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	// Create a first deal
	createDeal1 := types.MsgCreateDeal{Requester: testutil.Alice, CroId: "alicecro1", Amount: 1000, StartBlock: 10, EndBlock: 20}
	createResponse1, err := ms.CreateDeal(ctx, &createDeal1)
	require.Nil(t, err)

	deal1, found := k.GetDeal(ctx, createResponse1.DealId)
	require.True(t, found)

	// Create a second deal
	createDeal2 := types.MsgCreateDeal{Requester: testutil.Alice, CroId: "alicecro2", Amount: 1000, StartBlock: 20, EndBlock: 30}
	createResponse2, err := ms.CreateDeal(ctx, &createDeal2)
	require.Nil(t, err)

	deal2, found := k.GetDeal(ctx, createResponse2.DealId)
	require.True(t, found)

	// Create a third deal
	createDeal3 := types.MsgCreateDeal{Requester: testutil.Alice, CroId: "alicecro", Amount: 1000, StartBlock: 10, EndBlock: 20}
	createResponse3, err := ms.CreateDeal(ctx, &createDeal3)
	require.Nil(t, err)

	deal3, found := k.GetDeal(ctx, createResponse3.DealId)
	require.True(t, found)

	// Query for all the deals by alice
	queryDealsResponse, err := k.Deals(ctx, &types.QueryDealsRequest{Requester: testutil.Alice, Pagination: &qtypes.PageRequest{Limit: 2}})
	require.Nil(t, err)

	require.Equal(t, len(queryDealsResponse.Deals), 2)
	require.Contains(t, []types.Deal{deal1, deal2, deal3}, queryDealsResponse.Deals[0])
	require.Contains(t, []types.Deal{deal1, deal2, deal3}, queryDealsResponse.Deals[1])
}
