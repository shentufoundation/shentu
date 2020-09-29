package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// ClaimLock locks collaterals after a claim proposal is submitted.
func (k Keeper) ClaimLock(ctx sdk.Context, poolID uint64, loss sdk.Coins, purchaseTxHash string, proposalID uint64) error {
	pool, err := k.GetPool(ctx, poolID)
	if err != nil {
		return err
	}

	// update shield of purchase
	purchase, err := k.GetPurchase(ctx, purchaseTxHash)
	if err != nil {
		return err
	}
	purchase.Shield = purchase.Shield.Sub(loss)
	k.SetPurchase(ctx, purchaseTxHash, purchase)

	// update total collateral and shield of pool
	pool.TotalCollateral = pool.TotalCollateral.Sub(loss)
	pool.Shield = pool.Shield.Sub(loss)

	// update locked collaterals for community
	for i, collateral := range pool.Community {
		lockedCoins := GetLockedCoins(loss, pool.TotalCollateral, collateral.Amount)
		lockedCollateral := types.NewLockedCollateral(proposalID, lockedCoins)
		pool.Community[i].LockedCollaterals = append(pool.Community[i].LockedCollaterals, lockedCollateral)
		pool.Community[i].Amount = pool.Community[i].Amount.Sub(lockedCoins)
	}

	// update locked collateral for CertiK
	lockedCoins := GetLockedCoins(loss, pool.TotalCollateral, pool.CertiK.Amount)
	lockedCollateral := types.NewLockedCollateral(proposalID, lockedCoins)
	pool.CertiK.LockedCollaterals = append(pool.CertiK.LockedCollaterals, lockedCollateral)
	pool.CertiK.Amount = pool.CertiK.Amount.Sub(lockedCoins)

	k.SetPool(ctx, pool)
	return nil
}

// ClaimUnlock unlocks locked collaterals.
func (k Keeper) ClaimUnlock(ctx sdk.Context, poolID uint64, loss sdk.Coins, proposalID uint64) error {
	pool, err := k.GetPool(ctx, poolID)
	if err != nil {
		return err
	}
	pool.TotalCollateral = pool.TotalCollateral.Add(loss...)

	// unlock collaterals for community
	for i, collateral := range pool.Community {
		for j, locked := range collateral.LockedCollaterals {
			if locked.ProposalID == proposalID {
				collateral.Amount = collateral.Amount.Add(locked.LockedCoins...)
				collateral.LockedCollaterals[j] = collateral.LockedCollaterals[len(collateral.LockedCollaterals)-1]
				collateral.LockedCollaterals = collateral.LockedCollaterals[:len(collateral.LockedCollaterals)-1]
				break
			}
		}
		pool.Community[i] = collateral
	}

	// unlock collaterals for CertiK
	c := pool.CertiK
	for i, locked := range c.LockedCollaterals {
		if locked.ProposalID == proposalID {
			c.Amount = c.Amount.Add(locked.LockedCoins...)
			c.LockedCollaterals[i] = c.LockedCollaterals[len(c.LockedCollaterals)-1]
			c.LockedCollaterals = c.LockedCollaterals[:len(c.LockedCollaterals)-1]
			break
		}
	}
	pool.CertiK = c

	k.SetPool(ctx, pool)
	return nil
}

// RestoreShield restores shield for proposer.
func (k Keeper) RestoreShield(ctx sdk.Context, poolID uint64, loss sdk.Coins, purchaseTxHash string) error {
	pool, err := k.GetPool(ctx, poolID)
	if err != nil {
		return err
	}

	// update shield of purchase
	purchase, err := k.GetPurchase(ctx, purchaseTxHash)
	if err != nil {
		return err
	}
	purchase.Shield = purchase.Shield.Sub(loss)
	k.SetPurchase(ctx, purchaseTxHash, purchase)

	// update shield of pool
	pool.Shield = pool.Shield.Add(loss...)
	k.SetPool(ctx, pool)

	return nil
}
