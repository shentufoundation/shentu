package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/common"
	"github.com/certikfoundation/shentu/v2/x/shield/types"
	"github.com/certikfoundation/shentu/v2/x/shield/types/v1beta1"
)

// PurchaseShield purchases shield of a pool with standard fee rate.
func (k Keeper) PurchaseShield(ctx sdk.Context, poolID uint64, amount sdk.Coins, description string, purchaser sdk.AccAddress) (v1beta1.Purchase, error) {
	poolParams := k.GetPoolParams(ctx)
	if poolParams.MinShieldPurchase.IsAnyGT(amount) {
		return v1beta1.Purchase{}, types.ErrPurchaseTooSmall
	}

	bondDenom := k.BondDenom(ctx)
	if amount.AmountOf(bondDenom).Equal(sdk.ZeroInt()) {
		return v1beta1.Purchase{}, types.ErrInsufficientStaking
	}
	pool, found := k.GetPool(ctx, poolID)

	if !found {
		return v1beta1.Purchase{}, types.ErrNoPoolFound
	}
	if !pool.Active {
		return v1beta1.Purchase{}, types.ErrPoolInactive
	}
	if amount.Empty() {
		return v1beta1.Purchase{}, types.ErrNoShield
	}

	sp, err := k.AddStaking(ctx, poolID, purchaser, description, amount)

	return sp, err
}

func (k Keeper) GetPurchaserPurchases(ctx sdk.Context, address sdk.AccAddress) (res []v1beta1.Purchase) {
	purchases := k.GetAllPurchase(ctx)

	for _, p := range purchases {
		if p.Purchaser == address.String() {
			res = append(res, p)
		}
	}
	return
}

func (k Keeper) GetPoolPurchases(ctx sdk.Context, poolID uint64) (res []v1beta1.Purchase) {
	purchases := k.GetAllPurchase(ctx)

	for _, p := range purchases {
		if p.PoolId == poolID {
			res = append(res, p)
		}
	}
	return
}

// DistributeFees distributes rewards for current block plus leftover rewards for last block.
func (k Keeper) DistributeFees(ctx sdk.Context) {
	// Add leftover block service fees from last block
	serviceFees := k.GetServiceFees(ctx)

	// Distribute service fees.
	totalCollateral := k.GetTotalCollateral(ctx)
	if totalCollateral.IsZero() {
		return
	}
	providers := k.GetAllProviders(ctx)
	for _, provider := range providers {
		providerAddr, err := sdk.AccAddressFromBech32(provider.Address)
		if err != nil {
			panic(err)
		}

		// fees * providerCollateral / totalCollateral
		nativeFees := serviceFees.MulDec(sdk.NewDecFromInt(provider.Collateral).QuoInt(totalCollateral))
		provider.Rewards = provider.Rewards.Add(nativeFees...)
		k.SetProvider(ctx, providerAddr, provider)
		serviceFees = serviceFees.Sub(nativeFees)
	}

	// Store remaining block reward as new leftover
	k.SetServiceFees(ctx, serviceFees)
}

func (k Keeper) RecoverPurchases(ctx sdk.Context) {
	bondDenom := k.BondDenom(ctx)
	k.IteratePurchases(ctx, func(purchase v1beta1.Purchase) bool {
		var updated []v1beta1.RecoveringEntry
		pool, found := k.GetPool(ctx, purchase.PoolId)
		total := k.GetTotalShield(ctx)
		if !found {
			panic("pool not found for an existing purchase")
		}
		for _, e := range purchase.RecoveringEntries {
			if e.RecoverTime.Before(ctx.BlockTime()) {
				purchase.Shield = purchase.Shield.Add(e.Amount.AmountOf(bondDenom))
				pool.Shield = pool.Shield.Add(e.Amount.AmountOf(bondDenom))
				total = total.Add(e.Amount.AmountOf(bondDenom))
			} else {
				updated = append(updated, e)
			}
		}
		purchase.RecoveringEntries = updated
		k.SetPurchase(ctx, purchase)
		k.SetPool(ctx, pool)
		k.SetTotalShield(ctx, total)
		return false
	})
}

func (k Keeper) GetGlobalStakingPool(ctx sdk.Context) (pool sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetGlobalStakingPoolKey())
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
	store.Set(types.GetGlobalStakingPoolKey(), bz)
}

func (k Keeper) DeletePurchase(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetPurchaseKey(poolID, purchaser))
}

func (k Keeper) GetPurchase(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress) (purchase v1beta1.Purchase, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPurchaseKey(poolID, purchaser))
	if bz != nil {
		k.cdc.MustUnmarshalLengthPrefixed(bz, &purchase)
		found = true
	}
	return
}

func (k Keeper) SetPurchase(ctx sdk.Context, purchase v1beta1.Purchase) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&purchase)
	purchaser, err := sdk.AccAddressFromBech32(purchase.Purchaser)
	if err != nil {
		panic(err)
	}
	store.Set(types.GetPurchaseKey(purchase.PoolId, purchaser), bz)
}

func (k Keeper) AddStaking(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, description string, amount sdk.Coins) (v1beta1.Purchase, error) {
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return v1beta1.Purchase{}, types.ErrNoPoolFound
	}

	bondDenomAmt := amount.AmountOf(k.BondDenom(ctx))

	shieldAmt := bondDenomAmt.ToDec().Mul(pool.ShieldRate).TruncateInt()

	totalPurchaseAmt := pool.Shield.Add(shieldAmt)
	if pool.ShieldLimit.LT(totalPurchaseAmt) {
		return v1beta1.Purchase{}, types.ErrPurchaseExceededLimit
	}

	pool.Shield = totalPurchaseAmt
	k.SetPool(ctx, pool)

	gSPool := k.GetGlobalStakingPool(ctx)
	gSPool = gSPool.Add(bondDenomAmt)
	k.SetGlobalStakingPool(ctx, gSPool)

	sp, found := k.GetPurchase(ctx, poolID, purchaser)
	if !found {
		sp = v1beta1.NewPurchase(poolID, purchaser, description, bondDenomAmt, shieldAmt)
	} else {
		if sp.Locked {
			return v1beta1.Purchase{}, types.ErrPurchaseLocked
		}
		sp.Amount = sp.Amount.Add(bondDenomAmt)
		sp.Shield = sp.Shield.Add(shieldAmt)
	}
	sp.StartTime = ctx.BlockTime()
	k.SetPurchase(ctx, sp)

	totalShield := k.GetTotalShield(ctx)
	totalShield = totalShield.Add(shieldAmt)
	k.SetTotalShield(ctx, totalShield)

	if err := k.bk.SendCoinsFromAccountToModule(ctx, purchaser, types.ModuleName, amount); err != nil {
		return v1beta1.Purchase{}, err
	}
	return sp, nil
}

func (k Keeper) Unstake(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, amount sdk.Coins) error {
	bdAmount := amount.AmountOf(k.BondDenom(ctx))

	sp, found := k.GetPurchase(ctx, poolID, purchaser)
	if !found {
		return types.ErrPurchaseNotFound
	}
	if sp.Locked {
		return types.ErrPurchaseLocked
	}
	if sp.Amount.LT(bdAmount) {
		return types.ErrInsufficientStaking
	}
	poolParams := k.GetPoolParams(ctx)
	cd := poolParams.CooldownPeriod
	fees := sdk.ZeroInt()
	if sp.StartTime.Add(cd).After(ctx.BlockTime()) {
		fees = bdAmount.ToDec().Mul(poolParams.WithdrawFeesRate).QuoInt(sdk.NewInt(100)).TruncateInt()
		reserve := k.GetReserve(ctx)
		reserve.Amount = reserve.Amount.Add(fees)
		k.SetReserve(ctx, reserve)
	}

	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return types.ErrNoPoolFound
	}

	// update shield amount

	// calculate the shield amount to be deducted
	var recShield sdk.Coins
	for _, rc := range sp.RecoveringEntries {
		recShield.Add(rc.Amount...)
	}
	// TODO: support multiple coins (?)
	purchasedShield := recShield.AmountOf(k.BondDenom(ctx)).Add(sp.Shield)
	withdrawRatio := purchasedShield.ToDec().Quo(sp.Amount.ToDec())
	shieldReducAmt := common.MulCoins(amount, withdrawRatio)

	var updatedRE []v1beta1.RecoveringEntry
	for _, e := range sp.RecoveringEntries {
		if e.Amount.IsAllLTE(shieldReducAmt) {
			shieldReducAmt = shieldReducAmt.Sub(e.Amount)
			continue
		} else if shieldReducAmt.IsAllLTE(e.Amount) {
			e.Amount = e.Amount.Sub(shieldReducAmt)
			updatedRE = append(updatedRE, e)
			shieldReducAmt = sdk.NewCoins()
		} else if shieldReducAmt.Empty() {
			updatedRE = append(updatedRE, e)
		}
	}
	sp.RecoveringEntries = updatedRE

	totalShield := k.GetTotalShield(ctx)
	if !shieldReducAmt.IsZero() {
		sp.Shield = sp.Shield.Sub(shieldReducAmt.AmountOf(k.BondDenom(ctx)))

		// pool shield is already decreased for the loss amount when the claim proposal is submitted
		pool.Shield = pool.Shield.Sub(shieldReducAmt.AmountOf(k.BondDenom(ctx)))
		k.SetPool(ctx, pool)

		// update total shield
		newTotalShield := totalShield.Sub(shieldReducAmt.AmountOf(k.BondDenom(ctx)))
		k.SetTotalShield(ctx, newTotalShield)
	}

	sp.Amount = sp.Amount.Sub(bdAmount)
	if sp.Amount.Equal(sdk.ZeroInt()) {
		k.DeletePurchase(ctx, poolID, purchaser)
	} else {
		sp.StartTime = ctx.BlockTime()
		k.SetPurchase(ctx, sp)
	}

	// update global pool
	bondDenomAmt := bdAmount
	gSPool := k.GetGlobalStakingPool(ctx)
	gSPool = gSPool.Sub(bondDenomAmt)
	k.SetGlobalStakingPool(ctx, gSPool)

	withdraw := bdAmount.Sub(fees)
	withdrawCoins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), withdraw))
	return k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, purchaser, withdrawCoins)
}

func (k Keeper) FundShieldFees(ctx sdk.Context, amount sdk.Coins, sender sdk.AccAddress) error {
	if err := k.bk.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, amount); err != nil {
		return err
	}
	blockServiceFee := k.GetServiceFees(ctx)
	blockServiceFee = blockServiceFee.Add(sdk.NewDecCoinsFromCoins(amount...)...)
	k.SetServiceFees(ctx, blockServiceFee)
	return nil
}

func (k Keeper) GetAllPurchase(ctx sdk.Context) (purchases []v1beta1.Purchase) {
	k.IteratePurchases(ctx, func(purchase v1beta1.Purchase) bool {
		purchases = append(purchases, purchase)
		return false
	})
	return
}

// IteratePurchases iterates through purchase lists in a pool
func (k Keeper) IteratePurchases(ctx sdk.Context, callback func(purchase v1beta1.Purchase) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PurchaseKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var purchase v1beta1.Purchase
		k.cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &purchase)

		if callback(purchase) {
			break
		}
	}
}
