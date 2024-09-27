package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	query "github.com/cosmos/cosmos-sdk/types/query"

	keepertest "topchain/testutil/keeper"
	"topchain/x/subscription/types"
)

func TestDeal(t *testing.T) {
	keeper, ctx := keepertest.SubscriptionKeeper(t)
	deal := types.Deal{
		Id:        "deal1",
		Requester: "Requester1",
	}

	keeper.SetDeal(ctx, deal)

	req := &types.QueryDealRequest{Id: "deal1"}

	retrievedDeal, err := keeper.Deal(ctx, req)
	require.NoError(t, err)
	require.Equal(t, deal, retrievedDeal.Deal)
}

func TestDeals(t *testing.T) {
	keeper, ctx := keepertest.SubscriptionKeeper(t)
	deal1 := types.Deal{
		Id:        "deal1",
		Requester: "Requester1",
	}
	deal2 := types.Deal{
		Id:        "deal2",
		Requester: "Requester1",
	}
	deal3 := types.Deal{
		Id:        "deal3",
		Requester: "Requester2",
	}

	keeper.SetDeal(ctx, deal1)
	keeper.SetDeal(ctx, deal2)
	keeper.SetDeal(ctx, deal3)

	req := &types.QueryDealsRequest{Requester: "Requester1"}
	res, err := keeper.Deals(ctx, req)
	require.NoError(t, err)
	require.Len(t, res.Deals, 2)
	require.Contains(t, res.Deals, deal1)
	require.Contains(t, res.Deals, deal2)
}

func TestDealsWithPaginationOne(t *testing.T) {
	keeper, ctx := keepertest.SubscriptionKeeper(t)
	deal1 := types.Deal{
		Id:        "deal1",
		Requester: "Requester1",
	}
	deal2 := types.Deal{
		Id:        "deal2",
		Requester: "Requester1",
	}
	deal3 := types.Deal{
		Id:        "deal3",
		Requester: "Requester2",
	}

	keeper.SetDeal(ctx, deal1)
	keeper.SetDeal(ctx, deal2)
	keeper.SetDeal(ctx, deal3)

	req := &types.QueryDealsRequest{Requester: "Requester1", Pagination: &query.PageRequest{Limit: 1}}
	res, err := keeper.Deals(ctx, req)
	require.NoError(t, err)
	require.Len(t, res.Deals, 1)
	require.Contains(t, res.Deals, deal1)
	req = &types.QueryDealsRequest{Requester: "Requester1", Pagination: &query.PageRequest{Key: res.Pagination.NextKey, Limit: 1}}
	res, err = keeper.Deals(ctx, req)
	require.Len(t, res.Deals, 1)
	require.Contains(t, res.Deals, deal2)
}
