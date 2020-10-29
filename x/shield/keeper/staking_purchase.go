package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

func (k Keeper) GetGlobalStakingPurchasePool(ctx sdk.Context) (pool types.GlobalStakingPool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetStakingPurchasePoolKey())
	if bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &pool)
	}
	return types.GlobalStakingPool{
		Amount: sdk.NewInt(0),
	}
}

func (k Keeper) SetGlobalStakingPurchasePool(ctx sdk.Context, pool types.GlobalStakingPool) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(pool)
	store.Set(types.GetStakingPurchasePoolKey(), bz)
}

func (k Keeper) GetStakingPurchase(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress) (purchase types.StakingPurchase, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetStakingPurchaseKey(poolID, purchaser))
	if bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &purchase)
		found = true
	}
	return
}

func (k Keeper) SetStakingPurchase(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, purchase types.StakingPurchase) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(purchase)
	store.Set(types.GetStakingPurchaseKey(poolID, purchaser), bz)
}

func (k Keeper) AddGlobalStakingPurchasePool(ctx sdk.Context, amount sdk.Int) {
	pool := k.GetGlobalStakingPurchasePool(ctx)
	pool.Amount = pool.Amount.Add(amount)
	k.SetGlobalStakingPurchasePool(ctx, pool)
}

func (k Keeper) SubGlobalStakingPurchasePool(ctx sdk.Context, amount sdk.Int) {
	pool := k.GetGlobalStakingPurchasePool(ctx)
	pool.Amount = pool.Amount.Sub(amount)
	k.SetGlobalStakingPurchasePool(ctx, pool)
}

func (k Keeper) AddStaking(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, amount sdk.Int, endTime time.Time) error {
	k.AddGlobalStakingPurchasePool(ctx, amount)
	sp, found := k.GetStakingPurchase(ctx, poolID, purchaser)
	if !found {
		sp = types.NewStakingPurchase(poolID, purchaser, amount)
	}
	sp.Locked = sp.Locked.Add(amount)
	sp.Amount = sp.Amount.Add(amount)
	newExpiration := types.NewStakingExpiration(endTime, amount)
	sp.Expirations = append(sp.Expirations, newExpiration)
	err := k.supplyKeeper.SendCoinsFromAccountToModule(
		ctx, purchaser, types.ModuleName, sdk.NewCoins(sdk.NewCoin(k.sk.BondDenom(ctx), amount)))
	if err != nil {
		return err
	}
	k.SetStakingPurchase(ctx, poolID, purchaser, sp)
	return nil
}

func (k Keeper) WithdrawStaking(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, amount sdk.Int) error {
	sp, found := k.GetStakingPurchase(ctx, poolID, purchaser)
	if !found {
		return types.ErrPurchaseNotFound
	}
	newAmt := sp.Amount.Sub(amount)
	if newAmt.LT(sp.Locked) {
		return types.ErrNotEnoughStaked
	}
	if newAmt.IsZero() {
		store := ctx.KVStore(k.storeKey)
		store.Delete(types.GetStakingPurchaseKey(poolID, purchaser))
	}
	sp.Amount = newAmt
	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, purchaser, sdk.NewCoins(sdk.NewCoin(k.sk.BondDenom(ctx), amount))); err != nil {
		return err
	}
	k.SetStakingPurchase(ctx, poolID, purchaser, sp)
	return nil
}

func (k Keeper) FundShieldBlockRewards(ctx sdk.Context, amount sdk.Coins, sender sdk.AccAddress) error {
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, amount); err != nil {
		return err
	}
	blockServiceFee := k.GetBlockServiceFees(ctx)
	blockServiceFee = blockServiceFee.Add(types.NewMixedDecCoins(sdk.NewDecCoinsFromCoins(amount...), sdk.NewDecCoins()))
	k.SetBlockServiceFees(ctx, blockServiceFee)
	return nil
}

func (k Keeper) GetAllStakingPurchases(ctx sdk.Context) (purchases []types.StakingPurchase) {
	k.IterateStakingPurchases(ctx, func(purchase types.StakingPurchase) bool {
		purchases = append(purchases, purchase)
		return false
	})
	return
}

// IterateStakingPurchases iterates through purchase lists in a pool
func (k Keeper) IterateStakingPurchases(ctx sdk.Context, callback func(purchase types.StakingPurchase) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PurchaseListKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var purchase types.StakingPurchase
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &purchase)

		if callback(purchase) {
			break
		}
	}
}
