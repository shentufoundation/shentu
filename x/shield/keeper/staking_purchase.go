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

func (k Keeper) AddStaking(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, amount sdk.Int, endTime time.Time) {
	k.AddGlobalStakingPurchasePool(ctx, amount)
	sp, found := k.GetStakingPurchase(ctx, poolID, purchaser)
	if !found {
		sp = types.NewStakingPurchase(poolID, purchaser, amount)
	}
	sp.Locked = sp.Locked.Add(amount)
	sp.Amount = sp.Amount.Add(amount)
	newExpiration := types.NewStakingExpiration(endTime, amount)
	sp.Expirations = append(sp.Expirations, newExpiration)
	k.SetStakingPurchase(ctx, poolID, purchaser, sp)
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
