package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

// ============================== Grant Operations ==============================

// AddGrant adds or updates a grant for a theorem
func (k Keeper) AddGrant(ctx context.Context, theoremID uint64, grantor sdk.AccAddress, grantAmount sdk.Coins) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	// Check if theorem exists
	theorem, err := k.Theorems.Get(ctx, theoremID)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return status.Errorf(codes.NotFound, "theorem %d doesn't exist", theoremID)
		}
		return err
	}
	// Check theorem is still depositable
	if theorem.Status != types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD {
		return errors.Wrapf(types.ErrTheoremProposal, "%d", theoremID)
	}

	// validate deposit parameters - this will be handled by the caller (usually msgServer)

	// Update theorem
	theorem.TotalGrant = sdk.NewCoins(theorem.TotalGrant...).Add(grantAmount...)
	err = k.Theorems.Set(ctx, theorem.Id, theorem)
	if err != nil {
		return err
	}

	// Add or update grant object
	grant, err := k.Grants.Get(ctx, collections.Join(theoremID, grantor))
	switch {
	case err == nil:
		// grant exists
		grant.Amount = sdk.NewCoins(grant.Amount...).Add(grantAmount...)
	case errors.IsOf(err, collections.ErrNotFound):
		// grant doesn't exist
		grant = types.NewGrant(theoremID, grantor, grantAmount)
	default:
		// failed to get grant
		return err
	}

	// update the bounty module's account coins pool
	if err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, grantor, types.ModuleName, grantAmount); err != nil {
		return err
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeGrantTheorem,
			sdk.NewAttribute(types.AttributeKeyTheoremGrantor, grantor.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, grantAmount.String()),
			sdk.NewAttribute(types.AttributeKeyTheoremID, fmt.Sprintf("%d", theoremID)),
		),
	)

	return k.SetGrant(ctx, grant)
}

// RefundAndDeleteGrants refunds and deletes all the grants on a timeout theorem
func (k Keeper) RefundAndDeleteGrants(ctx context.Context, theoremID uint64) error {
	return k.IterateGrants(ctx, theoremID, func(key collections.Pair[uint64, sdk.AccAddress], grant types.Grant) (bool, error) {
		grantor := key.K2()
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, grantor, grant.Amount)
		if err != nil {
			return false, err
		}
		err = k.Grants.Remove(ctx, key)
		return false, err
	})
}

// IterateGrants iterates over all the theorem grants and performs a callback function
func (k Keeper) IterateGrants(ctx context.Context, theoremID uint64, cb func(key collections.Pair[uint64, sdk.AccAddress], value types.Grant) (bool, error)) error {
	rng := collections.NewPrefixedPairRange[uint64, sdk.AccAddress](theoremID)
	err := k.Grants.Walk(ctx, rng, cb)
	if err != nil {
		return err
	}
	return nil
}

// DistributionGrants distributes rewards to checker and prover
func (k Keeper) DistributionGrants(ctx context.Context, theoremID uint64, checker, prover sdk.AccAddress) error {
	param, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	theorem, err := k.Theorems.Get(ctx, theoremID)
	if err != nil {
		return err
	}
	totalGrant := sdk.NewDecCoinsFromCoins(theorem.TotalGrant...)

	cReward := totalGrant.MulDec(param.CheckerRate)
	pReward := totalGrant.Sub(cReward)

	checkerReward, err := k.Rewards.Get(ctx, checker)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			// not found
			checkerReward.Reward = cReward
		} else {
			return err
		}
	} else {
		checkerReward.Reward = checkerReward.Reward.Add(cReward...)
	}

	err = k.Rewards.Set(ctx, checker, checkerReward)
	if err != nil {
		return err
	}

	proverReward, err := k.Rewards.Get(ctx, prover)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			// not found
			proverReward.Reward = pReward
		} else {
			return err
		}
	} else {
		proverReward.Reward = proverReward.Reward.Add(pReward...)
	}
	err = k.Rewards.Set(ctx, prover, proverReward)
	if err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDistributeReward,
			sdk.NewAttribute(types.AttributeKeyTheoremID, fmt.Sprintf("%d", theoremID)),
			sdk.NewAttribute(types.AttributeKeyChecker, cReward.String()),
			sdk.NewAttribute(types.AttributeKeyProposer, pReward.String()),
		),
	)
	return nil
}

// SetGrant sets a grant in the store
func (k Keeper) SetGrant(ctx context.Context, grant types.Grant) error {
	grantor, err := k.authKeeper.AddressCodec().StringToBytes(grant.Grantor)
	if err != nil {
		return err
	}
	return k.Grants.Set(ctx, collections.Join(grant.TheoremId, sdk.AccAddress(grantor)), grant)
}

// ============================== Deposit Operations ==============================

// AddDeposit adds or updates a deposit for a proof
func (k Keeper) AddDeposit(ctx context.Context, proofID string, depositorAddr sdk.AccAddress, depositAmount sdk.Coins) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	// Check if proof exists
	proof, err := k.Proofs.Get(ctx, proofID)
	if err != nil {
		return err
	}
	// Check proof status
	if proof.Status != types.ProofStatus_PROOF_STATUS_HASH_LOCK_PERIOD {
		return types.ErrProofStatusInvalid
	}

	// validate deposit parameters - this will be handled by the caller (usually msgServer)

	// update the bounty module's account coins pool
	if err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, depositorAddr, types.ModuleName, depositAmount); err != nil {
		return err
	}

	// Add or update deposit object
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

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDepositProof,
			sdk.NewAttribute(types.AttributeKeyProofDepositor, depositorAddr.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, depositAmount.String()),
			sdk.NewAttribute(types.AttributeKeyProofID, proofID),
		),
	)

	return k.SetDeposit(ctx, deposit)
}

// SetDeposit sets a deposit in the store
func (k Keeper) SetDeposit(ctx context.Context, deposit types.Deposit) error {
	depositor, err := k.authKeeper.AddressCodec().StringToBytes(deposit.Depositor)
	if err != nil {
		return err
	}
	return k.Deposits.Set(ctx, collections.Join(deposit.ProofId, sdk.AccAddress(depositor)), deposit)
}

// RefundAndDeleteDeposit refunds and deletes a specific deposit
func (k Keeper) RefundAndDeleteDeposit(ctx context.Context, proofID string, depositorAddr sdk.AccAddress) error {
	deposit, err := k.Deposits.Get(ctx, collections.Join(proofID, depositorAddr))
	if err != nil {
		return err
	}

	// refund the deposit amount to the depositor
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositorAddr, deposit.Amount)
	if err != nil {
		return err
	}

	// remove the deposit from storage
	return k.Deposits.Remove(ctx, collections.Join(proofID, depositorAddr))
}

// IterateDeposits iterates over all the deposits for a proof and performs a callback function
func (k Keeper) IterateDeposits(ctx context.Context, proofID string, cb func(key collections.Pair[string, sdk.AccAddress], value types.Deposit) (bool, error)) error {
	rng := collections.NewPrefixedPairRange[string, sdk.AccAddress](proofID)
	err := k.Deposits.Walk(ctx, rng, cb)
	if err != nil {
		return err
	}
	return nil
}

// RefundAndDeleteDeposits refunds and deletes all the deposits for a proof
func (k Keeper) RefundAndDeleteDeposits(ctx context.Context, proofID string) error {
	return k.IterateDeposits(ctx, proofID, func(key collections.Pair[string, sdk.AccAddress], deposit types.Deposit) (bool, error) {
		depositor := key.K2()
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositor, deposit.Amount)
		if err != nil {
			return false, err
		}
		err = k.Deposits.Remove(ctx, key)
		return false, err
	})
}

// ============================== Funds Validation ==============================

// ValidateFunds validates funds amount and denomination against module parameters
// fundsType should be "grant" or "deposit" to determine which minimum amount to check against
func (k Keeper) ValidateFunds(ctx context.Context, amount sdk.Coins, fundsType string) (*types.Params, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params, err := k.Params.Get(sdkCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to get theorem parameters: %w", err)
	}

	// Check if the amount is valid and has no negative amount
	if !amount.IsValid() || amount.IsAnyNegative() {
		return nil, errors.Wrap(sdkerrors.ErrInvalidCoins, amount.String())
	}

	// Check against minimum amount based on funds type
	var minAmount sdk.Coins
	var errType error

	if fundsType == "grant" {
		minAmount = params.MinGrant
		errType = types.ErrMinGrantTooSmall
	} else { // deposit
		minAmount = params.MinDeposit
		errType = types.ErrMinDepositTooSmall
	}

	if !amount.IsAllGTE(minAmount) {
		return nil, errors.Wrapf(errType, "was (%s), need (%s)", amount, minAmount)
	}

	// Validate deposit denomination
	// Create a map for accepted denoms for quick lookup
	acceptedDenoms := make(map[string]bool, len(params.MinGrant))
	for _, coin := range params.MinGrant {
		acceptedDenoms[coin.Denom] = true
	}

	// Check if the deposited coins have valid denoms
	for _, coin := range amount {
		if _, ok := acceptedDenoms[coin.Denom]; !ok {
			// Build a slice of accepted denoms for error message
			denoms := make([]string, 0, len(acceptedDenoms))
			for denom := range acceptedDenoms {
				denoms = append(denoms, denom)
			}
			return nil, errors.Wrapf(types.ErrInvalidDepositDenom, "deposited %s, but bounty accepts only the following denom(s): %v", amount, denoms)
		}
	}

	return &params, nil
}
