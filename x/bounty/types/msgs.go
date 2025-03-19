package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgCreateProgram      = "create_program"
	TypeMsgEditProgram        = "edit_program"
	TypeMsgActivateProgram    = "activate_program"
	TypeMsgCloseProgram       = "close_program"
	TypeMsgSubmitFinding      = "submit_finding"
	TypeMsgEditFinding        = "edit_finding"
	TypeMsgActivateFinding    = "activate_finding"
	TypeMsgConfirmFinding     = "confirm_finding"
	TypeMsgConfirmFindingPaid = "confirm_finding_paid"
	TypeMsgCloseFinding       = "close_finding"
	TypeMsgPublishFinding     = "publish_finding"
)

var (
	_, _, _, _       sdk.Msg = &MsgCreateProgram{}, &MsgEditProgram{}, &MsgActivateProgram{}, &MsgCloseProgram{}
	_, _, _, _, _, _ sdk.Msg = &MsgSubmitFinding{}, &MsgEditFinding{}, &MsgActivateFinding{}, &MsgConfirmFinding{}, &MsgCloseFinding{}, &MsgPublishFinding{}
	_, _, _, _       sdk.Msg = &MsgCreateTheorem{}, &MsgGrant{}, &MsgSubmitProofHash{}, &MsgSubmitProofDetail{}
)

// NewMsgCreateProgram creates a new NewMsgCreateProgram instance.
// Delegator address and validator address are the same.
func NewMsgCreateProgram(pid, name, detail string, operator sdk.AccAddress) *MsgCreateProgram {
	return &MsgCreateProgram{
		ProgramId:       pid,
		Name:            name,
		Detail:          detail,
		OperatorAddress: operator.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgCreateProgram) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgCreateProgram) Type() string { return TypeMsgCreateProgram }

// NewMsgEditProgram edit a program.
func NewMsgEditProgram(pid, name, detail string, operator sdk.AccAddress) *MsgEditProgram {
	return &MsgEditProgram{
		ProgramId:       pid,
		Name:            name,
		Detail:          detail,
		OperatorAddress: operator.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgEditProgram) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgEditProgram) Type() string { return TypeMsgEditProgram }

// NewMsgSubmitFinding submit a new finding.
func NewMsgSubmitFinding(pid, fid, hash string, operator sdk.AccAddress, level SeverityLevel) *MsgSubmitFinding {
	return &MsgSubmitFinding{
		ProgramId:       pid,
		FindingId:       fid,
		FindingHash:     hash,
		OperatorAddress: operator.String(),
		SeverityLevel:   level,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgSubmitFinding) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgSubmitFinding) Type() string { return TypeMsgSubmitFinding }

// NewMsgEditFinding submit a new finding.
func NewMsgEditFinding(fid, hash, paymentHash string, operator sdk.AccAddress, level SeverityLevel) *MsgEditFinding {
	return &MsgEditFinding{
		FindingId:       fid,
		FindingHash:     hash,
		OperatorAddress: operator.String(),
		SeverityLevel:   level,
		PaymentHash:     paymentHash,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgEditFinding) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgEditFinding) Type() string { return TypeMsgEditFinding }

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

func NewMsgActivateFinding(findingID string, hostAddr sdk.AccAddress) *MsgActivateFinding {
	return &MsgActivateFinding{
		FindingId:       findingID,
		OperatorAddress: hostAddr.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgActivateFinding) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgActivateFinding) Type() string { return TypeMsgActivateFinding }

func NewMsgConfirmFinding(findingID, fingerprint string, hostAddr sdk.AccAddress) *MsgConfirmFinding {
	return &MsgConfirmFinding{
		FindingId:       findingID,
		OperatorAddress: hostAddr.String(),
		Fingerprint:     fingerprint,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgConfirmFinding) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgConfirmFinding) Type() string { return TypeMsgConfirmFinding }

func NewMsgConfirmFindingPaid(findingID string, hostAddr sdk.AccAddress) *MsgConfirmFindingPaid {
	return &MsgConfirmFindingPaid{
		FindingId:       findingID,
		OperatorAddress: hostAddr.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgConfirmFindingPaid) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgConfirmFindingPaid) Type() string { return TypeMsgConfirmFindingPaid }

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

// NewMsgPublishFinding publish finding.
func NewMsgPublishFinding(fid, desc, poc string, operator sdk.AccAddress) *MsgPublishFinding {
	return &MsgPublishFinding{
		FindingId:       fid,
		Description:     desc,
		ProofOfConcept:  poc,
		OperatorAddress: operator.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgPublishFinding) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgPublishFinding) Type() string { return TypeMsgPublishFinding }

func NewMsgCreateTheorem(title, desc, code, proposer string, initialGrant sdk.Coins) *MsgCreateTheorem {
	return &MsgCreateTheorem{
		Title:        title,
		Description:  desc,
		Code:         code,
		InitialGrant: initialGrant,
		Proposer:     proposer,
	}
}

func NewMsgGrant(theoremID uint64, grantor string, amount sdk.Coins) *MsgGrant {
	return &MsgGrant{
		TheoremId: theoremID,
		Grantor:   grantor,
		Amount:    amount,
	}
}

func NewMsgSubmitProofHash(theoremID uint64, prover, hash string, amount sdk.Coins) *MsgSubmitProofHash {
	return &MsgSubmitProofHash{
		TheoremId: theoremID,
		Prover:    prover,
		ProofHash: hash,
		Deposit:   amount,
	}
}

func NewMsgSubmitProofDetail(proofID, prover, detail string) *MsgSubmitProofDetail {
	return &MsgSubmitProofDetail{
		ProofId: proofID,
		Prover:  prover,
		Detail:  detail,
	}
}

func NewMsgSubmitProofVerification(proofID string, status ProofStatus, checker string) *MsgSubmitProofVerification {
	return &MsgSubmitProofVerification{
		ProofId: proofID,
		Status:  status,
		Checker: checker,
	}
}

func NewMsgWithdrawReward(address string) *MsgWithdrawReward {
	return &MsgWithdrawReward{
		Address: address,
	}
}
