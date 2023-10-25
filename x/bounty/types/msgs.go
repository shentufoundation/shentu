package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateProgram  = "create_program"
	TypeMsgSubmitFinding  = "submit_finding"
	TypeMsgAcceptFinding  = "accept_finding"
	TypeMsgRejectFinding  = "reject_finding"
	TypeMsgCancelFinding  = "cancel_finding"
	TypeMsgReleaseFinding = "release_finding"
	TypeMsgEndProgram     = "end_program"
)

// NewMsgCreateProgram creates a new NewMsgCreateProgram instance.
// Delegator address and validator address are the same.
func NewMsgCreateProgram(name, pid, desc string, operator sdk.AccAddress, members []string, levels []BountyLevel) (*MsgCreateProgram, error) {
	return &MsgCreateProgram{
		Name:            name,
		Description:     desc,
		OperatorAddress: operator.String(),
		MemberAccounts:  members,
		ProgramId:       pid,
		BountyLevels:    levels,
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
func NewMsgEditProgram(name, pid, desc string, operator sdk.AccAddress, members []string, levels []BountyLevel) (*MsgCreateProgram, error) {
	return &MsgCreateProgram{
		Name:            name,
		Description:     desc,
		OperatorAddress: operator.String(),
		MemberAccounts:  members,
		ProgramId:       pid,
		BountyLevels:    levels,
	}, nil
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
func NewMsgSubmitFinding(pid, fid, title, desc string, operator sdk.AccAddress, level SeverityLevel) *MsgSubmitFinding {

	return &MsgSubmitFinding{
		ProgramId:        pid,
		FindingId:        fid,
		Title:            title,
		Description:      desc,
		SubmitterAddress: operator.String(),
		SeverityLevel:    level,
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

func NewMsgModifyFindingStatus(findingID string, hostAddr sdk.AccAddress) *MsgModifyFindingStatus {
	return &MsgModifyFindingStatus{
		FindingId:       findingID,
		OperatorAddress: hostAddr.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgModifyFindingStatus) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgModifyFindingStatus) Type() string { return TypeMsgAcceptFinding }

// GetSignBytes returns the message bytes to sign over.
func (msg MsgModifyFindingStatus) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements the sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
func (msg MsgModifyFindingStatus) GetSigners() []sdk.AccAddress {
	// host should sign the message
	hostAddr, err := sdk.AccAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{hostAddr}
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgModifyFindingStatus) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid issuer address (%s)", err.Error())
	}

	if len(msg.FindingId) == 0 {
		return errors.New("empty finding-id is not allowed")
	}
	return nil
}

// NewReleaseFinding release finding.
func NewReleaseFinding(fid, desc string, operator sdk.AccAddress) *MsgReleaseFinding {
	return &MsgReleaseFinding{
		FindingId:       fid,
		Description:     desc,
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

func NewMsgOpenProgram(accAddr sdk.AccAddress, pid string) *MsgModifyProgramStatus {
	return &MsgModifyProgramStatus{
		ProgramId:       pid,
		OperatorAddress: accAddr.String(),
	}
}

// Route implements sdk.Msg interface.
func (msg MsgModifyProgramStatus) Route() string { return RouterKey }

// Type implements sdk.Msg interface.
func (msg MsgModifyProgramStatus) Type() string { return TypeMsgEndProgram }

// GetSigners implements sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
func (msg MsgModifyProgramStatus) GetSigners() []sdk.AccAddress {
	cAddr, _ := sdk.AccAddressFromBech32(msg.OperatorAddress)
	return []sdk.AccAddress{cAddr}
}

// GetSignBytes implements the sdk.Msg interface, returns the message bytes to sign over.
func (msg MsgModifyProgramStatus) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgModifyProgramStatus) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid address (%s)", err.Error())
	}
	return nil
}
