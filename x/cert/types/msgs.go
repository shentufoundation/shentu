package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgProposeCertifier   = "propose_certifier"
	TypeMsgCertifyValidator   = "certify_validator"
	TypeMsgDecertifyValidator = "decertify_validator"
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
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

// GetSigners defines whose signature is required.
func (m MsgProposeCertifier) GetSigners() []sdk.AccAddress {
	proposerAddr, err := sdk.AccAddressFromBech32(m.Proposer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{proposerAddr}
}

// NewMsgCertifyValidator returns a new validator node certification message.
func NewMsgCertifyValidator(certifier sdk.AccAddress, pk cryptotypes.PubKey) (*MsgCertifyValidator, error) {
	var pkAny *codectypes.Any
	if pk != nil {
		var err error
		if pkAny, err = codectypes.NewAnyWithValue(pk); err != nil {
			return nil, err
		}
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
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

// GetSigners defines whose signature is required.
func (m MsgCertifyValidator) GetSigners() []sdk.AccAddress {
	certifierAddr, err := sdk.AccAddressFromBech32(m.Certifier)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{certifierAddr}
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (m MsgCertifyValidator) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pubKey cryptotypes.PubKey
	return unpacker.UnpackAny(m.Pubkey, &pubKey)
}

// NewMsgDecertifyValidator returns a new validator node de-certification message.
func NewMsgDecertifyValidator(decertifier sdk.AccAddress, pk cryptotypes.PubKey) (*MsgDecertifyValidator, error) {
	var pkAny *codectypes.Any
	if pk != nil {
		var err error
		if pkAny, err = codectypes.NewAnyWithValue(pk); err != nil {
			return nil, err
		}
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
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

// GetSigners defines whose signature is required.
func (m MsgDecertifyValidator) GetSigners() []sdk.AccAddress {
	decertifierAddr, err := sdk.AccAddressFromBech32(m.Decertifier)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{decertifierAddr}
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (m MsgDecertifyValidator) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pubKey cryptotypes.PubKey
	return unpacker.UnpackAny(m.Pubkey, &pubKey)
}
