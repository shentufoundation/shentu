package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

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

	// Check coins to be deposited match the theorem's deposit params
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	if err := k.validateDepositDenom(ctx, params, grantAmount); err != nil {
		return err
	}

	if err := k.validateMinGrant(ctx, params, grantAmount); err != nil {
		return err
	}

	// update the bounty module's account coins pool
	if err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, grantor, types.ModuleName, grantAmount); err != nil {
		return err
	}

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
		// deposit exists
		grant.Amount = sdk.NewCoins(grant.Amount...).Add(grantAmount...)
	case errors.IsOf(err, collections.ErrNotFound):
		// deposit doesn't exist
		grant = types.NewGrant(theoremID, grantor, grantAmount)
	default:
		// failed to get deposit
		return err
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTheoremGrant,
			sdk.NewAttribute(types.AttributeKeyTheoremGrantor, grantor.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, grantAmount.String()),
			sdk.NewAttribute(types.AttributeKeyTheoremID, fmt.Sprintf("%d", theoremID)),
		),
	)

	return k.SetGrant(ctx, grant)
}

// RefundAndDeleteGrants refunds and deletes all the deposits on a timeout theorem.
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

// IterateGrants iterates over all the theorems deposits and performs a callback function
func (k Keeper) IterateGrants(ctx context.Context, theoremID uint64, cb func(key collections.Pair[uint64, sdk.AccAddress], value types.Grant) (bool, error)) error {
	rng := collections.NewPrefixedPairRange[uint64, sdk.AccAddress](theoremID)
	err := k.Grants.Walk(ctx, rng, cb)
	if err != nil {
		return err
	}
	return nil
}

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

	return nil
}

// validateMinGrant validates if initial grant is greater than or equal to the minimum
// required at the time of theorem submission. Returns nil on success, error otherwise.
func (k Keeper) validateMinGrant(ctx context.Context, params types.Params, initialDeposit sdk.Coins) error {
	// Check if the initial deposit is valid and has no negative amount
	if !initialDeposit.IsValid() || initialDeposit.IsAnyNegative() {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, initialDeposit.String())
	}

	// Check if the initial deposit meets the minimum required grant amount
	if !initialDeposit.IsAllGTE(params.MinGrant) {
		return errors.Wrapf(types.ErrMinGrantTooSmall, "was (%s), need (%s)", initialDeposit, params.MinGrant)
	}
	return nil
}

// validateDepositDenom validates if the deposit denom is accepted.
func (k Keeper) validateDepositDenom(ctx context.Context, params types.Params, depositAmount sdk.Coins) error {
	// Create a map for accepted denoms for quick lookup
	acceptedDenoms := make(map[string]bool, len(params.MinGrant))
	for _, coin := range params.MinGrant {
		acceptedDenoms[coin.Denom] = true
	}

	// Check if the deposited coins have valid denoms
	for _, coin := range depositAmount {
		if _, ok := acceptedDenoms[coin.Denom]; !ok {
			// Build a slice of accepted denoms for error message
			denoms := make([]string, 0, len(acceptedDenoms))
			for denom := range acceptedDenoms {
				denoms = append(denoms, denom)
			}
			return errors.Wrapf(types.ErrInvalidDepositDenom, "deposited %s, but bounty accepts only the following denom(s): %v", depositAmount, denoms)
		}
	}

	return nil
}

func (k Keeper) SetGrant(ctx context.Context, grant types.Grant) error {
	grantor, err := k.authKeeper.AddressCodec().StringToBytes(grant.Grantor)
	if err != nil {
		return err
	}
	return k.Grants.Set(ctx, collections.Join(grant.TheoremId, sdk.AccAddress(grantor)), grant)
}
