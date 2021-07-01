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

// NewMsgIssueCertificate returns a new certification message.
func NewMsgIssueCertificate(denomID, tokenID, name, uri, content, description string, certifier sdk.AccAddress) *MsgIssueCertificate {
	return &MsgIssueCertificate{
		DenomId:     denomID,
		TokenId:     tokenID,
		Name:        name,
		Uri:         uri,
		Content:     content,
		Description: description,
		Certifier:   certifier.String(),
	}
}

// Route returns the module name.
func (m MsgIssueCertificate) Route() string { return ModuleName }

// Type returns the action name.
func (m MsgIssueCertificate) Type() string { return "issue_certificate" }

// ValidateBasic runs stateless checks on the message.
func (m MsgIssueCertificate) ValidateBasic() error {
	certifierAddr, err := sdk.AccAddressFromBech32(m.Certifier)
	if err != nil {
		panic(err)
	}
	if certifierAddr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, certifierAddr.String())
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgIssueCertificate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

// GetSigners defines whose signature is required.
func (m MsgIssueCertificate) GetSigners() []sdk.AccAddress {
	certifierAddr, err := sdk.AccAddressFromBech32(m.Certifier)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{certifierAddr}
}

// NewMsgIssueCertificate returns a new certification message.
func NewMsgEditCertificate(denomID, tokenID, name, uri, content, description string, owner sdk.AccAddress) *MsgEditCertificate {
	return &MsgEditCertificate{
		DenomId:     denomID,
		TokenId:     tokenID,
		Name:        name,
		Uri:         uri,
		Content:     content,
		Description: description,
		Owner:       owner.String(),
	}
}

// Route returns the module name.
func (m MsgEditCertificate) Route() string { return ModuleName }

// Type returns the action name.
func (m MsgEditCertificate) Type() string { return "edit_certificate" }

// ValidateBasic runs stateless checks on the message.
func (m MsgEditCertificate) ValidateBasic() error {
	ownerAddr, err := sdk.AccAddressFromBech32(m.Owner)
	if err != nil {
		panic(err)
	}
	if ownerAddr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, ownerAddr.String())
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgEditCertificate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

// GetSigners defines whose signature is required.
func (m MsgEditCertificate) GetSigners() []sdk.AccAddress {
	ownerAddr, err := sdk.AccAddressFromBech32(m.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{ownerAddr}
}

// NewMsgRevokeCertificate creates a new instance of MsgRevokeCertificate.
func NewMsgRevokeCertificate(denomID, tokenID, description string, revoker sdk.AccAddress) *MsgRevokeCertificate {
	return &MsgRevokeCertificate{
		DenomId:     denomID,
		TokenId:     tokenID,
		Revoker:     revoker.String(),
		Description: description,
	}
}

// ValidateBasic runs stateless checks on the message.
func (m MsgRevokeCertificate) ValidateBasic() error {
	revokerAddr, err := sdk.AccAddressFromBech32(m.Revoker)
	if err != nil {
		panic(err)
	}
	if revokerAddr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, revokerAddr.String())
	}
	return nil
}

// Route returns the module name.
func (m MsgRevokeCertificate) Route() string { return ModuleName }

// Type returns the action name.
func (m MsgRevokeCertificate) Type() string { return "revoke_certificate" }

// GetSignBytes encodes the message for signing.
func (m MsgRevokeCertificate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

// GetSigners defines whose signature is required.
func (m MsgRevokeCertificate) GetSigners() []sdk.AccAddress {
	revokerAddr, err := sdk.AccAddressFromBech32(m.Revoker)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{revokerAddr}
}
