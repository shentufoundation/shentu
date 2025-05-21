package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewProgram(pid, name, detail string, admin sdk.AccAddress,
	status ProgramStatus, createTime time.Time) Program {

	return Program{
		ProgramId:    pid,
		Name:         name,
		Detail:       detail,
		AdminAddress: admin.String(),
		Status:       status,
		CreateTime:   createTime,
	}
}

func NewFinding(pid, fid, title, detail, hash string, operator sdk.AccAddress, createTime time.Time, level SeverityLevel) Finding {
	return Finding{
		ProgramId:        pid,
		FindingId:        fid,
		Title:            title,
		FindingHash:      hash,
		SubmitterAddress: operator.String(),
		SeverityLevel:    level,
		Status:           FindingStatusSubmitted,
		Detail:           detail,
		CreateTime:       createTime,
	}
}

func NewTheorem(id uint64, proposer sdk.AccAddress, title, desc, code string, submitTime, endTime time.Time) Theorem {
	return Theorem{
		Id:          id,
		Title:       title,
		Description: desc,
		Code:        code,
		Status:      TheoremStatus_THEOREM_STATUS_PROOF_PERIOD,
		SubmitTime:  &submitTime,
		EndTime:     &endTime,
		Proposer:    proposer.String(),
	}
}

func NewProof(theoremId uint64, proofHash, prover string, submitTime, endTime time.Time, deposit sdk.Coins) Proof {
	return Proof{
		TheoremId:  theoremId,
		Id:         proofHash,
		Status:     ProofStatus_PROOF_STATUS_HASH_LOCK_PERIOD,
		SubmitTime: &submitTime,
		EndTime:    &endTime,
		Prover:     prover,
		Deposit:    deposit,
	}
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
