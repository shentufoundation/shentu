package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_, _, _, _       sdk.Msg = &MsgCreateProgram{}, &MsgEditProgram{}, &MsgActivateProgram{}, &MsgCloseProgram{}
	_, _, _, _, _, _ sdk.Msg = &MsgSubmitFinding{}, &MsgEditFinding{}, &MsgActivateFinding{}, &MsgConfirmFinding{}, &MsgCloseFinding{}, &MsgPublishFinding{}
	_, _             sdk.Msg = &MsgCreateTheorem{}, &MsgGrant{}
	_, _, _          sdk.Msg = &MsgSubmitProofHash{}, &MsgSubmitProofDetail{}, &MsgSubmitProofVerification{}
	_                sdk.Msg = &MsgWithdrawReward{}
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

// NewMsgEditProgram edit a program.
func NewMsgEditProgram(pid, name, detail string, operator sdk.AccAddress) *MsgEditProgram {
	return &MsgEditProgram{
		ProgramId:       pid,
		Name:            name,
		Detail:          detail,
		OperatorAddress: operator.String(),
	}
}

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

func NewMsgActivateProgram(pid string, accAddr sdk.AccAddress) *MsgActivateProgram {
	return &MsgActivateProgram{
		ProgramId:       pid,
		OperatorAddress: accAddr.String(),
	}
}

func NewMsgCloseProgram(pid string, accAddr sdk.AccAddress) *MsgCloseProgram {
	return &MsgCloseProgram{
		ProgramId:       pid,
		OperatorAddress: accAddr.String(),
	}
}

func NewMsgActivateFinding(findingID string, hostAddr sdk.AccAddress) *MsgActivateFinding {
	return &MsgActivateFinding{
		FindingId:       findingID,
		OperatorAddress: hostAddr.String(),
	}
}

func NewMsgConfirmFinding(findingID, fingerprint string, hostAddr sdk.AccAddress) *MsgConfirmFinding {
	return &MsgConfirmFinding{
		FindingId:       findingID,
		OperatorAddress: hostAddr.String(),
		Fingerprint:     fingerprint,
	}
}

func NewMsgConfirmFindingPaid(findingID string, hostAddr sdk.AccAddress) *MsgConfirmFindingPaid {
	return &MsgConfirmFindingPaid{
		FindingId:       findingID,
		OperatorAddress: hostAddr.String(),
	}
}

func NewMsgCloseFinding(findingID string, hostAddr sdk.AccAddress) *MsgCloseFinding {
	return &MsgCloseFinding{
		FindingId:       findingID,
		OperatorAddress: hostAddr.String(),
	}
}

// NewMsgPublishFinding publish finding.
func NewMsgPublishFinding(fid, desc, poc string, operator sdk.AccAddress) *MsgPublishFinding {
	return &MsgPublishFinding{
		FindingId:       fid,
		Description:     desc,
		ProofOfConcept:  poc,
		OperatorAddress: operator.String(),
	}
}

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
