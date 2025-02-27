package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func (k Keeper) SubmitProofHash(ctx context.Context, theoremID uint64, proofID, prover string, deposit sdk.Coins) (*types.Proof, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	submitTime := sdkCtx.BlockHeader().Time

	// Check if theorem exists
	theorem, err := k.Theorems.Get(ctx, theoremID)
	if err != nil {
		return nil, err
	}
	// Check theorem is still depositable
	if theorem.Status != types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD {
		return nil, types.ErrTheoremStatusInvalid
	}

	param, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	proof, err := types.NewProof(theoremID, proofID, prover, submitTime, submitTime.Add(*param.ProofHashLockPeriod), deposit)
	if err != nil {
		return nil, err
	}

	if err := k.SetProof(ctx, proof); err != nil {
		return nil, err
	}

	return &proof, nil
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
	proof.Detail = detail

	return k.SetProof(ctx, proof)
}

func (k Keeper) AddDeposit(ctx context.Context, proofID string, depositorAddr sdk.AccAddress, depositAmount sdk.Coins) error {
	// Check if proof exists
	proof, err := k.Proofs.Get(ctx, proofID)
	if err != nil {
		return err
	}
	// Check proof status
	if proof.Status != types.ProofStatus_PROOF_STATUS_HASH_LOCK_PERIOD {
		return types.ErrProofStatusInvalid
	}

	// Check coins to be deposited match the theorem's deposit params
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	if err := k.validateDepositDenom(ctx, params, depositAmount); err != nil {
		return err
	}

	if err := k.validateMinDeposit(ctx, params, depositAmount); err != nil {
		return err
	}

	// update the bounty module's account coins pool
	if err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, depositorAddr, types.ModuleName, depositAmount); err != nil {
		return err
	}

	// Add or update grant object
	deposit, err := k.Deposits.Get(ctx, collections.Join(proofID, depositorAddr))
	switch {
	case err == nil:
		// deposit exists
		deposit.Amount = sdk.NewCoins(deposit.Amount...).Add(depositAmount...)
	case errors.IsOf(err, collections.ErrNotFound):
		// deposit doesn't exist
		deposit = types.NewDeposit(proofID, depositorAddr, depositAmount)
	default:
		// failed to get deposit
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProofDeposit,
			sdk.NewAttribute(types.AttributeKeyProofDepositor, depositorAddr.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, depositAmount.String()),
			sdk.NewAttribute(types.AttributeKeyProofID, proofID),
		),
	)

	return k.SetDeposit(ctx, deposit)
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

func (k Keeper) SetProof(ctx context.Context, proof types.Proof) error {
	return k.Proofs.Set(ctx, proof.Id, proof)
}
