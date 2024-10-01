package subscription_test

import (
	"testing"

	keepertest "topchain/testutil/keeper"
	"topchain/testutil/nullify"
	subscription "topchain/x/subscription/module"
	"topchain/x/subscription/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx, _ := keepertest.SubscriptionKeeper(t)
	subscription.InitGenesis(ctx, k, genesisState)
	got := subscription.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
