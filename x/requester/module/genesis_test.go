package requester_test

import (
	"testing"

	keepertest "topchain/testutil/keeper"
	"topchain/testutil/nullify"
	requester "topchain/x/requester/module"
	"topchain/x/requester/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.RequesterKeeper(t)
	requester.InitGenesis(ctx, k, genesisState)
	got := requester.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
