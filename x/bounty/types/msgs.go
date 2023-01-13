package types

import (
	"errors"
	"fmt"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateProgram     = "create_program"
	TypeMsgSubmitFinding     = "submit_finding"
	TypeMsgWithdrawalFinding = "withdrawal_finding"
	TypeMsgReactivateFinding = "reactivate_finding"
)

// NewMsgCreateProgram creates a new NewMsgCreateProgram instance.
// Delegator address and validator address are the same.
func NewMsgCreateProgram(
	creatorAddress string, description string, encKey []byte, commissionRate sdk.Dec, deposit sdk.Coins,
	submissionEndTime, judgingEndTime, claimEndTime time.Time,
) (*MsgCreateProgram, error) {
	var encAny *codectypes.Any
	if encKey != nil {
		encKeyMsg := EciesPubKey{
			PubKey: encKey,
		}

		var err error
		if encAny, err = codectypes.NewAnyWithValue(&encKeyMsg); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("encKey is empty")
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

// NewMsgSubmitFinding submit a new finding.
func NewMsgSubmitFinding(
	submitterAddress string, title, description string, programId uint64, severityLevel int32, poc string,
) (*MsgSubmitFinding, error) {
	if programId == 0 {
		return nil, errors.New("empty pid is not allowed")
	}

	return &MsgSubmitFinding{
		Title:            title,
		Desc:             description,
		Pid:              programId,
		SeverityLevel:    SeverityLevel(severityLevel),
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

// NewMsgWithdrawalFinding withdrawal a specific finding
func NewMsgWithdrawalFinding(accAddr sdk.AccAddress, findingId uint64) *MsgWithdrawalFinding {
	return &MsgWithdrawalFinding{
		From: accAddr.String(),
		Fid:  findingId,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgWithdrawalFinding) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgWithdrawalFinding) Type() string { return TypeMsgWithdrawalFinding }

// GetSigners implements the sdk.Msg interface
func (msg MsgWithdrawalFinding) GetSigners() []sdk.AccAddress {
	// creator should sign the message
	cAddr, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{cAddr}
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgWithdrawalFinding) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgWithdrawalFinding) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid issuer address (%s)", err.Error())
	}
	if msg.Fid == 0 {
		return errors.New("empty fid is not allowed")
	}
	return nil
}

// NewMsgReactivateFinding reactivate a specific finding
func NewMsgReactivateFinding(accAddr sdk.AccAddress, findingId uint64) *MsgReactivateFinding {
	return &MsgReactivateFinding{
		From: accAddr.String(),
		Fid:  findingId,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgReactivateFinding) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgReactivateFinding) Type() string { return TypeMsgReactivateFinding }

// GetSigners implements the sdk.Msg interface
func (msg MsgReactivateFinding) GetSigners() []sdk.AccAddress {
	// creator should sign the message
	cAddr, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{cAddr}
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgReactivateFinding) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgReactivateFinding) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid issuer address (%s)", err.Error())
	}
	if msg.Fid == 0 {
		return errors.New("empty fid is not allowed")
	}
	return nil
}
