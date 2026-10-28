package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = (*MsgBanAddress)(nil)
var _ sdk.Msg = (*MsgUnbanAddress)(nil)

// GetSigners returns the address of the message signer. Only the governance
// module account is authorized to sign these messages in practice; the actual
// authorization is enforced by the message server.
func (msg MsgBanAddress) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

// GetSigners returns the address of the message signer.
func (msg MsgUnbanAddress) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

// ValidateBasic performs stateless validation on MsgBanAddress.
func (msg MsgBanAddress) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid authority")
	}

	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid address")
	}
	if addr.Empty() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "address is empty")
	}

	if msg.Category == CATEGORY_UNSPECIFIED {
		return errorsmod.Wrap(ErrInvalidCategory, "category must be set")
	}

	return nil
}

// ValidateBasic performs stateless validation on MsgUnbanAddress.
func (msg MsgUnbanAddress) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid authority")
	}

	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid address")
	}
	if addr.Empty() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "address is empty")
	}

	return nil
}
