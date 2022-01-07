package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/x/shield/types"
)

func (k Keeper) GetGlobalStakingPool(ctx sdk.Context) (pool sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetGlobalStakeForShieldPoolKey())
	if bz == nil {
		return sdk.NewInt(0)
	}

	ip := sdk.IntProto{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &ip)
	return ip.Int
}

func (k Keeper) SetGlobalStakingPool(ctx sdk.Context, value sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&sdk.IntProto{Int: value})
	store.Set(types.GetGlobalStakeForShieldPoolKey(), bz)
}

func (k Keeper) GetStakingPurchase(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress) (purchase types.StakingPurchase, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetStakeForShieldKey(poolID, purchaser))
	if bz != nil {
		k.cdc.MustUnmarshalLengthPrefixed(bz, &purchase)
		found = true
	}
	return
}

func (k Keeper) SetStakingPurchase(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, purchase types.StakingPurchase) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&purchase)
	store.Set(types.GetStakeForShieldKey(poolID, purchaser), bz)
}

func (k Keeper) AddStaking(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, purchaseID uint64, stakingAmt sdk.Int) error {
	stakingCoins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), stakingAmt))
	if err := k.bk.SendCoinsFromAccountToModule(ctx, purchaser, types.ModuleName, stakingCoins); err != nil {
		return err
	}

	pool := k.GetGlobalStakingPool(ctx)
	pool = pool.Add(stakingAmt)
	k.SetGlobalStakingPool(ctx, pool)

	sFS, found := k.GetStakingPurchase(ctx, poolID, purchaser)
	if !found {
		sFS = types.NewStakingPurchase(poolID, purchaser, stakingAmt)
	} else {
		sFS.Amount = sFS.Amount.Add(stakingAmt)
	}
	k.SetStakingPurchase(ctx, poolID, purchaser, sFS)
	return nil
}

func (k Keeper) UnstakeFromShield(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, amount sdk.Int) error {
	sp, found := k.GetStakingPurchase(ctx, poolID, purchaser)
	if !found {
		return types.ErrPurchaseNotFound
	}
	k.SetStakingPurchase(ctx, poolID, purchaser, sp)
	return nil
}

func (k Keeper) FundShieldBlockRewards(ctx sdk.Context, amount sdk.Coins, sender sdk.AccAddress) error {
	if err := k.bk.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, amount); err != nil {
		return err
	}
	blockServiceFee := k.GetBlockServiceFees(ctx)
	blockServiceFee = blockServiceFee.Add(types.NewMixedDecCoins(sdk.NewDecCoinsFromCoins(amount...), sdk.NewDecCoins()))
	k.SetBlockServiceFees(ctx, blockServiceFee)
	return nil
}

func (k Keeper) GetAllStakingPurchase(ctx sdk.Context) (purchases []types.StakingPurchase) {
	k.IterateStakingPurchases(ctx, func(purchase types.StakingPurchase) bool {
		purchases = append(purchases, purchase)
		return false
	})
	return
}

// IterateStakingPurchases iterates through purchase lists in a pool
func (k Keeper) IterateStakingPurchases(ctx sdk.Context, callback func(purchase types.StakingPurchase) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.StakeForShieldKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var purchase types.StakingPurchase
		k.cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &purchase)

		if callback(purchase) {
			break
		}
	}
}
