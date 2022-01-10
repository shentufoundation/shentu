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

func (k Keeper) DeleteStakingPurchase(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetStakingPurchaseKey(poolID, purchaser))
}

func (k Keeper) GetStakingPurchase(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress) (purchase types.StakingPurchase, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetStakingPurchaseKey(poolID, purchaser))
	if bz != nil {
		k.cdc.MustUnmarshalLengthPrefixed(bz, &purchase)
		found = true
	}
	return
}

func (k Keeper) SetStakingPurchase(ctx sdk.Context, purchase types.StakingPurchase) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&purchase)
	purchaser, err := sdk.AccAddressFromBech32(purchase.Purchaser)
	if err != nil {
		panic(err)
	}
	store.Set(types.GetStakingPurchaseKey(purchase.PoolId, purchaser), bz)
}

func (k Keeper) AddStaking(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, amount sdk.Coins) (types.StakingPurchase, error) {

	if err := k.bk.SendCoinsFromAccountToModule(ctx, purchaser, types.ModuleName, amount); err != nil {
		return types.StakingPurchase{}, err
	}

	bondDenomAmt := amount.AmountOf(k.BondDenom(ctx))
	pool := k.GetGlobalStakingPool(ctx)
	pool = pool.Add(bondDenomAmt)
	k.SetGlobalStakingPool(ctx, pool)

	sp, found := k.GetStakingPurchase(ctx, poolID, purchaser)
	if !found {
		sp = types.NewStakingPurchase(poolID, purchaser, bondDenomAmt)
	} else {
		sp.Amount = sp.Amount.Add(bondDenomAmt)
	}
	sp.StartTime = ctx.BlockTime()
	k.SetStakingPurchase(ctx, sp)
	return sp, nil
}

func (k Keeper) Unstake(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, amount sdk.Int) error {
	sp, found := k.GetStakingPurchase(ctx, poolID, purchaser)
	if !found {
		return types.ErrPurchaseNotFound
	}
	if sp.Amount.LT(amount) {
		return types.ErrInsufficientStaking
	}
	poolParams := k.GetPoolParams(ctx)
	cd := poolParams.CooldownPeriod
	if sp.StartTime.Add(cd).After(ctx.BlockTime()) {
		return types.ErrBeforeCooldownEnd
	}
	sp.Amount = sp.Amount.Sub(amount)
	if sp.Amount.Equal(sdk.ZeroInt()) {
		k.DeleteStakingPurchase(ctx, poolID, purchaser)
	} else {
		k.SetStakingPurchase(ctx, sp)
	}

	withdrawCoins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), amount))

	return k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, purchaser, withdrawCoins)
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
	iterator := sdk.KVStorePrefixIterator(store, types.StakingPurchaseKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var purchase types.StakingPurchase
		k.cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &purchase)

		if callback(purchase) {
			break
		}
	}
}
