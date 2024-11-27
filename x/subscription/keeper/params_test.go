package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"mandu/x/subscription/types"
)

func TestGetParams(t *testing.T) {
	k, ctx, _ := MockSubscriptionKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
