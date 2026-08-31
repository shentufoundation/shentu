package types

import (
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgWithdrawRewards = "withdraw_rewards"
)

// NewMsgWithdrawRewards creates a new MsgWithdrawRewards instance.
func NewMsgWithdrawRewards(sender sdk.AccAddress) *MsgWithdrawRewards {
	return &MsgWithdrawRewards{
		From: sender.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgWithdrawRewards) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgWithdrawRewards) Type() string { return TypeMsgWithdrawRewards }

// GetSigners implements the sdk.Msg interface.
func (msg MsgWithdrawRewards) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgWithdrawRewards) ValidateBasic() error {
	if msg.From == "" {
		return ErrEmptySender
	}

	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return errorsmod.Wrapf(err, "invalid sender address %s", msg.From)
	}

	return nil
}
