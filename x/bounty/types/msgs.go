package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateProgram   = "create_program"
	TypeMsgActivateProgram = "activate_program"
	TypeMsgCloseProgram    = "close_program"
	TypeMsgSubmitFinding   = "submit_finding"
	TypeMsgEditFinding     = "edit_finding"
	TypeMsgConfirmFinding  = "confirm_finding"
	TypeMsgCloseFinding    = "close_finding"
	TypeMsgReleaseFinding  = "release_finding"
)

// NewMsgCreateProgram creates a new NewMsgCreateProgram instance.
// Delegator address and validator address are the same.
func NewMsgCreateProgram(pid, name, detail string, operator sdk.AccAddress, levels []BountyLevel) *MsgCreateProgram {
	return &MsgCreateProgram{
		ProgramId:       pid,
		Name:            name,
		Detail:          detail,
		BountyLevels:    levels,
		OperatorAddress: operator.String(),
	}
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
	cAddr, err := sdk.AccAddressFromBech32(msg.OperatorAddress)
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
	_, err := sdk.AccAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid issuer address (%s)", err.Error())
	}
	return nil
}

// NewMsgEditProgram edit a program.
func NewMsgEditProgram(pid, name, detail string, operator sdk.AccAddress, levels []BountyLevel) *MsgEditProgram {
	return &MsgEditProgram{
		ProgramId:       pid,
		Name:            name,
		Detail:          detail,
		BountyLevels:    levels,
		OperatorAddress: operator.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgEditProgram) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgEditProgram) Type() string { return TypeMsgCreateProgram }

// GetSigners implements the sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
// If the validator address is not same as delegator's, then the validator must
// sign the msg as well.
func (msg MsgEditProgram) GetSigners() []sdk.AccAddress {
	// creator should sign the message
	cAddr, err := sdk.AccAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{cAddr}
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgEditProgram) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgEditProgram) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid issuer address (%s)", err.Error())
	}
	return nil
}

// NewMsgSubmitFinding submit a new finding.
func NewMsgSubmitFinding(pid, fid, title, detail, hash string, operator sdk.AccAddress, level SeverityLevel) *MsgSubmitFinding {
	return &MsgSubmitFinding{
		ProgramId:        pid,
		FindingId:        fid,
		Title:            title,
		FindingHash:      hash,
		SubmitterAddress: operator.String(),
		SeverityLevel:    level,
		Detail:           detail,
	}
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
	if len(msg.ProgramId) == 0 {
		return errors.New("empty pid is not allowed")
	}
	return nil
}

// NewMsgEditFinding submit a new finding.
func NewMsgEditFinding(pid, fid, title, detail, hash string, operator sdk.AccAddress, level SeverityLevel) *MsgEditFinding {
	return &MsgEditFinding{
		ProgramId:        pid,
		FindingId:        fid,
		Title:            title,
		FindingHash:      hash,
		SubmitterAddress: operator.String(),
		SeverityLevel:    level,
		Detail:           detail,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgEditFinding) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgEditFinding) Type() string { return TypeMsgEditFinding }

// GetSigners implements the sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
// If the validator address is not same as delegator's, then the validator must
// sign the msg as well.
func (msg MsgEditFinding) GetSigners() []sdk.AccAddress {
	// creator should sign the message
	cAddr, err := sdk.AccAddressFromBech32(msg.SubmitterAddress)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{cAddr}
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgEditFinding) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgEditFinding) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.SubmitterAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid issuer address (%s)", err.Error())
	}
	if len(msg.ProgramId) == 0 {
		return errors.New("empty pid is not allowed")
	}
	return nil
}

func NewMsgActivateProgram(pid string, accAddr sdk.AccAddress) *MsgActivateProgram {
	return &MsgActivateProgram{
		ProgramId:       pid,
		OperatorAddress: accAddr.String(),
	}
}

// Route implements sdk.Msg interface.
func (msg MsgActivateProgram) Route() string { return RouterKey }

// Type implements sdk.Msg interface.
func (msg MsgActivateProgram) Type() string { return TypeMsgActivateProgram }

// GetSigners implements sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
func (msg MsgActivateProgram) GetSigners() []sdk.AccAddress {
	cAddr, _ := sdk.AccAddressFromBech32(msg.OperatorAddress)
	return []sdk.AccAddress{cAddr}
}

// GetSignBytes implements the sdk.Msg interface, returns the message bytes to sign over.
func (msg MsgActivateProgram) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgActivateProgram) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid address (%s)", err.Error())
	}
	return nil
}

func NewMsgCloseProgram(pid string, accAddr sdk.AccAddress) *MsgCloseProgram {
	return &MsgCloseProgram{
		ProgramId:       pid,
		OperatorAddress: accAddr.String(),
	}
}

// Route implements sdk.Msg interface.
func (msg MsgCloseProgram) Route() string { return RouterKey }

// Type implements sdk.Msg interface.
func (msg MsgCloseProgram) Type() string { return TypeMsgCloseProgram }

// GetSigners implements sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
func (msg MsgCloseProgram) GetSigners() []sdk.AccAddress {
	cAddr, _ := sdk.AccAddressFromBech32(msg.OperatorAddress)
	return []sdk.AccAddress{cAddr}
}

// GetSignBytes implements the sdk.Msg interface, returns the message bytes to sign over.
func (msg MsgCloseProgram) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgCloseProgram) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid address (%s)", err.Error())
	}
	return nil
}

func NewMsgConfirmFinding(findingID string, hostAddr sdk.AccAddress) *MsgConfirmFinding {
	return &MsgConfirmFinding{
		FindingId:       findingID,
		OperatorAddress: hostAddr.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgConfirmFinding) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgConfirmFinding) Type() string { return TypeMsgConfirmFinding }

// GetSignBytes returns the message bytes to sign over.
func (msg MsgConfirmFinding) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements the sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
func (msg MsgConfirmFinding) GetSigners() []sdk.AccAddress {
	// host should sign the message
	hostAddr, err := sdk.AccAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{hostAddr}
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgConfirmFinding) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid issuer address (%s)", err.Error())
	}

	if len(msg.FindingId) == 0 {
		return errors.New("empty finding-id is not allowed")
	}
	return nil
}

func NewMsgCloseFinding(findingID string, hostAddr sdk.AccAddress) *MsgCloseFinding {
	return &MsgCloseFinding{
		FindingId:       findingID,
		OperatorAddress: hostAddr.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgCloseFinding) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgCloseFinding) Type() string { return TypeMsgCloseFinding }

// GetSignBytes returns the message bytes to sign over.
func (msg MsgCloseFinding) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements the sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
func (msg MsgCloseFinding) GetSigners() []sdk.AccAddress {
	// host should sign the message
	hostAddr, err := sdk.AccAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{hostAddr}
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgCloseFinding) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid issuer address (%s)", err.Error())
	}

	if len(msg.FindingId) == 0 {
		return errors.New("empty finding-id is not allowed")
	}
	return nil
}

// NewMsgReleaseFinding release finding.
func NewMsgReleaseFinding(fid, desc, poc string, operator sdk.AccAddress) *MsgReleaseFinding {
	return &MsgReleaseFinding{
		FindingId:       fid,
		Description:     desc,
		ProofOfConcept:  poc,
		OperatorAddress: operator.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgReleaseFinding) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgReleaseFinding) Type() string { return TypeMsgReleaseFinding }

// GetSigners implements the sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
func (msg MsgReleaseFinding) GetSigners() []sdk.AccAddress {
	// releaser should sign the message
	cAddr, err := sdk.AccAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{cAddr}
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgReleaseFinding) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgReleaseFinding) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid issuer address (%s)", err.Error())
	}

	if len(msg.FindingId) == 0 {
		return errors.New("empty fid is not allowed")
	}
	return nil
}
