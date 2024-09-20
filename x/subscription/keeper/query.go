package keeper

import (
	"topchain/x/subscription/types"
)

var _ types.QueryServer = Keeper{}
