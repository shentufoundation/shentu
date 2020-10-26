package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// ClaimLock locks collaterals after a claim proposal is submitted.
func (k Keeper) ClaimLock(ctx sdk.Context, proposalID, poolID uint64, purchaser sdk.AccAddress, purchaseID uint64, loss sdk.Coins, lockPeriod time.Duration) error {
	// Provider.Locked unnecessary?
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
	//totalWithdrawing := k.GetTotalWithdrawing(ctx) //TODO: No need for this variable?
	totalLocked := k.GetTotalLocked(ctx)
	newLockAmt := totalLocked.Add(lossAmt)
	// Ensure that total collateral (withdrawing and non-withdrawing)
	// can cover the new total lock amount.
	if newLockAmt.GT(totalCollateral) {
		return types.ErrNotEnoughCollateral
	}

	// Lock proportional amount from each provider.
	providers := k.GetAllProviders(ctx)
	proportion := newLockAmt.ToDec().Quo(totalCollateral.ToDec())
	remaining := lossAmt
	for i := range providers {
		var lockAmt sdk.Int
		if i < len(providers)-1 {
			lockAmt = providers[i].Collateral.ToDec().Mul(proportion).TruncateInt()
			if lockAmt.LT(providers[i].Collateral) && lockAmt.LT(remaining) {
				lockAmt = lockAmt.Add(sdk.OneInt())
			}
			remaining = remaining.Sub(lockAmt)
		} else {
			lockAmt = remaining
		}
		k.LockProvider(ctx, providers[i], types.NewLockedCollateral(proposalID, lockAmt), lockPeriod)
	}

	// Update shield amount and delete time of the purchase.
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
	k.DequeuePurchase(ctx, purchaseList, purchase.DeletionTime)
	purchase.Shield = purchase.Shield.Sub(lossAmt)
	votingEndTime := ctx.BlockTime().Add(lockPeriod)
	if purchase.DeletionTime.Before(votingEndTime) {
		// TODO: correctly update delete time & protection end time
		purchase.DeletionTime = votingEndTime // temp
	}
	k.SetPurchaseList(ctx, purchaseList)
	k.InsertExpiringPurchaseQueue(ctx, purchaseList, purchase.DeletionTime)

	// Update pool and global pool states
	pool.Shield = pool.Shield.Sub(lossAmt)
	k.SetPool(ctx, pool)

	totalLocked = newLockAmt // totalLocked.Add(lossAmt)
	totalCollateral = totalCollateral.Sub(lossAmt)
	totalShield := k.GetTotalShield(ctx)
	totalShield = totalShield.Sub(lossAmt)
	k.SetTotalLocked(ctx, totalLocked)
	k.SetTotalCollateral(ctx, totalCollateral)
	k.SetTotalShield(ctx, totalShield)

	return nil
}

// LockProvider locks the specified amount of collaterals from a provider.
// If necessary, it extends withdrawing collaterals and, if exist, their
// linked unbondings as well.
func (k Keeper) LockProvider(ctx sdk.Context, provider types.Provider, newLock types.LockedCollateral, lockPeriod time.Duration) {
	provider.Locked = provider.Locked.Add(newLock.Amount)
	provider.Collateral = provider.Collateral.Sub(newLock.Amount)
	provider.LockedCollaterals = append(provider.LockedCollaterals, newLock)

	// If there are enough bonded delegations backing 
	// locked collaterals, we are done.
	if provider.DelegationBonded.GTE(provider.Locked) {
		k.SetProvider(ctx, provider.Address, provider)
		return
	}

	// Lenient check:
	// Consider the amount of non-withdrawing collateral.
	if provider.Locked.GT(provider.Collateral.Sub(provider.Withdrawing)) {
		// Stricter check:
		// Consider the amount of non-withdrawing and withdrawing
		// collaterals 4 days from now.
		endTime := ctx.BlockTime().Add(lockPeriod)
		impendingWithdrawAmount := k.ComputeWithdrawAmountByTime(ctx, endTime)
		applicableCollateralAmt := provider.Collateral.Sub(impendingWithdrawAmount)
		if provider.Locked.GT(applicableCollateralAmt) {
			// Delay some withdrawals among the ones expiring within 4 days.
			amountToDelay := provider.Locked.Sub(applicableCollateralAmt)
			k.DelayWithdraws(ctx, lockPeriod, amountToDelay, provider.Address)
		}
	}
	k.SetProvider(ctx, provider.Address, provider)
}

func (k Keeper) ClaimUnlock(ctx sdk.Context, proposalID, poolID uint64, loss sdk.Coins) error {
	lossAmt := loss.AmountOf(k.sk.BondDenom(ctx))

	totalCollateral := k.GetTotalCollateral(ctx)
	totalLocked := k.GetTotalLocked(ctx)

	totalCollateral = totalCollateral.Add(lossAmt)
	totalLocked = totalLocked.Sub(lossAmt)
	k.SetTotalCollateral(ctx, totalCollateral)
	k.SetTotalLocked(ctx, totalLocked)

	// update providers
	providers := k.GetAllProviders(ctx)
	for i := range providers {
		for j := range providers[i].LockedCollaterals {
			provider := providers[i]
			lockedCollateral := providers[i].LockedCollaterals[j]

			if lockedCollateral.ProposalID == proposalID {
				provider.Locked = provider.Locked.Sub(lockedCollateral.Amount)
				provider.Collateral = provider.Collateral.Add(lockedCollateral.Amount)
				provider.LockedCollaterals = append(provider.LockedCollaterals[:j], provider.LockedCollaterals[j+1:]...)
				k.SetProvider(ctx, provider.Address, provider)
				break
			}
		}
	} // for each provider

	return nil
}

func (k Keeper) RestoreShield(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, id uint64, loss sdk.Coins) error {
	lossAmt := loss.AmountOf(k.sk.BondDenom(ctx))

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
