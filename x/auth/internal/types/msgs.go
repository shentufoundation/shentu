package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	ModuleName = "auth"
	RouterKey  = ModuleName
)

// MsgManualVesting unlocks the specified amount in a manual vesting account.
type MsgManualVesting struct {
	Certifier    sdk.AccAddress `json:"certifier" yaml:"certifier"`
	Account      sdk.AccAddress `json:"account_address" yaml:"account_address"`
	UnlockAmount sdk.Coin       `json:"unlock_amount" yaml:"unlock_amount"`
}

var _ sdk.Msg = MsgManualVesting{}

// NewMsgManualVesting returns a MsgManualVesting object.
func NewMsgManualVesting(certifier, account sdk.AccAddress, unlockAmount sdk.Coin) MsgManualVesting {
	return MsgManualVesting{
		Certifier:    certifier,
		Account:      account,
		UnlockAmount: unlockAmount,
	}
}

// Route returns the name of the module.
func (m MsgManualVesting) Route() string { return ModuleName }

// Type returns a human-readable string for the message.
func (m MsgManualVesting) Type() string { return "manual_vesting" }

// ValidateBasic runs stateless checks on the message.
func (m MsgManualVesting) ValidateBasic() error {
	if m.Certifier.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing from address")
	}
	if m.Account.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing account address")
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgManualVesting) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required.
func (m MsgManualVesting) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Certifier}
}

// MsgLockedSend transfers coins and have them vesting
// in the receiver's manual vesting account.
type MsgLockedSend struct {
	From   sdk.AccAddress `json:"from" yaml:"from"`
	To     sdk.AccAddress `json:"to" yaml:"to"`
	Amount sdk.Coin       `json:"amount" yaml:"amount"`
}

var _ sdk.Msg = MsgLockedSend{}

// NewMsgLockedSend returns a MsgLockedSend object.
func NewMsgLockedSend(from, to sdk.AccAddress, amount sdk.Coin) MsgLockedSend {
	return MsgLockedSend{
		From:   from,
		To:     to,
		Amount: amount,
	}
}

// Route returns the name of the module.
func (m MsgLockedSend) Route() string { return ModuleName }

// Type returns a human-readable string for the message.
func (m MsgLockedSend) Type() string { return "locked_send" }

// ValidateBasic runs stateless checks on the message.
func (m MsgLockedSend) ValidateBasic() error {
	if m.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing from address")
	}
	if m.To.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing to address")
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgLockedSend) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required.
func (m MsgLockedSend) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.From}
}
