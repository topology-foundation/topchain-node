package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCancelSubscription{}

func NewMsgCancelSubscription(creator string, subscriptionid string) *MsgCancelSubscription {
	return &MsgCancelSubscription{
		Creator:        creator,
		Subscriptionid: subscriptionid,
	}
}

func (msg *MsgCancelSubscription) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
