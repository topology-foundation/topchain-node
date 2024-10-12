package keeper

import (
	"topchain/x/challenge/types"
)

var _ types.QueryServer = Keeper{}
