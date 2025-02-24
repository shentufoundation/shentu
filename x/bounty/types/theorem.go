package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// DefaultStartingTheoremID is 1
	DefaultStartingTheoremID uint64 = 1
)

func NewTheorem(id uint64, proposer sdk.AccAddress, title, desc, code string, submitTime, grantEndTime time.Time, proofTime time.Duration) (Theorem, error) {

	proofEndTime := submitTime.Add(proofTime)
	return Theorem{
		Id:             id,
		Title:          title,
		Description:    desc,
		Code:           code,
		Status:         TheoremStatus_THEOREM_STATUS_GRANT_PERIOD,
		SubmitTime:     &submitTime,
		GrantEndTime:   &grantEndTime,
		ProofStartTime: &submitTime,
		ProofEndTime:   &proofEndTime,
		Proposer:       proposer.String(),
	}, nil
}

func NewProof(theoremId uint64, proofHash, prover string, submitTime time.Time, deposit sdk.Coins) (Proof, error) {
	return Proof{
		TheoremId:  theoremId,
		Id:         proofHash,
		Status:     ProofStatus_PROOF_STATUS_HASH_LOCK_PERIOD,
		SubmitTime: &submitTime,
		Prover:     prover,
		Deposit:    deposit,
	}, nil
}

func NewGrant(theoremID uint64, grantor sdk.AccAddress, amount sdk.Coins) Grant {
	return Grant{
		TheoremsId: theoremID,
		Grantor:    grantor.String(),
		Amount:     amount,
	}
}
