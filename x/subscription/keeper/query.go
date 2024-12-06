package keeper

import (
	"mandu/x/subscription/types"
)

var _ types.QueryServer = Keeper{}
