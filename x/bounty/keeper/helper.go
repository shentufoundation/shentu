package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"cosmossdk.io/collections"

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

// HasActiveProofs checks if a theorem has any proofs in hash lock or detail period
func (k Keeper) HasActiveProofs(ctx context.Context, theoremId uint64) (bool, string, error) {
	var activeProofId string
	hasActiveProof := false

	rng := collections.NewPrefixedPairRange[uint64, string](theoremId)
	err := k.ProofsByTheorem.Walk(ctx, rng, func(key collections.Pair[uint64, string], _ []byte) (bool, error) {
		proof, err := k.Proofs.Get(ctx, key.K2())
		if err != nil {
			return false, err
		}
		if proof.Status == types.ProofStatus_PROOF_STATUS_HASH_LOCK_PERIOD ||
			proof.Status == types.ProofStatus_PROOF_STATUS_HASH_DETAIL_PERIOD {
			hasActiveProof = true
			activeProofId = proof.Id
			return true, nil
		}
		return false, nil
	})

	return hasActiveProof, activeProofId, err
}

// GetProgramFindings retrieves all findings associated with a program
func (k Keeper) GetProgramFindings(ctx context.Context, programID string) ([]string, error) {
	var findings []string

	// Create a range for all keys with the given program ID prefix
	rng := collections.NewPrefixedPairRange[string, string](programID)

	// Walk through all program-finding pairs for this program
	err := k.ProgramFindings.Walk(ctx, rng, func(key collections.Pair[string, string]) (bool, error) {
		findings = append(findings, key.K2())
		return false, nil
	})

	return findings, err
}
