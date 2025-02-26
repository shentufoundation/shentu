package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func (k Keeper) AddGrant(ctx context.Context, theoremID uint64, grantor sdk.AccAddress, grantAmount sdk.Coins) error {
	// Check if theorem exists
	theorem, err := k.Theorems.Get(ctx, theoremID)
	if err != nil {
		return err
	}
	// Check theorem is still depositable
	if theorem.Status != types.TheoremStatus_THEOREM_STATUS_GRANT_PERIOD &&
		theorem.Status != types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD {
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
	err = k.SetTheorem(ctx, theorem)
	if err != nil {
		return err
	}

	// Add or update grant object
	grant, err := k.Grants.Get(ctx, collections.Join(grantor, theoremID))
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

	// TODO add event
	return k.SetGrant(ctx, grant)
}

func (k Keeper) SetGrant(ctx context.Context, grant types.Grant) error {
	grantor, err := k.authKeeper.AddressCodec().StringToBytes(grant.Grantor)
	if err != nil {
		return err
	}
	return k.Grants.Set(ctx, collections.Join(sdk.AccAddress(grantor), grant.TheoremId), grant)
}

func (k Keeper) SetDeposit(ctx context.Context, deposit types.Deposit) error {
	depositor, err := k.authKeeper.AddressCodec().StringToBytes(deposit.Depositor)
	if err != nil {
		return err
	}
	return k.Deposits.Set(ctx, collections.Join(sdk.AccAddress(depositor), deposit.ProofId), deposit)
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

// validateMinDeposit validates if deposit is greater than or equal to the minimum
// required at the time of proof submission. Returns nil on success, error otherwise.
func (k Keeper) validateMinDeposit(ctx context.Context, params types.Params, initialDeposit sdk.Coins) error {
	// Check if the initial deposit is valid and has no negative amount
	if !initialDeposit.IsValid() || initialDeposit.IsAnyNegative() {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, initialDeposit.String())
	}

	// Check if the initial deposit meets the minimum required grant amount
	if !initialDeposit.IsAllGTE(params.MinDeposit) {
		return errors.Wrapf(types.ErrMinDepositTooSmall, "was (%s), need (%s)", initialDeposit, params.MinGrant)
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
