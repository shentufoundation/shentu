package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/irisnet/irismod/modules/nft/types"
)

// NewMsgCreateNFTAdmin returns a new certifier proposal message.
func NewMsgCreateNFTAdmin(issuer, address string) *MsgCreateNFTAdmin {
	return &MsgCreateNFTAdmin{
		Issuer:issuer,
		Address: address,
	}
}

// Route returns the module name.
func (m MsgCreateNFTAdmin) Route() string { return types.ModuleName }

// Type returns the action name.
func (m MsgCreateNFTAdmin) Type() string { return "create_admin" }

// ValidateBasic runs stateless checks on the message.
func (m MsgCreateNFTAdmin) ValidateBasic() error {
	certifierAddr, err := sdk.AccAddressFromBech32(m.Issuer)
	if err != nil {
		panic(err)
	}
	if certifierAddr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, certifierAddr.String())
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgCreateNFTAdmin) GetSignBytes() []byte {
	bz := types.ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

// GetSigners defines whose signature is required.
func (m MsgCreateNFTAdmin) GetSigners() []sdk.AccAddress {
	proposerAddr, err := sdk.AccAddressFromBech32(m.Issuer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{proposerAddr}
}

// NewMsgRevokeNFTAdmin returns a new certifier proposal message.
func NewMsgRevokeNFTAdmin(issuer, address string) *MsgRevokeNFTAdmin {
	return &MsgRevokeNFTAdmin{
		Revoker:issuer,
		Address: address,
	}
}

// Route returns the module name.
func (m MsgRevokeNFTAdmin) Route() string { return types.ModuleName }

// Type returns the action name.
func (m MsgRevokeNFTAdmin) Type() string { return "revoke_admin" }

// ValidateBasic runs stateless checks on the message.
func (m MsgRevokeNFTAdmin) ValidateBasic() error {
	certifierAddr, err := sdk.AccAddressFromBech32(m.Revoker)
	if err != nil {
		panic(err)
	}
	if certifierAddr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, certifierAddr.String())
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgRevokeNFTAdmin) GetSignBytes() []byte {
	bz := types.ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

// GetSigners defines whose signature is required.
func (m MsgRevokeNFTAdmin) GetSigners() []sdk.AccAddress {
	proposerAddr, err := sdk.AccAddressFromBech32(m.Revoker)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{proposerAddr}
}
