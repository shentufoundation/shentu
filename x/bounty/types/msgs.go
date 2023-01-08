package types

import (
	"errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgCreateProgram = "create_program"

	TypeMsgSubmitFinding = "submit_finding"
)

// NewMsgCreateProgram creates a new MsgCreateValidator instance.
// Delegator address and validator address are the same.
func NewMsgCreateProgram(
	creatorAddress string, description string, encKey cryptotypes.PubKey, commissionRate sdk.Dec, deposit sdk.Coins,
	submissionEndTime, judgingEndTime, claimEndTime time.Time,
) (*MsgCreateProgram, error) {
	var encAny *codectypes.Any
	if encKey != nil {
		var err error
		if encAny, err = codectypes.NewAnyWithValue(encKey); err != nil {
			return nil, err
		}
	}
	return &MsgCreateProgram{
		Description:       description,
		CommissionRate:    commissionRate,
		SubmissionEndTime: submissionEndTime,
		CreatorAddress:    creatorAddress,
		EncryptionKey:     encAny,
		Deposit:           deposit,
	}, nil
}

// Route implements the sdk.Msg interface.
func (msg MsgCreateProgram) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgCreateProgram) Type() string { return TypeMsgCreateProgram }

// GetSigners implements the sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
// If the validator address is not same as delegator's, then the validator must
// sign the msg as well.
func (msg MsgCreateProgram) GetSigners() []sdk.AccAddress {
	// creator should sign the message
	cAddr, err := sdk.AccAddressFromBech32(msg.CreatorAddress)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{cAddr}
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgCreateProgram) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgCreateProgram) ValidateBasic() error {
	// TODO: implement ValidateBasic
	return nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgCreateProgram) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pubKey cryptotypes.PubKey
	return unpacker.UnpackAny(msg.EncryptionKey, &pubKey)
}

// NewMsgSubmitFinding submit a new finding.
func NewMsgSubmitFinding(
	submitterAddress string, title, description string, programId uint64, severityLevel uint32, poc string,
) (*MsgSubmitFinding, error) {
	if programId == 0 {
		return nil, errors.New("empty pid is not allowed")
	}

	return &MsgSubmitFinding{
		Title:            title,
		Desc:             description,
		Pid:              programId,
		SeverityLevel:    severityLevel,
		Poc:              poc,
		SubmitterAddress: submitterAddress,
	}, nil
}

// Route implements the sdk.Msg interface.
func (msg MsgSubmitFinding) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgSubmitFinding) Type() string { return TypeMsgSubmitFinding }

// GetSigners implements the sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
// If the validator address is not same as delegator's, then the validator must
// sign the msg as well.
func (msg MsgSubmitFinding) GetSigners() []sdk.AccAddress {
	// creator should sign the message
	cAddr, err := sdk.AccAddressFromBech32(msg.SubmitterAddress)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{cAddr}
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgSubmitFinding) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgSubmitFinding) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.SubmitterAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid issuer address (%s)", err.Error())
	}

	if msg.Pid == 0 {
		return errors.New("empty pid is not allowed")
	}
	return nil
}
