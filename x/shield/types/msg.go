package types

import (
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

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgWithdrawRewards) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgWithdrawRewards) ValidateBasic() error {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	if from.Empty() {
		return ErrEmptySender
	}

	return nil
}
