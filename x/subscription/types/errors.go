package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/subscription module sentinel errors
var (
	ErrInvalidSigner = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrSample        = sdkerrors.Register(ModuleName, 1101, "sample error")
	ErrTopologyRPC   = sdkerrors.Register(ModuleName, 1102, "topology rpc failure")
)
