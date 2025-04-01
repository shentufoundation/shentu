package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
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

	err = k.ProofsByTheorem.Remove(ctx, collections.Join(proof.TheoremId, proof.Id))
	if err != nil {
		return err
	}

	addrBytes, err := k.authKeeper.AddressCodec().StringToBytes(proof.Prover)
	if err != nil {
		return err
	}
	err = k.Deposits.Remove(ctx, collections.Join(proof.Id, sdk.AccAddress(addrBytes)))
	if err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDeleteProof,
			sdk.NewAttribute(types.AttributeKeyProofID, proofID),
		),
	)

	return nil
}
