package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgRequestSubscription{}

func NewMsgRequestSubscription(creator string, croid string, amount uint64, duration int32) *MsgRequestSubscription {
	return &MsgRequestSubscription{
		Creator:  creator,
		CroId:    croid,
		Amount:   amount,
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
