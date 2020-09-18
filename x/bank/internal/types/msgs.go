package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// RouterKey is they name of the bank module
const RouterKey = bank.ModuleName

// MsgSendLock transfers coins and have them vesting
// in the receiver's manual vesting account.
type MsgSendLock struct {
	From   sdk.AccAddress `json:"from" yaml:"from"`
	To     sdk.AccAddress `json:"to" yaml:"to"`
	Amount sdk.Coin       `json:"amount" yaml:"amount"`
}

var _ sdk.Msg = MsgSendLock{}

// NewMsgSendLock returns a MsgSendLock object.
func NewMsgSendLock(from, to sdk.AccAddress, amount sdk.Coin) MsgSendLock {
	return MsgSendLock{
		From:   from,
		To:     to,
		Amount: amount,
	}
}

// Route returns the name of the module.
func (m MsgSendLock) Route() string { return bank.ModuleName }

// Type returns a human-readable string for the message.
func (m MsgSendLock) Type() string { return "send_lock" }

// ValidateBasic runs stateless checks on the message.
func (m MsgSendLock) ValidateBasic() error {
	if m.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing from address")
	}
	if m.To.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing to address")
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgSendLock) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required.
func (m MsgSendLock) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.From}
}
