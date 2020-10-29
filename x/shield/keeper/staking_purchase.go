package keeper

import (
	"github.com/certikfoundation/shentu/x/shield/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetGlobalStakingPurchasePool(ctx sdk.Context) (pool types.GlobalStakingPool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetStakingPurchasePoolKey())
	if bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &pool)
	}
	return
}

func (k Keeper) SetGlobalStakingPurchasePool(ctx sdk.Context, pool types.GlobalStakingPool) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(pool)
	store.Set(types.GetStakingPurchasePoolKey(), bz)
	return
}

func (k Keeper) GetStakingPurchase(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress) (purchase types.StakingPurchase) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetStakingPurchaseKey(poolID, purchaser))
	if bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &purchase)
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

func (k Keeper) AddStaking(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, amount sdk.Int) {
	k.AddGlobalStakingPurchasePool(ctx, amount)
	sp := k.GetStakingPurchase(ctx, poolID, purchaser)
	sp.Locked = sp.Locked.Add(amount)
	sp.Amount = sp.Amount.Add(amount)
	k.SetStakingPurchase(ctx, poolID, purchaser, sp)
}

func (k Keeper) WithdrawStaking(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, amount sdk.Int) error {
	sp := k.GetStakingPurchase(ctx, poolID, purchaser)
	newAmt := sp.Amount.Sub(amount)
	if newAmt.LT(sp.Locked) {
		return types.ErrNotEnoughStaked
	}
	sp.Amount = newAmt
	k.SetStakingPurchase(ctx, poolID, purchaser, sp)
	return nil
}
