package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgRequestSubscription{}

func NewMsgRequestSubscription(creator string, croid string, ammount string, duration string) *MsgRequestSubscription {
	return &MsgRequestSubscription{
		Creator:  creator,
		Croid:    croid,
		Ammount:  ammount,
		Duration: duration,
	}
}

func (msg *MsgRequestSubscription) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
