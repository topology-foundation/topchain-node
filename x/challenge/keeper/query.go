package keeper

import (
	"mandu/x/challenge/types"
)

var _ types.QueryServer = Keeper{}
