package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
)

func (k Keeper) DeleteProof(ctx context.Context, proofID string) error {
	proof, err := k.Proofs.Get(ctx, proofID)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return status.Errorf(codes.NotFound, "proof %d doesn't exist", proofID)
		}
		return err
	}

	err = k.Proofs.Remove(ctx, proof.Id)
	if err != nil {
		return err
	}

	err = k.ActiveProofsQueue.Remove(ctx, collections.Join(*proof.EndTime, proof.Id))
	if err != nil {
		return err
	}

	err = k.TheoremProof.Remove(ctx, proof.TheoremId)
	if err != nil {
		return err
	}

	err = k.ProofsByTheorem.Remove(ctx, collections.Join(proof.TheoremId, proof.Id))
	if err != nil {
		return err
	}

	return nil
}
