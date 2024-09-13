package keeper

import (
	"topchain/x/pin/types"
)

var _ types.QueryServer = Keeper{}
