package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

// ============================== Grant Operations ==============================

// AddGrant adds or updates a grant for a theorem
func (k Keeper) AddGrant(ctx context.Context, theoremID uint64, grantor sdk.AccAddress, grantAmount sdk.Coins) error {
	// Check if theorem exists and verify status
	theorem, err := k.Theorems.Get(ctx, theoremID)
	if err != nil {
		return err
	}
	if theorem.Status != types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD {
		return errors.Wrapf(types.ErrTheoremProposal, "%d", theoremID)
	}

	// Transfer funds to module account
	if err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, grantor, types.ModuleName, grantAmount); err != nil {
		return err
	}

	// Update theorem total grant
	theorem.TotalGrant = sdk.NewCoins(theorem.TotalGrant...).Add(grantAmount...)
	if err = k.Theorems.Set(ctx, theorem.Id, theorem); err != nil {
		return err
	}

	// Update or create grant record
	if err = k.updateOrCreateGrant(ctx, theoremID, grantor, grantAmount); err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeGrantTheorem,
			sdk.NewAttribute(types.AttributeKeyTheoremGrantor, grantor.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, grantAmount.String()),
			sdk.NewAttribute(types.AttributeKeyTheoremID, fmt.Sprintf("%d", theoremID)),
		),
	)

	return nil
}

// updateOrCreateGrant updates an existing grant or creates a new one
func (k Keeper) updateOrCreateGrant(ctx context.Context, theoremID uint64, grantor sdk.AccAddress, amount sdk.Coins) error {
	grant, err := k.Grants.Get(ctx, collections.Join(theoremID, grantor))
	switch {
	case err == nil:
		grant.Amount = sdk.NewCoins(grant.Amount...).Add(amount...)
	case errors.IsOf(err, collections.ErrNotFound):
		grant = types.NewGrant(theoremID, grantor, amount)
	default:
		return fmt.Errorf("failed to get grant: %w", err)
	}

	return k.SetGrant(ctx, grant)
}

// RefundAndDeleteGrants refunds and deletes all the grants for a theorem
func (k Keeper) RefundAndDeleteGrants(ctx context.Context, theoremID uint64) error {
	return k.IterateGrants(ctx, theoremID, func(key collections.Pair[uint64, sdk.AccAddress], grant types.Grant) (bool, error) {
		grantor := key.K2()
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, grantor, grant.Amount); err != nil {
			return false, err
		}
		return false, k.Grants.Remove(ctx, key)
	})
}

// IterateGrants iterates over all the theorem grants and performs a callback function
func (k Keeper) IterateGrants(ctx context.Context, theoremID uint64, cb func(key collections.Pair[uint64, sdk.AccAddress], value types.Grant) (bool, error)) error {
	rng := collections.NewPrefixedPairRange[uint64, sdk.AccAddress](theoremID)
	return k.Grants.Walk(ctx, rng, cb)
}

// DistributionGrants distributes rewards to checker and prover
func (k Keeper) DistributionGrants(ctx context.Context, theoremID uint64, checker, prover sdk.AccAddress) error {
	// Get parameters and theorem
	param, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	theorem, err := k.Theorems.Get(ctx, theoremID)
	if err != nil {
		return err
	}

	// Calculate rewards
	totalGrant := sdk.NewDecCoinsFromCoins(theorem.TotalGrant...)
	checkerRewards := totalGrant.MulDec(param.CheckerRate)
	proverRewards := totalGrant.Sub(checkerRewards)

	// Update rewards for checker and prover
	if err := k.updateReward(ctx, checker, checkerRewards); err != nil {
		return fmt.Errorf("failed to update checker reward: %w", err)
	}

	if err := k.updateReward(ctx, prover, proverRewards); err != nil {
		return fmt.Errorf("failed to update prover reward: %w", err)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDistributeReward,
			sdk.NewAttribute(types.AttributeKeyTheoremID, fmt.Sprintf("%d", theoremID)),
			sdk.NewAttribute(types.AttributeKeyChecker, checkerRewards.String()),
			sdk.NewAttribute(types.AttributeKeyProposer, proverRewards.String()),
		),
	)

	return nil
}

// updateReward updates the reward for a given address
func (k Keeper) updateReward(ctx context.Context, addr sdk.AccAddress, reward sdk.DecCoins) error {
	existingReward, err := k.Rewards.Get(ctx, addr)

	if err != nil && !errors.IsOf(err, collections.ErrNotFound) {
		return err
	}

	if errors.IsOf(err, collections.ErrNotFound) {
		existingReward = types.Reward{Reward: reward}
	} else {
		existingReward.Reward = existingReward.Reward.Add(reward...)
	}

	return k.Rewards.Set(ctx, addr, existingReward)
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
	// Check if proof exists and verify status
	proof, err := k.Proofs.Get(ctx, proofID)
	if err != nil {
		return err
	}

	if proof.Status != types.ProofStatus_PROOF_STATUS_HASH_LOCK_PERIOD {
		return types.ErrProofStatusInvalid
	}

	// Transfer funds and update deposit
	if err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, depositorAddr, types.ModuleName, depositAmount); err != nil {
		return err
	}

	if err = k.updateOrCreateDeposit(ctx, proofID, depositorAddr, depositAmount); err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDepositProof,
			sdk.NewAttribute(types.AttributeKeyProofDepositor, depositorAddr.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, depositAmount.String()),
			sdk.NewAttribute(types.AttributeKeyProofID, proofID),
		),
	)

	return nil
}

// updateOrCreateDeposit updates an existing deposit or creates a new one
func (k Keeper) updateOrCreateDeposit(ctx context.Context, proofID string, depositor sdk.AccAddress, amount sdk.Coins) error {
	deposit, err := k.Deposits.Get(ctx, collections.Join(proofID, depositor))
	switch {
	case err == nil:
		deposit.Amount = sdk.NewCoins(deposit.Amount...).Add(amount...)
	case errors.IsOf(err, collections.ErrNotFound):
		deposit = types.NewDeposit(proofID, depositor, amount)
	default:
		return fmt.Errorf("failed to get deposit: %w", err)
	}

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
	key := collections.Join(proofID, depositorAddr)
	deposit, err := k.Deposits.Get(ctx, key)
	if err != nil {
		return err
	}

	// Transfer funds back to depositor
	if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositorAddr, deposit.Amount); err != nil {
		return err
	}

	return k.Deposits.Remove(ctx, key)
}

// IterateDeposits iterates over all the deposits for a proof and performs a callback function
func (k Keeper) IterateDeposits(ctx context.Context, proofID string, cb func(key collections.Pair[string, sdk.AccAddress], value types.Deposit) (bool, error)) error {
	rng := collections.NewPrefixedPairRange[string, sdk.AccAddress](proofID)
	return k.Deposits.Walk(ctx, rng, cb)
}

// RefundAndDeleteDeposits refunds and deletes all the deposits for a proof
func (k Keeper) RefundAndDeleteDeposits(ctx context.Context, proofID string) error {
	return k.IterateDeposits(ctx, proofID, func(key collections.Pair[string, sdk.AccAddress], deposit types.Deposit) (bool, error) {
		depositor := key.K2()
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositor, deposit.Amount); err != nil {
			return false, err
		}
		return false, k.Deposits.Remove(ctx, key)
	})
}

// ============================== Funds Validation ==============================

// ValidateFunds validates funds amount and denomination against module parameters
// fundsType should be "grant" or "deposit" to determine which minimum amount to check against
func (k Keeper) ValidateFunds(ctx context.Context, amount sdk.Coins, fundsType string) (*types.Params, error) {
	// Get params
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get theorem parameters: %w", err)
	}

	// Basic validation
	if !amount.IsValid() || amount.IsAnyNegative() {
		return nil, errors.Wrap(sdkerrors.ErrInvalidCoins, amount.String())
	}

	// Determine minimum amount and error type based on funds type
	var minAmount sdk.Coins
	var errType error

	switch fundsType {
	case types.FundTypeGrant:
		minAmount = params.MinGrant
		errType = types.ErrMinGrantTooSmall
	case types.FundTypeDeposit:
		minAmount = params.MinDeposit
		errType = types.ErrMinDepositTooSmall
	default:
		return nil, fmt.Errorf("invalid funds type: %s", fundsType)
	}

	// Check minimum amount
	if !amount.IsAllGTE(minAmount) {
		return nil, errors.Wrapf(errType, "was (%s), need (%s)", amount, minAmount)
	}

	// Build accepted denominations map once for efficiency
	acceptedDenoms := make(map[string]bool, len(params.MinGrant))
	for _, coin := range params.MinGrant {
		acceptedDenoms[coin.Denom] = true
	}

	// Validate denominations
	invalidDenoms := make([]string, 0)
	for _, coin := range amount {
		if !acceptedDenoms[coin.Denom] {
			invalidDenoms = append(invalidDenoms, coin.Denom)
		}
	}

	if len(invalidDenoms) > 0 {
		// Build accepted denominations list for error message
		acceptedDenomsList := make([]string, 0, len(acceptedDenoms))
		for denom := range acceptedDenoms {
			acceptedDenomsList = append(acceptedDenomsList, denom)
		}
		return nil, errors.Wrapf(types.ErrInvalidDepositDenom,
			"deposited invalid denominations: %v, bounty accepts only the following denom(s): %v",
			invalidDenoms, acceptedDenomsList)
	}

	return &params, nil
}
