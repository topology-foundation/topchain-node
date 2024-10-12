package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) PricePerVertexChallenge(ctx sdk.Context, challenger_address string, provider_id string) int64 {
	// TODO - devise a formula, set of arguments might be different
	return 1
}
