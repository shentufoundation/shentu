package types

import (
	"encoding/json"

	"gopkg.in/yaml.v2"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgProposeCertifier is the message for proposing new certifier.
type MsgProposeCertifier struct {
	Proposer    sdk.AccAddress `json:"proposer" yaml:"proposer"`
	Alias       string         `json:"alias" yaml:"alias"`
	Certifier   sdk.AccAddress `json:"certifier" yaml:"certifier"`
	Description string         `json:"description" yaml:"description"`
}

// NewMsgProposeCertifier returns a new certifier proposal message.
func NewMsgProposeCertifier(proposer, certifier sdk.AccAddress, alias string, description string) MsgProposeCertifier {
	return MsgProposeCertifier{
		Proposer:    proposer,
		Certifier:   certifier,
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
	if m.Certifier.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Certifier.String())
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
	return []sdk.AccAddress{m.Proposer}
}

// MsgCertifyValidator is the message for certifying a validator node.
type MsgCertifyValidator struct {
	Certifier sdk.AccAddress `json:"certifier" yaml:"certifier"`
	Validator crypto.PubKey  `json:"validator" yaml:"validator"`
}

type msgCertifyValidatorPretty struct {
	Certifier sdk.AccAddress `json:"certifier" yaml:"certifier"`
	Validator string         `json:"validator" yaml:"validator"`
}

// NewMsgCertifyValidator returns a new validator node certification message.
func NewMsgCertifyValidator(certifier sdk.AccAddress, validator crypto.PubKey) MsgCertifyValidator {
	return MsgCertifyValidator{
		Certifier: certifier,
		Validator: validator,
	}
}

// Route returns the module name.
func (m MsgCertifyValidator) Route() string { return ModuleName }

// Type returns the action name.
func (m MsgCertifyValidator) Type() string { return "certify_validator" }

// ValidateBasic runs stateless checks on the message.
func (m MsgCertifyValidator) ValidateBasic() error {
	if m.Validator == nil {
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
	return []sdk.AccAddress{m.Certifier}
}

// MarshalYAML implements a custom marshal yaml function due to consensus pubkey.
func (m MsgCertifyValidator) MarshalYAML() (interface{}, error) {
	d, err := yaml.Marshal(struct {
		Certifier sdk.AccAddress
		Validator string
	}{
		Certifier: m.Certifier,
		Validator: sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, m.Validator),
	})
	if err != nil {
		return nil, err
	}
	return string(d), nil
}

// Custom implementation due to the pubkey.
func (m MsgCertifyValidator) MarshalJSON() ([]byte, error) {
	var pk string
	var err error
	if m.Validator != nil {
		pk, err = sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, m.Validator)
		if err != nil {
			return nil, err
		}
	}
	return json.Marshal(struct {
		Certifier sdk.AccAddress
		Validator string
	}{
		m.Certifier,
		pk,
	})
}

// Custom implementation due to the pubkey.
func (m *MsgCertifyValidator) UnmarshalJSON(bz []byte) error {
	var alias msgCertifyValidatorPretty
	if err := json.Unmarshal(bz, &alias); err != nil {
		return err
	}
	if alias.Validator != "" {
		pk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, alias.Validator)
		if err != nil {
			return err
		}
		m.Validator = pk
	}
	m.Certifier = alias.Certifier
	return nil
}

// MsgDecertifyValidator is the message for de-certifying a validator node.
type MsgDecertifyValidator struct {
	Decertifier sdk.AccAddress `json:"decertifier" yaml:"decertifier"`
	Validator   crypto.PubKey  `json:"validator" yaml:"validator"`
}

type msgDecertifyValidatorPretty struct {
	Decertifier sdk.AccAddress `json:"decertifier" yaml:"decertifier"`
	Validator   string         `json:"validator" yaml:"validator"`
}

// NewMsgDecertifyValidator returns a new validator node de-certification message.
func NewMsgDecertifyValidator(decertifier sdk.AccAddress, validator crypto.PubKey) MsgDecertifyValidator {
	return MsgDecertifyValidator{
		Decertifier: decertifier,
		Validator:   validator,
	}
}

// Route returns the module name.
func (m MsgDecertifyValidator) Route() string { return ModuleName }

// Type returns the action name.
func (m MsgDecertifyValidator) Type() string { return "decertify_validator" }

// ValidateBasic runs stateless checks on the message.
func (m MsgDecertifyValidator) ValidateBasic() error {
	if m.Validator == nil {
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
	return []sdk.AccAddress{m.Decertifier}
}

// MarshalYAML implements a custom marshal yaml function due to consensus pubkey.
func (m MsgDecertifyValidator) MarshalYAML() (interface{}, error) {
	d, err := yaml.Marshal(struct {
		Decertifier sdk.AccAddress
		Validator   string
	}{
		Decertifier: m.Decertifier,
		Validator:   sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, m.Validator),
	})
	if err != nil {
		return nil, err
	}
	return string(d), nil
}

// Custom implementation due to the pubkey.
func (m MsgDecertifyValidator) MarshalJSON() ([]byte, error) {
	var pk string
	var err error
	if m.Validator != nil {
		pk, err = sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, m.Validator)
		if err != nil {
			return nil, err
		}
	}
	return json.Marshal(struct {
		Decertifier sdk.AccAddress
		Validator   string
	}{
		m.Decertifier,
		pk,
	})
}

// Custom implementation due to the pubkey.
func (m *MsgDecertifyValidator) UnmarshalJSON(bz []byte) error {
	var alias msgDecertifyValidatorPretty
	if err := json.Unmarshal(bz, &alias); err != nil {
		return err
	}
	if alias.Validator != "" {
		pk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, alias.Validator)
		if err != nil {
			return err
		}
		m.Validator = pk
	}
	m.Decertifier = alias.Decertifier
	return nil
}

// MsgCertifyGeneral is the message for issuing a general certificate.
type MsgCertifyGeneral struct {
	CertificateType    string         `json:"certificate_type" yaml:"certificate_type"`
	RequestContentType string         `json:"request_content_type" yaml:"request_content_type"`
	RequestContent     string         `json:"request_content" yaml:"request_content"`
	Description        string         `json:"description" yaml:"description"`
	Certifier          sdk.AccAddress `json:"certifier" yaml:"certiifer"`
}

// NewMsgCertifyGeneral returns a new general certification message.
func NewMsgCertifyGeneral(
	certificateType, requestContentType, requestContent, description string, certifier sdk.AccAddress,
) MsgCertifyGeneral {
	return MsgCertifyGeneral{
		CertificateType:    certificateType,
		RequestContentType: requestContentType,
		RequestContent:     requestContent,
		Description:        description,
		Certifier:          certifier,
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
	return []sdk.AccAddress{m.Certifier}
}

// MsgRevokeCertificate returns a certificate revoking operation.
type MsgRevokeCertificate struct {
	Revoker     sdk.AccAddress `json:"revoker" yaml:"revoker"`
	ID          CertificateID  `json:"id" yaml:"id"`
	Description string         `json:"description" yaml:"description"`
}

// NewMsgRevokeCertificate creates a new instance of MsgRevokeCertificate.
func NewMsgRevokeCertificate(revoker sdk.AccAddress, id CertificateID, description string) MsgRevokeCertificate {
	return MsgRevokeCertificate{
		Revoker:     revoker,
		ID:          id,
		Description: description,
	}
}

// ValidateBasic runs stateless checks on the message.
func (m MsgRevokeCertificate) ValidateBasic() error {
	if m.Revoker.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Revoker.String())
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
	return []sdk.AccAddress{m.Revoker}
}

// MsgCertifyCompilation is the message for certifying a compilation.
type MsgCertifyCompilation struct {
	SourceCodeHash string         `json:"sourcecodehash" yaml:"sourcecodehash"`
	Compiler       string         `json:"compiler" yaml:"compiler"`
	BytecodeHash   string         `json:"bytecodehash" yaml:"bytecodehash"`
	Description    string         `json:"description" yaml:"description"`
	Certifier      sdk.AccAddress `json:"certifier" yaml:"certifier"`
}

// NewMsgCertifyCompilation returns a compilation certificate message.
func NewMsgCertifyCompilation(sourceCodeHash, compiler, bytecodeHash, description string, certifier sdk.AccAddress) MsgCertifyCompilation {
	return MsgCertifyCompilation{
		SourceCodeHash: sourceCodeHash,
		Compiler:       compiler,
		BytecodeHash:   bytecodeHash,
		Description:    description,
		Certifier:      certifier,
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
	return []sdk.AccAddress{m.Certifier}
}

// MsgCertifyPlatform is the message for certifying a validator's host platform.
type MsgCertifyPlatform struct {
	Certifier sdk.AccAddress `json:"certifier" yaml:"certifier"`
	Validator crypto.PubKey  `json:"validator" yaml:"validator"`
	Platform  string         `json:"platform" yaml:"platform"`
}

type msgCertifyPlatformPretty struct {
	Certifier sdk.AccAddress `json:"certifier" yaml:"certifier"`
	Validator string         `json:"validator" yaml:"validator"`
	Platform  string         `json:"platform" yaml:"platform"`
}

// NewMsgCertifyPlatform returns a new validator host platform certification
// message.
func NewMsgCertifyPlatform(certifier sdk.AccAddress, validator crypto.PubKey, platform string) MsgCertifyPlatform {
	return MsgCertifyPlatform{
		Certifier: certifier,
		Validator: validator,
		Platform:  platform,
	}
}

// Route returns the module name.
func (m MsgCertifyPlatform) Route() string { return ModuleName }

// Type returns the action name.
func (m MsgCertifyPlatform) Type() string { return "certify_platform" }

// ValidateBasic runs stateless checks on the message.
func (m MsgCertifyPlatform) ValidateBasic() error {
	if m.Validator == nil {
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
	return []sdk.AccAddress{m.Certifier}
}

// MarshalYAML implements a custom marshal yaml function due to consensus pubkey.
func (m MsgCertifyPlatform) MarshalYAML() (interface{}, error) {
	d, err := yaml.Marshal(struct {
		Certifier sdk.AccAddress
		Validator string
		Platform  string
	}{
		Certifier: m.Certifier,
		Validator: sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, m.Validator),
		Platform:  m.Platform,
	})
	if err != nil {
		return nil, err
	}
	return string(d), nil
}

// Custom implementation due to the pubkey.
func (m MsgCertifyPlatform) MarshalJSON() ([]byte, error) {
	var pk string
	var err error
	if m.Validator != nil {
		pk, err = sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, m.Validator)
		if err != nil {
			return nil, err
		}
	}
	return json.Marshal(struct {
		Certifier sdk.AccAddress
		Validator string
		Platform  string
	}{
		m.Certifier,
		pk,
		m.Platform,
	})
}

// Custom implementation due to the pubkey.
func (m *MsgCertifyPlatform) UnmarshalJSON(bz []byte) error {
	var alias msgCertifyPlatformPretty
	if err := json.Unmarshal(bz, &alias); err != nil {
		return err
	}
	if alias.Validator != "" {
		pk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, alias.Validator)
		if err != nil {
			return err
		}
		m.Validator = pk
	}
	m.Certifier = alias.Certifier
	m.Platform = alias.Platform
	return nil
}
