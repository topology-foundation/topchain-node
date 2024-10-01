package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "topchain/testutil/keeper"
	"topchain/x/subscription/types"
)

func TestGetParams(t *testing.T) {
	k, ctx, _ := keepertest.SubscriptionKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
