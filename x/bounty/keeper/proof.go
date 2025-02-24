package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func (k Keeper) SubmitProofHash(ctx context.Context, theoremID uint64, proofID, prover string, deposit sdk.Coins) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	submitTime := sdkCtx.BlockHeader().Time

	// Check if theorem exists
	theorem, err := k.Theorems.Get(ctx, theoremID)
	if err != nil {
		return err
	}
	// Check theorem is still depositable
	if theorem.Status != types.TheoremStatus_THEOREM_STATUS_GRANT_PERIOD &&
		theorem.Status != types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD {
		// todo add error info
		return types.ErrTheoremStatusInvalid
	}

	// Check coins to be deposited match the proposal's deposit params
	_, err = k.Params.Get(ctx)
	if err != nil {
		return err
	}

	proof, err := types.NewProof(theoremID, proofID, prover, submitTime, deposit)
	if err != nil {
		return err
	}

	return k.SetProof(ctx, proof)
}

func (k Keeper) SubmitProofDetail(ctx context.Context, proofId string, detail string) error {
	// Check if proof exists
	proof, err := k.Proofs.Get(ctx, proofId)
	if err != nil {
		return err
	}
	// Check proof status
	if proof.Status != types.ProofStatus_PROOF_STATUS_HASH_LOCK_PERIOD {
		return types.ErrProofStatusInvalid
	}

	// update proof
	proof.ProofContent = detail

	return k.SetProof(ctx, proof)
}

func (k Keeper) SetProof(ctx context.Context, proof types.Proof) error {
	return k.Proofs.Set(ctx, proof.Id, proof)
}
