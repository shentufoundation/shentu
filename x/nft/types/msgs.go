package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/nft/types"
)

// NewMsgCreateAdmin returns a new create NFT admin message.
func NewMsgCreateAdmin(creator, address string) *MsgCreateAdmin {
	return &MsgCreateAdmin{
		Creator: creator,
		Address: address,
	}
}

// Route returns the module name.
func (m MsgCreateAdmin) Route() string { return types.ModuleName }

// Type returns the action name.
func (m MsgCreateAdmin) Type() string { return "create_admin" }

// ValidateBasic runs stateless checks on the message.
func (m MsgCreateAdmin) ValidateBasic() error {
	certifierAddr, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		panic(err)
	}
	if certifierAddr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, certifierAddr.String())
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgCreateAdmin) GetSignBytes() []byte {
	bz := types.ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

// GetSigners defines whose signature is required.
func (m MsgCreateAdmin) GetSigners() []sdk.AccAddress {
	proposerAddr, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{proposerAddr}
}

// NewMsgRevokeAdmin returns a new revoke NFT admin message.
func NewMsgRevokeAdmin(issuer, address string) *MsgRevokeAdmin {
	return &MsgRevokeAdmin{
		Revoker: issuer,
		Address: address,
	}
}

// Route returns the module name.
func (m MsgRevokeAdmin) Route() string { return types.ModuleName }

// Type returns the action name.
func (m MsgRevokeAdmin) Type() string { return "revoke_admin" }

// ValidateBasic runs stateless checks on the message.
func (m MsgRevokeAdmin) ValidateBasic() error {
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
func (m MsgRevokeAdmin) GetSignBytes() []byte {
	bz := types.ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

// GetSigners defines whose signature is required.
func (m MsgRevokeAdmin) GetSigners() []sdk.AccAddress {
	proposerAddr, err := sdk.AccAddressFromBech32(m.Revoker)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{proposerAddr}
}
