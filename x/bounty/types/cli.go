package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewProgram(pid, name, detail string, admin sdk.AccAddress,
	status ProgramStatus, createTime time.Time) (Program, error) {

	return Program{
		ProgramId:    pid,
		Name:         name,
		Detail:       detail,
		AdminAddress: admin.String(),
		Status:       status,
		CreateTime:   createTime,
	}, nil
}

// ValidProgramStatus returns true if the program status is valid and false
// otherwise.
func ValidProgramStatus(status ProgramStatus) bool {
	if status == ProgramStatusInactive ||
		status == ProgramStatusActive ||
		status == ProgramStatusClosed {
		return true
	}
	return false
}

func NewTheorem(id uint64, proposer sdk.AccAddress, title, desc, code string, submitTime, endTime time.Time) (Theorem, error) {
	return Theorem{
		Id:          id,
		Title:       title,
		Description: desc,
		Code:        code,
		Status:      TheoremStatus_THEOREM_STATUS_PROOF_PERIOD,
		SubmitTime:  &submitTime,
		EndTime:     &endTime,
		Proposer:    proposer.String(),
	}, nil
}

func NewProof(theoremId uint64, proofHash, prover string, submitTime, endTime time.Time, deposit sdk.Coins) (Proof, error) {
	return Proof{
		TheoremId:  theoremId,
		Id:         proofHash,
		Status:     ProofStatus_PROOF_STATUS_HASH_LOCK_PERIOD,
		SubmitTime: &submitTime,
		EndTime:    &endTime,
		Prover:     prover,
		Deposit:    deposit,
	}, nil
}

func NewGrant(theoremID uint64, grantor sdk.AccAddress, amount sdk.Coins) Grant {
	return Grant{
		TheoremId: theoremID,
		Grantor:   grantor.String(),
		Amount:    amount,
	}
}

func NewDeposit(proofID string, depositor sdk.AccAddress, amount sdk.Coins) Deposit {
	return Deposit{
		ProofId:   proofID,
		Depositor: depositor.String(),
		Amount:    amount,
	}
}
