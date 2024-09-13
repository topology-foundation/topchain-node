package keeper_test

import (
	"context"
	"strconv"
	"testing"

	keepertest "topchain/testutil/keeper"
	"topchain/testutil/nullify"
	"topchain/x/pin/keeper"
	"topchain/x/pin/types"

	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNPinRequest(keeper keeper.Keeper, ctx context.Context, n int) []types.PinRequest {
	items := make([]types.PinRequest, n)
	for i := range items {
		items[i].Index = strconv.Itoa(i)

		keeper.SetPinRequest(ctx, items[i])
	}
	return items
}

func TestPinRequestGet(t *testing.T) {
	keeper, ctx := keepertest.PinKeeper(t)
	items := createNPinRequest(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetPinRequest(ctx,
			item.Index,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestPinRequestRemove(t *testing.T) {
	keeper, ctx := keepertest.PinKeeper(t)
	items := createNPinRequest(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemovePinRequest(ctx,
			item.Index,
		)
		_, found := keeper.GetPinRequest(ctx,
			item.Index,
		)
		require.False(t, found)
	}
}

func TestPinRequestGetAll(t *testing.T) {
	keeper, ctx := keepertest.PinKeeper(t)
	items := createNPinRequest(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllPinRequest(ctx)),
	)
}
