package keeper

import (
	"encoding/hex"

	"github.com/tendermint/tendermint/crypto/tmhash"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// SetPurchase sets a purchase of shield.
func (k Keeper) SetPurchase(ctx sdk.Context, txhash string, purchase types.Purchase) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(purchase)
	store.Set(types.GetPurchaseTxHashKey(txhash), bz)
}

// GetPurchase gets a purchase from store by txhash.
func (k Keeper) GetPurchase(ctx sdk.Context, txhash string) (types.Purchase, error) {
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
func (k Keeper) DeletePurchase(ctx sdk.Context, txhash string) error {
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
	if !pool.Shield.Add(shield...).IsAllGTE(pool.TotalCollateral) {
		return types.Purchase{}, types.ErrNotEnoughShield
	}

	// send tokens to shield module account
	shieldDec := sdk.NewDecCoinsFromCoins(shield...)
	premium, _ := shieldDec.MulDec(poolParams.ShieldFeesRate).TruncateDecimal()
	if err := k.DepositNativePremium(ctx, premium, purchaser); err != nil {
		return types.Purchase{}, err
	}

	// update pool premium and shield
	premiumMixedDec := types.NewMixedDecCoins(sdk.NewDecCoinsFromCoins(premium...), sdk.DecCoins{})
	pool.Premium = pool.Premium.Add(premiumMixedDec)
	pool.Shield = pool.Shield.Add(shield...)
	k.SetPool(ctx, pool)

	// set purchase
	txhash := hex.EncodeToString(tmhash.Sum(ctx.TxBytes()))
	protectionEndTime := ctx.BlockTime().Add(poolParams.ProtectionPeriod)
	claimPeriodEndTime := ctx.BlockTime().Add(claimParams.ClaimPeriod)
	purchase := types.NewPurchase(poolID, shield, ctx.BlockHeight(), protectionEndTime, claimPeriodEndTime, description, purchaser)
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
			store.Delete(iterator.Key())
		}
	}
}
