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
func NewMsgCreateProgram(name, creatorAddr, pid string, detail ProgramDetail, memberAddrs []string) (*MsgCreateProgram, error) {

	return &MsgCreateProgram{
		Name:           name,
		Detail:         detail,
		CreatorAddress: creatorAddr,
		MemberAccounts: memberAddrs,
		ProgramId:      pid,
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
	_, err := sdk.AccAddressFromBech32(msg.CreatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid issuer address (%s)", err.Error())
	}
	return nil
}

// NewMsgEditProgram edit a program.
func NewMsgEditProgram(name, creatorAddr string, desc ProgramDetail, memberAddrs []string) MsgEditProgram {

	return MsgEditProgram{
		Name:           name,
		Detail:         desc,
		CreatorAddress: creatorAddr,
		MemberAccounts: memberAddrs,
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
	cAddr, err := sdk.AccAddressFromBech32(msg.CreatorAddress)
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
	_, err := sdk.AccAddressFromBech32(msg.CreatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid issuer address (%s)", err.Error())
	}
	return nil
}

// NewMsgSubmitFinding submit a new finding.
func NewMsgSubmitFinding(
	programId, findingId, title string, detail FindingDetail, accAddr sdk.AccAddress) *MsgSubmitFinding {

	return &MsgSubmitFinding{
		ProgramId:        programId,
		FindingId:        findingId,
		Title:            title,
		Detail:           detail,
		SubmitterAddress: accAddr.String(),
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

func NewMsgHostAcceptFinding(findingID string, hostAddr sdk.AccAddress) *MsgHostAcceptFinding {
	return &MsgHostAcceptFinding{
		FindingId:   findingID,
		HostAddress: hostAddr.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgHostAcceptFinding) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgHostAcceptFinding) Type() string { return TypeMsgAcceptFinding }

// GetSignBytes returns the message bytes to sign over.
func (msg MsgHostAcceptFinding) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements the sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
func (msg MsgHostAcceptFinding) GetSigners() []sdk.AccAddress {
	// host should sign the message
	hostAddr, err := sdk.AccAddressFromBech32(msg.HostAddress)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{hostAddr}
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgHostAcceptFinding) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.HostAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid issuer address (%s)", err.Error())
	}

	if len(msg.FindingId) == 0 {
		return errors.New("empty finding-id is not allowed")
	}
	return nil
}

func NewMsgHostRejectFinding(findingID string, hostAddr sdk.AccAddress) *MsgHostRejectFinding {
	return &MsgHostRejectFinding{
		FindingId:   findingID,
		HostAddress: hostAddr.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgHostRejectFinding) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgHostRejectFinding) Type() string { return TypeMsgRejectFinding }

// GetSignBytes returns the message bytes to sign over.
func (msg MsgHostRejectFinding) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements the sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
func (msg MsgHostRejectFinding) GetSigners() []sdk.AccAddress {
	// host should sign the message
	hostAddr, err := sdk.AccAddressFromBech32(msg.HostAddress)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{hostAddr}
}

// ValidateBasic implements the sdk.Msg interface.
func (msg *MsgHostRejectFinding) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.HostAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid issuer address (%s)", err.Error())
	}

	if len(msg.FindingId) == 0 {
		return errors.New("empty finding-id is not allowed")
	}
	return nil
}

// NewMsgCancelFinding cancel a specific finding
func NewMsgCancelFinding(accAddr sdk.AccAddress, findingID string) *MsgCancelFinding {
	return &MsgCancelFinding{
		SubmitterAddress: accAddr.String(),
		FindingId:        findingID,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgCancelFinding) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgCancelFinding) Type() string { return TypeMsgCancelFinding }

// GetSigners implements the sdk.Msg interface
func (msg MsgCancelFinding) GetSigners() []sdk.AccAddress {
	// creator should sign the message
	cAddr, err := sdk.AccAddressFromBech32(msg.SubmitterAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{cAddr}
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgCancelFinding) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgCancelFinding) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.SubmitterAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid issuer address (%s)", err.Error())
	}
	if len(msg.FindingId) == 0 {
		return errors.New("empty finding-id is not allowed")
	}
	return nil
}

// NewReleaseFinding release finding.
func NewReleaseFinding(
	hostAddr, fid string, findingDesc, findingPoc, findingComment string,
) *MsgReleaseFinding {
	return &MsgReleaseFinding{
		FindingId:   fid,
		Desc:        findingDesc,
		Poc:         findingPoc,
		HostAddress: hostAddr,
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
	cAddr, err := sdk.AccAddressFromBech32(msg.HostAddress)
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
	_, err := sdk.AccAddressFromBech32(msg.HostAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid issuer address (%s)", err.Error())
	}

	if len(msg.FindingId) == 0 {
		return errors.New("empty fid is not allowed")
	}
	return nil
}

func NewMsgOpenProgram(accAddr sdk.AccAddress, pid string) *MsgOpenProgram {
	return &MsgOpenProgram{
		OpenAddress: accAddr.String(),
		ProgramId:   pid,
	}
}

// Route implements sdk.Msg interface.
func (msg MsgOpenProgram) Route() string { return RouterKey }

// Type implements sdk.Msg interface.
func (msg MsgOpenProgram) Type() string { return TypeMsgEndProgram }

// GetSigners implements sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
func (msg MsgOpenProgram) GetSigners() []sdk.AccAddress {
	cAddr, _ := sdk.AccAddressFromBech32(msg.OpenAddress)
	return []sdk.AccAddress{cAddr}
}

// GetSignBytes implements the sdk.Msg interface, returns the message bytes to sign over.
func (msg MsgOpenProgram) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgOpenProgram) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OpenAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid address (%s)", err.Error())
	}
	return nil
}

func NewMsgCloseProgram(accAddr sdk.AccAddress, pid string) *MsgCloseProgram {
	return &MsgCloseProgram{
		CloseAddress: accAddr.String(),
		ProgramId:    pid,
	}
}

// Route implements sdk.Msg interface.
func (msg MsgCloseProgram) Route() string { return RouterKey }

// Type implements sdk.Msg interface.
func (msg MsgCloseProgram) Type() string { return TypeMsgEndProgram }

// GetSigners implements sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
func (msg MsgCloseProgram) GetSigners() []sdk.AccAddress {
	cAddr, _ := sdk.AccAddressFromBech32(msg.CloseAddress)
	return []sdk.AccAddress{cAddr}
}

// GetSignBytes implements the sdk.Msg interface, returns the message bytes to sign over.
func (msg MsgCloseProgram) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgCloseProgram) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.CloseAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid address (%s)", err.Error())
	}
	return nil
}
