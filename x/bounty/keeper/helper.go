package keeper

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func (k Keeper) GetProgramFingerprintHash(program *types.Program) string {
	programFingerprint := &types.ProgramFingerprint{
		ProgramId:    program.ProgramId,
		Name:         program.Name,
		Detail:       program.Detail,
		AdminAddress: program.AdminAddress,
		Status:       program.Status,
	}

	bz := k.cdc.MustMarshal(programFingerprint)
	hash := sha256.Sum256(bz)
	return hex.EncodeToString(hash[:])
}

func (k Keeper) GetFindingFingerprintHash(finding *types.Finding) string {
	findingFingerprint := &types.FindingFingerprint{
		ProgramId:   finding.ProgramId,
		FindingId:   finding.FindingId,
		FindingHash: finding.FindingHash,
		Status:      finding.Status,
		PaymentHash: finding.PaymentHash,
	}

	bz := k.cdc.MustMarshal(findingFingerprint)
	hash := sha256.Sum256(bz)
	return hex.EncodeToString(hash[:])
}

func (k Keeper) GetProofHash(theoremId uint64, prover, detail string) string {
	proofHash := &types.ProofHash{
		TheoremId: theoremId,
		Prover:    prover,
		Detail:    detail,
	}
	bz := k.cdc.MustMarshal(proofHash)
	hash := sha256.Sum256(bz)
	return hex.EncodeToString(hash[:])
}
