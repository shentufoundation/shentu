package types

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/gogoproto/proto"
)

const (
	TypeMsgUpdateCertifier    = "update_certifier"
	TypeMsgCertifyValidator   = "certify_validator"
	TypeMsgDecertifyValidator = "decertify_validator"
	TypeMsgCertifyGeneral     = "certify_general"
	TypeMsgRevokeCertificate  = "revoke_certificate"
	TypeMsgCertifyCompilation = "certify_compilation"
)

// NewMsgUpdateCertifier returns a new governance-authorized certifier update message.
func NewMsgUpdateCertifier(
	authority sdk.AccAddress,
	certifier sdk.AccAddress,
	description string,
	operation AddOrRemove,
) *MsgUpdateCertifier {
	return &MsgUpdateCertifier{
		Authority:   authority.String(),
		Certifier:   certifier.String(),
		Description: description,
		Operation:   operation.ToProto(),
	}
}

// Route returns the module name.
func (m MsgUpdateCertifier) Route() string { return ModuleName }

// Type returns the action name.
func (m MsgUpdateCertifier) Type() string { return TypeMsgUpdateCertifier }

// ValidateBasic runs stateless checks on the message.
func (m MsgUpdateCertifier) ValidateBasic() error {
	authorityAddr, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, m.Authority)
	}
	if authorityAddr.Empty() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "<empty>")
	}

	certifierAddr, err := sdk.AccAddressFromBech32(m.Certifier)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, m.Certifier)
	}
	if certifierAddr.Empty() {
		return ErrEmptyCertifier
	}

	if _, err := AddOrRemoveFromProto(m.Operation); err != nil {
		return err
	}
	return nil
}

// GetSigners defines whose signature is required.
func (m MsgUpdateCertifier) GetSigners() []sdk.AccAddress {
	authorityAddr, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authorityAddr}
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
	certifierAddr, err := sdk.AccAddressFromBech32(m.Certifier)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, m.Certifier)
	}
	if certifierAddr.Empty() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "<empty>")
	}

	certType, err := certificateTypeFromAny(m.Content)
	if err != nil {
		return err
	}

	if certType == CertificateTypeCompilation {
		if m.Compiler == "" {
			return ErrCompiler
		}
		if m.BytecodeHash == "" {
			return ErrBytecodeHash
		}
	}
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
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, m.Revoker)
	}
	if revokerAddr.Empty() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "<empty>")
	}
	if m.Id == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "certificate id must be positive")
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

func certificateTypeFromAny(contentAny *codectypes.Any) (CertificateType, error) {
	if contentAny == nil || contentAny.TypeUrl == "" {
		return CertificateTypeNil, ErrInvalidRequestContentType
	}

	if content, ok := contentAny.GetCachedValue().(Content); ok {
		switch content.(type) {
		case *Compilation:
			return CertificateTypeCompilation, nil
		case *Auditing:
			return CertificateTypeAuditing, nil
		case *Proof:
			return CertificateTypeProof, nil
		case *OracleOperator:
			return CertificateTypeOracleOperator, nil
		case *ShieldPoolCreator:
			return CertificateTypeShieldPoolCreator, nil
		case *Identity:
			return CertificateTypeIdentity, nil
		case *General:
			return CertificateTypeGeneral, nil
		case *BountyAdmin:
			return CertificateTypeBountyAdmin, nil
		case *OpenMath:
			return CertificateTypeOpenMath, nil
		}
	}

	typeName := contentAny.TypeUrl
	if idx := strings.LastIndex(typeName, "/"); idx >= 0 {
		typeName = typeName[idx+1:]
	}

	switch typeName {
	case "shentu.cert.v1alpha1.Compilation", "Compilation":
		return CertificateTypeCompilation, nil
	case "shentu.cert.v1alpha1.Auditing", "Auditing":
		return CertificateTypeAuditing, nil
	case "shentu.cert.v1alpha1.Proof", "Proof":
		return CertificateTypeProof, nil
	case "shentu.cert.v1alpha1.OracleOperator", "OracleOperator":
		return CertificateTypeOracleOperator, nil
	case "shentu.cert.v1alpha1.ShieldPoolCreator", "ShieldPoolCreator":
		return CertificateTypeShieldPoolCreator, nil
	case "shentu.cert.v1alpha1.Identity", "Identity":
		return CertificateTypeIdentity, nil
	case "shentu.cert.v1alpha1.General", "General":
		return CertificateTypeGeneral, nil
	case "shentu.cert.v1alpha1.BountyAdmin", "BountyAdmin":
		return CertificateTypeBountyAdmin, nil
	case "shentu.cert.v1alpha1.OpenMath", "OpenMath":
		return CertificateTypeOpenMath, nil
	default:
		return CertificateTypeNil, ErrInvalidRequestContentType
	}
}
