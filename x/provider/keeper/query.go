package keeper

import (
	"topchain/x/provider/types"
)

var _ types.QueryServer = Keeper{}
