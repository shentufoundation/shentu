package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// SecureCollaterals is called after a claim is submitted to secure
// the given amount of collaterals for the duration and adjust shield
// module states accordingly.
func (k Keeper) SecureCollaterals(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, purchaseID uint64, loss sdk.Coins, duration time.Duration) error {
	lossAmt := loss.AmountOf(k.sk.BondDenom(ctx))

	// Verify shield.
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return types.ErrNoPoolFound
	}
	if lossAmt.GT(pool.Shield) {
		return types.ErrNotEnoughShield
	}

	// Verify collateral availability.
	totalCollateral := k.GetTotalCollateral(ctx)
	totalClaimed := k.GetTotalClaimed(ctx)
	totalClaimed = totalClaimed.Add(lossAmt)
	if totalClaimed.GT(totalCollateral) {
		panic("total claimed surpassed total collateral")
	}

	// Secure the updated loss ratio from each provider to cover total claimed.
	providers := k.GetAllProviders(ctx)
	lossRatio := totalClaimed.ToDec().Quo(totalCollateral.ToDec())
	for i := range providers {
		secureAmt := providers[i].Collateral.ToDec().Mul(lossRatio).TruncateInt()
		// Require each provider to secure one more unit, if possible,
		// so that the last provider does not have to cover combined 
		// truncated amounts.
		if secureAmt.LT(totalClaimed) && secureAmt.LT(providers[i].Collateral) {
			secureAmt = secureAmt.Add(sdk.OneInt())
		}
		k.SecureFromProvider(ctx, providers[i], secureAmt, duration)
		totalClaimed = totalClaimed.Sub(secureAmt)
	}

	// Update purchase states.
	purchaseList, found := k.GetPurchaseList(ctx, poolID, purchaser)
	if !found {
		return types.ErrPurchaseNotFound
	}
	var index int
	for i, entry := range purchaseList.Entries {
		if entry.PurchaseID == purchaseID {
			index = i
			break
		}
	}
	purchase := &purchaseList.Entries[index]
	if lossAmt.GT(purchase.Shield) {
		return types.ErrNotEnoughShield
	}
	purchase.Shield = purchase.Shield.Sub(lossAmt)
	votingEndTime := ctx.BlockTime().Add(duration)
	if purchase.DeletionTime.Before(votingEndTime) {
		purchase.DeletionTime = votingEndTime
	}
	k.SetPurchaseList(ctx, purchaseList)

	// Update pool and global pool states.
	pool.Shield = pool.Shield.Sub(lossAmt)
	k.SetPool(ctx, pool)

	totalShield := k.GetTotalShield(ctx)
	totalShield = totalShield.Sub(lossAmt)
	k.SetTotalShield(ctx, totalShield)
	k.SetTotalClaimed(ctx, totalClaimed)

	return nil
}

// SecureFromProvider secures the specified amount of collaterals from
// the provider for the duration. If necessary, it extends withdrawing
// collaterals and, if exist, their linked unbondings as well.
func (k Keeper) SecureFromProvider(ctx sdk.Context, provider types.Provider, amount sdk.Int, duration time.Duration) {	
	// Lenient check:
	// Check if non-withdrawing, bonded delegation-backed collaterals
	// cannot cover the amount.
	if provider.Collateral.Sub(provider.Withdrawing).GTE(amount) && provider.DelegationBonded.GTE(provider.DelegationBonded) {
		return
	}

	// Secure the given amount of collaterals until the end of the
	// lock period by delaying withdrawals, if necessary.
	endTime := ctx.BlockTime().Add(duration)
	upcomingWithdrawAmount := k.ComputeWithdrawAmountByTime(ctx, provider.Address, endTime)
	notWithdrawnSoon := provider.Collateral.Sub(upcomingWithdrawAmount)

	// Secure the given amount of staking (bonded or unbonding) until
	// the end of the lock period by delaying unbondings, if necessary.
	totalUnbondingAmount := k.ComputeTotalUnbondingAmount(ctx, provider.Address)
	upcomingUnbondingAmount := k.ComputeUnbondingAmountByTime(ctx, provider.Address, endTime)
	notUnbondedSoon := provider.DelegationBonded.Add(totalUnbondingAmount.Sub(upcomingUnbondingAmount))

	if amount.GT(notWithdrawnSoon) || amount.GT(notUnbondedSoon) {
		withdrawDelayAmt := amount.Sub(notWithdrawnSoon)
		k.DelayWithdraws(ctx, provider.Address, withdrawDelayAmt, duration)

		unbondingDelayAmt := amount.Sub(notUnbondedSoon)
		k.DelayUnbonding(ctx, provider.Address, unbondingDelayAmt, duration)
	}
}

func (k Keeper) ClaimEnd(ctx sdk.Context, id, poolID uint64, loss sdk.Coins) {
	lossAmt := loss.AmountOf(k.sk.BondDenom(ctx))
	totalClaimed := k.GetTotalClaimed(ctx)
	totalClaimed = totalClaimed.Sub(lossAmt)
	k.SetTotalClaimed(ctx, totalClaimed)
}

func (k Keeper) RestoreShield(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, id uint64, loss sdk.Coins) error {
	lossAmt := loss.AmountOf(k.sk.BondDenom(ctx))

	// Update the total shield.
	totalShield := k.GetTotalShield(ctx)
	totalShield = totalShield.Add(lossAmt)
	k.SetTotalShield(ctx, totalShield)

	// Update shield of the pool.
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return types.ErrNoPoolFound
	}
	pool.Shield = pool.Shield.Add(lossAmt)
	k.SetPool(ctx, pool)

	// Update shield of the purchase.
	purchaseList, found := k.GetPurchaseList(ctx, poolID, purchaser)
	if !found {
		return types.ErrPurchaseNotFound
	}
	for i := range purchaseList.Entries {
		if purchaseList.Entries[i].PurchaseID == id {
			purchaseList.Entries[i].Shield = purchaseList.Entries[i].Shield.Add(lossAmt)
			break
		}
	}
	k.SetPurchaseList(ctx, purchaseList)

	return nil
}
