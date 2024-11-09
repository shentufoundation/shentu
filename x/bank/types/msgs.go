package types

import (
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// TypeMsgLockedSend bank message types
const (
	TypeMsgLockedSend = "locked_send"
)

var _ sdk.Msg = &MsgLockedSend{}

// NewMsgLockedSend returns a MsgLockedSend object.
func NewMsgLockedSend(from, to sdk.AccAddress, unlocker string, amount sdk.Coins) *MsgLockedSend {
	return &MsgLockedSend{
		FromAddress:     from.String(),
		ToAddress:       to.String(),
		UnlockerAddress: unlocker,
		Amount:          amount,
	}
}

// Route returns the name of the module.
func (m MsgLockedSend) Route() string { return bankTypes.RouterKey }

// Type returns a human-readable string for the message.
func (m MsgLockedSend) Type() string { return TypeMsgLockedSend }

// ValidateBasic runs stateless checks on the message.
func (m MsgLockedSend) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(m.ToAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid recipient address (%s)", err)
	}

	return nil
}

// GetSigners defines whose signature is required.
func (m MsgLockedSend) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}
