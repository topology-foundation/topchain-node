package provider_test

import (
	"testing"

	keepertest "topchain/testutil/keeper"
	"topchain/testutil/nullify"
	provider "topchain/x/provider/module"
	"topchain/x/provider/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.ProviderKeeper(t)
	provider.InitGenesis(ctx, k, genesisState)
	got := provider.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
