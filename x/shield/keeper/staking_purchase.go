package keeper

import (
	"encoding/binary"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

func (k Keeper) GetGlobalStakeForShieldPool(ctx sdk.Context) (pool sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetGlobalStakeForShieldPoolKey())
	if bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &pool)
	}
	return sdk.NewInt(0)
}

func (k Keeper) SetGlobalShieldStakingPool(ctx sdk.Context, pool sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(pool)
	store.Set(types.GetGlobalStakeForShieldPoolKey(), bz)
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

func (k Keeper) GetStakeForShield(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress) (purchase types.ShieldStaking, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetStakeForShieldKey(poolID, purchaser))
	if bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &purchase)
		found = true
	}
	return
}

func (k Keeper) SetStakeForShield(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, purchase types.ShieldStaking) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(purchase)
	store.Set(types.GetStakeForShieldKey(poolID, purchaser), bz)
}

func (k Keeper) AddStaking(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, purchaseID uint64, amount sdk.Int) error {
	pool := k.GetGlobalStakeForShieldPool(ctx)
	pool = pool.Add(amount)
	k.SetGlobalShieldStakingPool(ctx, pool)
	sp, found := k.GetStakeForShield(ctx, poolID, purchaser)
	if !found {
		sp = types.NewShieldStaking(poolID, purchaser, amount)
	}

	sp.Amount = sp.Amount.Add(amount)

	if err := k.supplyKeeper.SendCoinsFromAccountToModule(
		ctx, purchaser, types.ModuleName, sdk.NewCoins(sdk.NewCoin(k.sk.BondDenom(ctx), amount))); err != nil {
		return err
	}
	k.SetStakeForShield(ctx, poolID, purchaser, sp)
	k.SetOriginalStaking(ctx, purchaseID, amount)
	return nil
}

func (k Keeper) UnstakeFromShield(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, purchaseID uint64, amount sdk.Int) error {
	sp, found := k.GetStakeForShield(ctx, poolID, purchaser)
	if !found {
		return types.ErrPurchaseNotFound
	}
	if sp.WithdrawRequested.Add(amount).GT(sp.Amount) {
		return types.ErrNotEnoughStaked
	}
	k.SetStakeForShield(ctx, poolID, purchaser, sp)
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

func (k Keeper) GetAllStakeForShields(ctx sdk.Context) (purchases []types.ShieldStaking) {
	k.IterateStakeForShields(ctx, func(purchase types.ShieldStaking) bool {
		purchases = append(purchases, purchase)
		return false
	})
	return
}

// IterateStakeForShields iterates through purchase lists in a pool
func (k Keeper) IterateStakeForShields(ctx sdk.Context, callback func(purchase types.ShieldStaking) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.StakeForShieldKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var purchase types.ShieldStaking
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &purchase)

		if callback(purchase) {
			break
		}
	}
}

func (k Keeper) GetAllOriginalStakings(ctx sdk.Context) (originalStakings []types.OriginalStaking) {
	k.IterateOriginalStakings(ctx, func(newOS types.OriginalStaking) bool {
		originalStakings = append(originalStakings, newOS)
		return false
	})
	return
}

// IterateStakeForShields iterates through purchase lists in a pool
func (k Keeper) IterateOriginalStakings(ctx sdk.Context, callback func(original types.OriginalStaking) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.OriginalStakingKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var amount sdk.Int
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &amount)
		bz := iterator.Key()[1:]
		id := binary.LittleEndian.Uint64(bz)
		newOS := types.NewOriginalStaking(id, amount)

		if callback(newOS) {
			break
		}
	}
}

func (k Keeper) ProcessStakeForShieldExpiration(ctx sdk.Context, poolID, purchaseID uint64, bondDenom string, purchaser sdk.AccAddress) error {
	stakingPurchase, found := k.GetStakeForShield(ctx, poolID, purchaser)
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
	k.SetStakeForShield(ctx, poolID, purchaser, stakingPurchase)
	if renew.IsZero() {
		return nil
	}

	sPRate := k.GetStakeForShieldRate(ctx)
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
