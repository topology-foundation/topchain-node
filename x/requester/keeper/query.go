package keeper

import (
	"topchain/x/requester/types"
)

var _ types.QueryServer = Keeper{}
