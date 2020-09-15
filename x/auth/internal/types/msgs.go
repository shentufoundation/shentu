package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	ModuleName = "auth"
	RouterKey  = ModuleName
)

// MsgTriggerVesting triggers vesting of the specified account
type MsgTriggerVesting struct {
	Certifier sdk.AccAddress `json:"certifier" yaml:"certifier"`
	Account   sdk.AccAddress `json:"account_address" yaml:"account_address"`
}

var _ sdk.Msg = MsgTriggerVesting{}

func NewMsgTriggerVesting(certifier, account sdk.AccAddress) MsgTriggerVesting {
	return MsgTriggerVesting{
		Certifier: certifier,
		Account:   account,
	}
}

// Route returns the name of the module.
func (m MsgTriggerVesting) Route() string { return ModuleName }

// Type returns a human-readable string for the message.
func (m MsgTriggerVesting) Type() string { return "trigger_vesting" }

// ValidateBasic runs stateless checks on the message.
func (m MsgTriggerVesting) ValidateBasic() error {
	if m.Certifier.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing from address")
	}
	if m.Account.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing account address")
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgTriggerVesting) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required.
func (m MsgTriggerVesting) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Certifier}
}
