package keeper

import (
	"github.com/tendermint/tendermint/crypto/tmhash"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// SetPurchase sets a purchase of shield.
func (k Keeper) SetPurchase(ctx sdk.Context, txhash []byte, purchase types.Purchase) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(purchase)
	store.Set(types.GetPurchaseTxHashKey(txhash), bz)
}

// GetPurchase gets a purchase from store by txhash.
func (k Keeper) GetPurchase(ctx sdk.Context, txhash []byte) (types.Purchase, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPurchaseTxHashKey(txhash))
	if bz != nil {
		var purchase types.Purchase
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &purchase)
		return purchase, nil
	}
	return types.Purchase{}, types.ErrPurchaseNotFound
}

// DeletePurchase deletes a purchase of shield.
func (k Keeper) DeletePurchase(ctx sdk.Context, txhash []byte) error {
	store := ctx.KVStore(k.storeKey)
	_, err := k.GetPurchase(ctx, txhash)
	if err != nil {
		return err
	}
	store.Delete(types.GetPurchaseTxHashKey(txhash))
	return nil
}

// PurchaseShield purchases shield of a pool.
func (k Keeper) PurchaseShield(
	ctx sdk.Context, poolID uint64, shield sdk.Coins, description string, purchaser sdk.AccAddress,
) (types.Purchase, error) {
	pool, err := k.GetPool(ctx, poolID)
	if err != nil {
		return types.Purchase{}, err
	}
	poolParams := k.GetPoolParams(ctx)
	claimParams := k.GetClaimProposalParams(ctx)

	// check preconditions
	if !pool.Active {
		return types.Purchase{}, types.ErrPoolInactive
	}
	if pool.EndTime <= ctx.BlockTime().Unix() + types.DefaultWithdrawalPeriod.Milliseconds()/1000 {
		return types.Purchase{}, types.ErrPoolLifeTooShort
	}
	shieldAmt := shield.AmountOf(k.sk.BondDenom(ctx))
	if shieldAmt.GT(pool.Available) {
		return types.Purchase{}, types.ErrNotEnoughShield
	}

	// send tokens to shield module account
	shieldDec := sdk.NewDecCoinsFromCoins(shield...)
	premium, _ := shieldDec.MulDec(poolParams.ShieldFeesRate).TruncateDecimal()
	if err := k.DepositNativePremium(ctx, premium, purchaser); err != nil {
		return types.Purchase{}, err
	}

	// update pool premium, shield and available
	premiumMixedDec := types.NewMixedDecCoins(sdk.NewDecCoinsFromCoins(premium...), sdk.DecCoins{})
	pool.Premium = pool.Premium.Add(premiumMixedDec)
	pool.Shield = pool.Shield.Add(shield...)
	pool.Available = pool.Available.Sub(shieldAmt)
	k.SetPool(ctx, pool)

	// set purchase
	txhash := tmhash.Sum(ctx.TxBytes())
	protectionEndTime := ctx.BlockTime().Add(poolParams.ProtectionPeriod)
	claimPeriodEndTime := ctx.BlockTime().Add(claimParams.ClaimPeriod)
	purchase := types.NewPurchase(txhash, poolID, shield, ctx.BlockHeight(), protectionEndTime, claimPeriodEndTime, description, purchaser)
	k.SetPurchase(ctx, txhash, purchase)

	return purchase, nil
}

// IterateAllPurchases iterates over the all the stored purchases and performs a callback function.
func (k Keeper) IterateAllPurchases(ctx sdk.Context, callback func(purchase types.Purchase) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PurchaseKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var purchase types.Purchase
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &purchase)

		if callback(purchase) {
			break
		}
	}
}

// TODO improve the performance
// RemoveExpiredPurchases removes purchases whose claim period end time is before current block time.
func (k Keeper) RemoveExpiredPurchases(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PurchaseKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var purchase types.Purchase
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &purchase)
		if purchase.ClaimPeriodEndTime.Before(ctx.BlockTime()) {
			pool, err := k.GetPool(ctx, purchase.PoolID)
			if err == nil {
				pool.Available = pool.Available.Add(purchase.Shield.AmountOf(k.sk.BondDenom(ctx)))
				pool.Shield = pool.Shield.Sub(purchase.Shield)
				k.SetPool(ctx, pool)
			}
			store.Delete(iterator.Key())
		}
	}
}

// GetOnesPurchases returns a purchaser's all purchases.
func (k Keeper) GetOnesPurchases(ctx sdk.Context, address sdk.AccAddress) (purchases []types.Purchase) {
	k.IterateAllPurchases(ctx, func(purchase types.Purchase) bool {
		if purchase.Purchaser.Equals(address) {
			purchases = append(purchases, purchase)
		}
		return false
	})
	return purchases
}

// GetPoolPurchases returns a all purchases in a given pool.
func (k Keeper) GetPoolPurchases(ctx sdk.Context, poolID uint64) (purchases []types.Purchase) {
	k.IterateAllPurchases(ctx, func(purchase types.Purchase) bool {
		if purchase.PoolID == poolID {
			purchases = append(purchases, purchase)
		}
		return false
	})
	return purchases
}

// IteratePurchases iterates through purchases in a pool
func (k Keeper) IteratePurchases(ctx sdk.Context, callback func(purchase types.Purchase) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PurchaseKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var purchase types.Purchase
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &purchase)

		if callback(purchase) {
			break
		}
	}
}

// GetAllPurchases retrieves all purchases.
func (k Keeper) GetAllPurchases(ctx sdk.Context) (purchases []types.Purchase) {
	k.IteratePurchases(ctx, func(purchase types.Purchase) bool {
		purchases = append(purchases, purchase)
		return false
	})
	return
}
