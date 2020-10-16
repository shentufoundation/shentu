package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

const (
	ModuleName = auth.ModuleName
	RouterKey  = ModuleName
)

// MsgUnlock unlocks the specified amount in a manual vesting account.
type MsgUnlock struct {
	Issuer       sdk.AccAddress `json:"issuer" yaml:"issuer"`
	Account      sdk.AccAddress `json:"account_address" yaml:"account_address"`
	UnlockAmount sdk.Coins      `json:"unlock_amount" yaml:"unlock_amount"`
}

var _ sdk.Msg = MsgUnlock{}

// NewMsgUnlock returns a MsgUnlock object.
func NewMsgUnlock(issuer, account sdk.AccAddress, unlockAmount sdk.Coins) MsgUnlock {
	return MsgUnlock{
		Issuer:       issuer,
		Account:      account,
		UnlockAmount: unlockAmount,
	}
}

// Route returns the name of the module.
func (m MsgUnlock) Route() string { return ModuleName }

// Type returns a human-readable string for the message.
func (m MsgUnlock) Type() string { return "unlock" }

// ValidateBasic runs stateless checks on the message.
func (m MsgUnlock) ValidateBasic() error {
	if m.Issuer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing from address")
	}
	if m.Account.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing account address")
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgUnlock) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required.
func (m MsgUnlock) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Issuer}
}
