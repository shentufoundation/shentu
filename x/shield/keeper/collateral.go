package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/x/shield/types"
)

// DepositCollateral deposits a community member's collateral for a pool.
func (k Keeper) DepositCollateral(ctx sdk.Context, from sdk.AccAddress, amount sdk.Int) error {
	provider, found := k.GetProvider(ctx, from)
	if !found {
		provider = k.addProvider(ctx, from)
	}
	// Check if there are enough delegations backing collaterals.
	if ctx.BlockHeight() < common.Update1Height && provider.DelegationBonded.LT(provider.Collateral.Add(amount).Sub(provider.Withdrawing)) ||
		ctx.BlockHeight() >= common.Update1Height && provider.DelegationBonded.LT(provider.Collateral.Add(amount)) {
		return types.ErrInsufficientStaking
	}

	// Update provider.
	provider.Collateral = provider.Collateral.Add(amount)
	k.SetProvider(ctx, from, provider)

	// Update total collateral.
	totalCollateral := k.GetTotalCollateral(ctx)
	totalCollateral = totalCollateral.Add(amount)
	k.SetTotalCollateral(ctx, totalCollateral)

	return nil
}

// WithdrawCollateral withdraws a community member's collateral for a pool.
// In case of unbonding-initiated withdraw, store the validator address and
// the creation height.
func (k Keeper) WithdrawCollateral(ctx sdk.Context, from sdk.AccAddress, amount sdk.Int) error {
	if amount.IsZero() {
		return nil
	}

	// Check the collateral can be withdrew.
	provider, found := k.GetProvider(ctx, from)
	if !found {
		return types.ErrProviderNotFound
	}

	// Do not need to consider shield for withdrawable amount
	// because the withdraw period is the same as the shield protection period.
	withdrawable := provider.Collateral.Sub(provider.Withdrawing)
	if amount.GT(withdrawable) {
		return types.ErrOverWithdraw
	}

	// Insert into withdraw queue.
	poolParams := k.GetPoolParams(ctx)
	completionTime := ctx.BlockHeader().Time.Add(poolParams.WithdrawPeriod)
	withdraw := types.NewWithdraw(from, amount, completionTime)
	k.InsertWithdrawQueue(ctx, withdraw)

	// Update provider's withdrawing.
	provider.Withdrawing = provider.Withdrawing.Add(amount)
	k.SetProvider(ctx, provider.Address, provider)

	// Update total withdrawing.
	totalWithdrawing := k.GetTotalWithdrawing(ctx)
	totalWithdrawing = totalWithdrawing.Add(amount)
	k.SetTotalWithdrawing(ctx, totalWithdrawing)

	return nil
}
