package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// RouterKey is they name of the bank module
const RouterKey = bank.ModuleName

// MsgLockedSend transfers coins and have them vesting
// in the receiver's manual vesting account.
type MsgLockedSend struct {
	From   sdk.AccAddress `json:"from" yaml:"from"`
	To     sdk.AccAddress `json:"to" yaml:"to"`
	Amount sdk.Coins      `json:"amount" yaml:"amount"`
}

var _ sdk.Msg = MsgLockedSend{}

// NewMsgLockedSend returns a MsgLockedSend object.
func NewMsgLockedSend(from, to sdk.AccAddress, amount sdk.Coins) MsgLockedSend {
	return MsgLockedSend{
		From:   from,
		To:     to,
		Amount: amount,
	}
}

// Route returns the name of the module.
func (m MsgLockedSend) Route() string { return RouterKey }

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
