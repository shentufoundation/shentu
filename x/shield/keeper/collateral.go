package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// DepositCollateral deposits a community member's collateral for a pool.
func (k Keeper) DepositCollateral(ctx sdk.Context, from sdk.AccAddress, amount sdk.Int) error {
	// Check eligibility.
	provider, found := k.GetProvider(ctx, from)
	if !found {
		provider = k.addProvider(ctx, from)
	}
	provider.Collateral = provider.Collateral.Add(amount)
	if amount.GT(provider.Available) {
		return types.ErrInsufficientStaking
	}

	// Update states.
	provider.Available = provider.Available.Sub(amount)
	k.SetProvider(ctx, from, provider)

	totalCollateral := k.GetTotalCollateral(ctx)
	totalCollateral = totalCollateral.Add(amount)
	k.SetTotalCollateral(ctx, totalCollateral)

	return nil
}

// WithdrawCollateral withdraws a community member's collateral for a pool.
func (k Keeper) WithdrawCollateral(ctx sdk.Context, from sdk.AccAddress, amount sdk.Int) error {
	if amount.IsZero() {
		return nil
	}

	provider, found := k.GetProvider(ctx, from)
	if !found {
		return types.ErrProviderNotFound
	}
	withdrawable := provider.Collateral.Sub(provider.Withdrawing)
	if amount.GT(withdrawable) {
		return types.ErrOverWithdraw
	}

	// Insert into withdraw queue.
	poolParams := k.GetPoolParams(ctx)
	completionTime := ctx.BlockHeader().Time.Add(poolParams.WithdrawPeriod)
	withdraw := types.NewWithdraw(from, amount, completionTime)
	k.InsertWithdrawQueue(ctx, withdraw)

	provider.Withdrawing = provider.Withdrawing.Add(amount)
	k.SetProvider(ctx, provider.Address, provider)

	return nil
}
