package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gogoproto/proto"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgProposeCertifier   = "propose_certifier"
	TypeMsgCertifyValidator   = "certify_validator"
	TypeMsgDecertifyValidator = "decertify_validator"
	TypeMsgCertifyGeneral     = "certify_general"
	TypeMsgRevokeCertificate  = "revoke_certificate"
	TypeMsgCertifyCompilation = "certify_compilation"
	TypeMsgCertifyPlatform    = "certify_platform"
)

// NewMsgProposeCertifier returns a new certifier proposal message.
func NewMsgProposeCertifier(proposer, certifier sdk.AccAddress, alias string, description string) *MsgProposeCertifier {
	return &MsgProposeCertifier{
		Proposer:    proposer.String(),
		Certifier:   certifier.String(),
		Alias:       alias,
		Description: description,
	}
}

// Route returns the module name.
func (m MsgProposeCertifier) Route() string { return ModuleName }

// Type returns the action name.
func (m MsgProposeCertifier) Type() string { return "propose_certifier" }

// ValidateBasic runs stateless checks on the message.
func (m MsgProposeCertifier) ValidateBasic() error {
	certifierAddr, err := sdk.AccAddressFromBech32(m.Certifier)
	if err != nil {
		panic(err)
	}
	if certifierAddr.Empty() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, certifierAddr.String())
	}
	return nil
}

// GetSigners defines whose signature is required.
func (m MsgProposeCertifier) GetSigners() []sdk.AccAddress {
	proposerAddr, err := sdk.AccAddressFromBech32(m.Proposer)
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
	contentAny, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		panic(err)
	}
	return &MsgIssueCertificate{
		Content:      contentAny,
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
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, revokerAddr.String())
	}
	return nil
}

// Route returns the module name.
func (m MsgRevokeCertificate) Route() string { return ModuleName }

// Type returns the action name.
func (m MsgRevokeCertificate) Type() string { return "revoke_certificate" }

// GetSigners defines whose signature is required.
func (m MsgRevokeCertificate) GetSigners() []sdk.AccAddress {
	revokerAddr, err := sdk.AccAddressFromBech32(m.Revoker)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{revokerAddr}
}

// NewMsgCertifyPlatform returns a new validator host platform certification
// message.
func NewMsgCertifyPlatform(certifier sdk.AccAddress, pk cryptotypes.PubKey, platform string) (*MsgCertifyPlatform, error) {
	var pkAny *codectypes.Any
	if pk != nil {
		var err error
		if pkAny, err = codectypes.NewAnyWithValue(pk); err != nil {
			return nil, err
		}
	}

	return &MsgCertifyPlatform{Certifier: certifier.String(), ValidatorPubkey: pkAny, Platform: platform}, nil
}

// Route returns the module name.
func (m MsgCertifyPlatform) Route() string { return ModuleName }

// Type returns the action name.
func (m MsgCertifyPlatform) Type() string { return "certify_platform" }

// ValidateBasic runs stateless checks on the message.
func (m MsgCertifyPlatform) ValidateBasic() error {
	certifierAddr, err := sdk.AccAddressFromBech32(m.Certifier)
	if err != nil {
		panic(err)
	}
	if certifierAddr.Empty() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "<empty>")
	}

	if m.ValidatorPubkey == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "<empty>")
	}
	return nil
}

// GetSigners defines whose signature is required.
func (m MsgCertifyPlatform) GetSigners() []sdk.AccAddress {
	certifierAddr, err := sdk.AccAddressFromBech32(m.Certifier)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{certifierAddr}
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (m MsgCertifyPlatform) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pubKey cryptotypes.PubKey
	return unpacker.UnpackAny(m.ValidatorPubkey, &pubKey)
}
