package validation

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func ValidateNonEmptyString(value string) error {
	if len(strings.TrimSpace(value)) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "string is empty")
	}
	return nil
}

func ValidatePositiveAmount(amount uint64) error {
	if amount <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount must be greater than 0")
	}
	return nil
}

func ValidateAddress(address string) error {
	if err := ValidateNonEmptyString(address); err != nil {
		return err
	}
	if _, err := sdk.AccAddressFromBech32(address); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid address")
	}
	return nil
}

func ValidateBlockRange(startBlock, endBlock uint64) error {
	if startBlock >= endBlock {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "start block must be less than end block")
	}
	return nil
}
