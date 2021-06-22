package types

import (
	fmt "fmt"

	"github.com/gogo/protobuf/proto"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
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
func NewMsgIssueCertificate(
	content Content, compiler, bytecodeHash, description string, certifier sdk.AccAddress,
) *MsgIssueCertificate {
	msg, ok := content.(proto.Message)
	if !ok {
		panic(fmt.Errorf("%T does not implement proto.Message", content))
	}
	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		panic(err)
	}
	return &MsgIssueCertificate{
		Content:      any,
		Compiler:     compiler,
		BytecodeHash: bytecodeHash,
		Description:  description,
		Certifier:    certifier.String(),
	}
}

// Route returns the module name.
func (m MsgIssueCertificate) Route() string { return ModuleName }

// Type returns the action name.
func (m MsgIssueCertificate) Type() string { return "issue_certificate" }

// ValidateBasic runs stateless checks on the message.
func (m MsgIssueCertificate) ValidateBasic() error {
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

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces.
func (m MsgIssueCertificate) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var content Content
	return unpacker.UnpackAny(m.Content, &content)
}

// NewMsgRevokeCertificate creates a new instance of MsgRevokeCertificate.
func NewMsgRevokeCertificate(revoker sdk.AccAddress, id uint64, description string) *MsgRevokeCertificate {
	return &MsgRevokeCertificate{
		Revoker:     revoker.String(),
		Id:          id,
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
