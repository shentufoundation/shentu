package types

import (
	"encoding/json"

	"github.com/tendermint/tendermint/crypto"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, certifierAddr.String())
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgProposeCertifier) GetSignBytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required.
func (m MsgProposeCertifier) GetSigners() []sdk.AccAddress {
	proposerAddr, err := sdk.AccAddressFromBech32(m.Proposer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{proposerAddr}
}

type msgCertifyValidatorPretty struct {
	Certifier sdk.AccAddress `json:"certifier" yaml:"certifier"`
	Validator string         `json:"validator" yaml:"validator"`
}

// NewMsgCertifyValidator returns a new validator node certification message.
func NewMsgCertifyValidator(certifier sdk.AccAddress, pk crypto.PubKey) (*MsgCertifyValidator, error) {
	pkAny, err := codectypes.PackAny(pk)
	if err != nil {
		return nil, err
	}
		
	return &MsgCertifyValidator{Certifier: certifier.String(), Pubkey: pkAny}, nil
}

// Route returns the module name.
func (m MsgCertifyValidator) Route() string { return ModuleName }

// Type returns the action name.
func (m MsgCertifyValidator) Type() string { return "certify_validator" }

// ValidateBasic runs stateless checks on the message.
func (m MsgCertifyValidator) ValidateBasic() error {
	if m.Pubkey == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidPubKey, "<empty>")
	}

	certifierAddr, err := sdk.AccAddressFromBech32(m.Certifier)
	if err != nil {
		panic(err)
	}
	if certifierAddr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "<empty>")
	}

	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgCertifyValidator) GetSignBytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required.
func (m MsgCertifyValidator) GetSigners() []sdk.AccAddress {
	certifierAddr, err := sdk.AccAddressFromBech32(m.Certifier)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{certifierAddr}
}

type msgDecertifyValidatorPretty struct {
	Decertifier sdk.AccAddress `json:"decertifier" yaml:"decertifier"`
	Validator   string         `json:"validator" yaml:"validator"`
}

// NewMsgDecertifyValidator returns a new validator node de-certification message.
func NewMsgDecertifyValidator(decertifier sdk.AccAddress, pk crypto.PubKey) (*MsgDecertifyValidator, error) {
	pkAny, err := codectypes.PackAny(pk)
	if err != nil {
		return nil, err
	}
	
	return &MsgDecertifyValidator{Decertifier: decertifier.String(), Pubkey: pkAny}, nil
}

// Route returns the module name.
func (m MsgDecertifyValidator) Route() string { return ModuleName }

// Type returns the action name.
func (m MsgDecertifyValidator) Type() string { return "decertify_validator" }

// ValidateBasic runs stateless checks on the message.
func (m MsgDecertifyValidator) ValidateBasic() error {
	if m.Pubkey == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidPubKey, "<empty>")
	}

	certifierAddr, err := sdk.AccAddressFromBech32(m.Decertifier)
	if err != nil {
		panic(err)
	}
	if certifierAddr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "<empty>")
	}

	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgDecertifyValidator) GetSignBytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required.
func (m MsgDecertifyValidator) GetSigners() []sdk.AccAddress {
	decertifierAddr, err := sdk.AccAddressFromBech32(m.Decertifier)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{decertifierAddr}
}

// NewMsgCertifyGeneral returns a new general certification message.
func NewMsgCertifyGeneral(
	certificateType, requestContentType, requestContent, description string, certifier sdk.AccAddress,
) *MsgCertifyGeneral {
	return &MsgCertifyGeneral{
		CertificateType:    certificateType,
		RequestContentType: requestContentType,
		RequestContent:     requestContent,
		Description:        description,
		Certifier:          certifier.String(),
	}
}

// Route returns the module name.
func (m MsgCertifyGeneral) Route() string { return ModuleName }

// Type returns the action name.
func (m MsgCertifyGeneral) Type() string { return "certify_general" }

// ValidateBasic runs stateless checks on the message.
func (m MsgCertifyGeneral) ValidateBasic() error {
	if certificateType := CertificateTypeFromString(m.CertificateType); certificateType == CertificateTypeNil {
		return ErrInvalidCertificateType
	}
	if requestContentType := RequestContentTypeFromString(m.RequestContentType); requestContentType == RequestContentTypeNil {
		return ErrInvalidRequestContentType
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgCertifyGeneral) GetSignBytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required.
func (m MsgCertifyGeneral) GetSigners() []sdk.AccAddress {
	certifierAddr, err := sdk.AccAddressFromBech32(m.Certifier)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{certifierAddr}
}

// NewMsgRevokeCertificate creates a new instance of MsgRevokeCertificate.
func NewMsgRevokeCertificate(revoker sdk.AccAddress, id CertificateID, description string) *MsgRevokeCertificate {
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
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required.
func (m MsgRevokeCertificate) GetSigners() []sdk.AccAddress {
	revokerAddr, err := sdk.AccAddressFromBech32(m.Revoker)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{revokerAddr}
}

// NewMsgCertifyCompilation returns a compilation certificate message.
func NewMsgCertifyCompilation(sourceCodeHash, compiler, bytecodeHash, description string, certifier sdk.AccAddress) *MsgCertifyCompilation {
	return &MsgCertifyCompilation{
		SourceCodeHash: sourceCodeHash,
		Compiler:       compiler,
		BytecodeHash:   bytecodeHash,
		Description:    description,
		Certifier:      certifier.String(),
	}
}

// Route returns the module name.
func (m MsgCertifyCompilation) Route() string { return ModuleName }

// Type returns the action name.
func (m MsgCertifyCompilation) Type() string { return "certify_compilation" }

// ValidateBasic runs stateless checks on the message.
func (m MsgCertifyCompilation) ValidateBasic() error {
	if m.SourceCodeHash == "" {
		return sdkerrors.Wrap(ErrSourceCodeHash, "<empty>")
	}
	if m.Compiler == "" {
		return sdkerrors.Wrap(ErrCompiler, "<empty>")
	}
	if m.BytecodeHash == "" {
		return sdkerrors.Wrap(ErrBytecodeHash, "<empty>")
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgCertifyCompilation) GetSignBytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required.
func (m MsgCertifyCompilation) GetSigners() []sdk.AccAddress {
	certifierAddr, err := sdk.AccAddressFromBech32(m.Certifier)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{certifierAddr}
}

type msgCertifyPlatformPretty struct {
	Certifier sdk.AccAddress `json:"certifier" yaml:"certifier"`
	Validator string         `json:"validator" yaml:"validator"`
	Platform  string         `json:"platform" yaml:"platform"`
}

// NewMsgCertifyPlatform returns a new validator host platform certification
// message.
func NewMsgCertifyPlatform(certifier sdk.AccAddress, pk crypto.PubKey, platform string) (*MsgCertifyPlatform, error) {
	pkAny, err := codectypes.PackAny(pk)
	if err != nil {
		return nil, err
	}
	
	return &MsgCertifyPlatform{Certifier: certifier.String(), ValidatorPubkey: pkAny, Platform:  platform}, nil
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
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "<empty>")
	}

	if m.ValidatorPubkey == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "<empty>")
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgCertifyPlatform) GetSignBytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required.
func (m MsgCertifyPlatform) GetSigners() []sdk.AccAddress {
	certifierAddr, err := sdk.AccAddressFromBech32(m.Certifier)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{certifierAddr}
}
