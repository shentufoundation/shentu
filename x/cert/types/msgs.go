package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gogoproto/proto"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	proposer sdk.AccAddress,
) *MsgUpdateCertifier {
	msg := &MsgUpdateCertifier{
		Authority:   authority.String(),
		Certifier:   certifier.String(),
		Description: description,
		Operation:   operation.String(),
	}
	if len(proposer) > 0 {
		msg.Proposer = proposer.String()
	}
	return msg
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

	if _, err := AddOrRemoveFromString(m.Operation); err != nil {
		return err
	}

	if m.Proposer != "" {
		proposerAddr, err := sdk.AccAddressFromBech32(m.Proposer)
		if err != nil {
			return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, m.Proposer)
		}
		if proposerAddr.Empty() {
			return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "<empty>")
		}
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

