package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// ClaimLock locks collaterals after a claim proposal is submitted.
func (k Keeper) ClaimLock(ctx sdk.Context, proposalID, poolID uint64, purchaser sdk.AccAddress, purchaseID uint64, loss sdk.Coins, lockPeriod time.Duration) error {
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
	totalWithdrawing := k.GetTotalWithdrawing(ctx)
	totalLocked := k.GetTotalLocked(ctx)
	// No need to check if TotalCollateral - TotalWithdrawing > lossAmt + TotalLocked(ReservedCollaterals)
	if totalLocked.Add(lossAmt).GT(totalCollateral.Sub(totalWithdrawing)) {
		// (1) Compute collaterals NOT expiring within 4 days.
		impendingWithdrawAmount := k.ComputeWithdrawAmount(ctx, ctx.BlockHeader().Time.Add(k.gk.GetVotingParams(ctx).VotingPeriod * 2))
		lockableCollateral := totalCollateral.Sub(impendingWithdrawAmount)

		// (2) Verify that amount from (1) >= lossAmt + TotalLocked
		if lossAmt.Add(totalLocked).GT(lockableCollateral) {
			return types.ErrNotEnoughCollateral
		}
	}
	
	// TODO: Update purchase
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

	purchaseDeleteTime := purchase.ProtectionEndTime.Add(k.GetPurchaseDeletionPeriod(ctx))
	k.DequeuePurchase(ctx, purchaseList, purchaseDeleteTime)	
	purchase.Shield = purchase.Shield.Sub(lossAmt)
	votingEndTime := ctx.BlockTime().Add(lockPeriod)
	if purchaseDeleteTime.Before(votingEndTime) {
		// TODO: update delete time & protection end time?
	}
	k.SetPurchaseList(ctx, purchaseList)
	k.InsertPurchaseQueue(ctx, purchaseList, purchaseDeleteTime)
	
	// Update the pool.
	pool.Shield = pool.Shield.Sub(lossAmt)
	k.SetPool(ctx, pool)

	// Update global pool.
	totalShield := k.GetTotalShield(ctx)

	totalShield = totalShield.Sub(lossAmt)
	totalLocked = totalLocked.Sub(lossAmt)
	totalCollateral = totalCollateral.Sub(lossAmt)
	k.SetTotalShield(ctx, totalShield)
	k.SetTotalLocked(ctx, totalLocked)
	k.SetTotalCollateral(ctx, totalCollateral)

	return nil
}

func (k Keeper) ClaimUnlock(ctx sdk.Context, proposalID, poolID uint64, loss sdk.Coins) error {
 	lossAmt := loss.AmountOf(k.sk.BondDenom(ctx))
	 
	totalCollateral := k.GetTotalCollateral(ctx)
	totalLocked := k.GetTotalLocked(ctx)
	
	totalCollateral = totalCollateral.Sub(lossAmt)
	totalLocked = totalLocked.Sub(lossAmt)
	k.SetTotalCollateral(ctx, totalCollateral)
	k.SetTotalLocked(ctx, totalLocked)

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
