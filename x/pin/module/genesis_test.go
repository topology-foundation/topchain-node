package pin_test

import (
	"testing"

	keepertest "topchain/testutil/keeper"
	"topchain/testutil/nullify"
	pin "topchain/x/pin/module"
	"topchain/x/pin/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		PinRequestList: []types.PinRequest{
			{
				Index: "0",
			},
			{
				Index: "1",
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.PinKeeper(t)
	pin.InitGenesis(ctx, k, genesisState)
	got := pin.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.PinRequestList, got.PinRequestList)
	// this line is used by starport scaffolding # genesis/test/assert
}
