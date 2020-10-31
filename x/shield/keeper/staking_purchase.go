package keeper

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

func (k Keeper) GetGlobalStakingPurchasePool(ctx sdk.Context) (pool sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetGlobalStakingPurchasePoolKey())
	if bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &pool)
	}
	return sdk.NewInt(0)
}

func (k Keeper) SetGlobalStakingPurchasePool(ctx sdk.Context, pool sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(pool)
	store.Set(types.GetGlobalStakingPurchasePoolKey(), bz)
}

func (k Keeper) GetOriginalStaking(ctx sdk.Context, purchaseID uint64) (amount sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetOriginalStakingKey(purchaseID))
	if bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &amount)
	}
	return sdk.NewInt(0)
}

func (k Keeper) SetOriginalStaking(ctx sdk.Context, purchaseID uint64, amount sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(amount)
	store.Set(types.GetOriginalStakingKey(purchaseID), bz)
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

func (k Keeper) AddStaking(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, purchaseID uint64, amount sdk.Int) error {
	pool := k.GetGlobalStakingPurchasePool(ctx)
	pool = pool.Add(amount)
	k.SetGlobalStakingPurchasePool(ctx, pool)
	sp, found := k.GetStakingPurchase(ctx, poolID, purchaser)
	if !found {
		sp = types.NewStakingPurchase(poolID, purchaser, amount)
	}

	sp.Amount = sp.Amount.Add(amount)

	err := k.supplyKeeper.SendCoinsFromAccountToModule(
		ctx, purchaser, types.ModuleName, sdk.NewCoins(sdk.NewCoin(k.sk.BondDenom(ctx), amount)))
	if err != nil {
		return err
	}
	k.SetStakingPurchase(ctx, poolID, purchaser, sp)
	k.SetOriginalStaking(ctx, purchaseID, amount)
	return nil
}

func (k Keeper) WithdrawStaking(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, purchaseID uint64, amount sdk.Int) error {
	sp, found := k.GetStakingPurchase(ctx, poolID, purchaser)
	if !found {
		return types.ErrPurchaseNotFound
	}
	if sp.WithdrawRequested.Add(amount).GT(sp.Amount) {
		return types.ErrNotEnoughStaked
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
	iterator := sdk.KVStorePrefixIterator(store, types.StakingPurchaseKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var purchase types.StakingPurchase
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &purchase)

		if callback(purchase) {
			break
		}
	}
}

func (k Keeper) ProcessStakingPurchaseExpiration(ctx sdk.Context, poolID, purchaseID uint64, bondDenom string, purchaser sdk.AccAddress) error {
	stakingPurchase, found := k.GetStakingPurchase(ctx, poolID, purchaser)
	if !found {
		return nil
	}
	amount := k.GetOriginalStaking(ctx, purchaseID)
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetOriginalStakingKey(purchaseID))

	renew := amount.Sub(stakingPurchase.WithdrawRequested)
	if renew.IsNegative() {
		renew = sdk.NewInt(0)
	}
	stakingPurchase.WithdrawRequested = sdk.MaxInt(sdk.ZeroInt(), stakingPurchase.WithdrawRequested.Sub(amount))
	stakingPurchase.Amount = stakingPurchase.Amount.Sub(amount)
	k.SetStakingPurchase(ctx, poolID, purchaser, stakingPurchase)
	if renew.IsZero() {
		return nil
	}

	sPRate := k.GetStakingPurchaseRate(ctx)
	renewShieldInt := sPRate.QuoInt(amount).TruncateInt()
	renewShield := sdk.NewCoins(sdk.NewCoin(bondDenom, renewShieldInt))
	desc := fmt.Sprintf(`renewed from PurchaseID %s`, strconv.FormatUint(purchaseID, 10))
	if _, err := k.purchaseShield(ctx, poolID, renewShield, desc, purchaser,
		sdk.NewCoins(), sdk.NewCoins(sdk.NewCoin(bondDenom, renew))); err != nil {
		panic(err)
	}
	withdrawAmt := amount.Sub(renew)
	withdrawCoins := sdk.NewCoins(sdk.NewCoin(bondDenom, withdrawAmt))
	k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, purchaser, withdrawCoins)
	return nil
}
