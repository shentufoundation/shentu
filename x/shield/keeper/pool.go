package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

func (k Keeper) SetTotalCollateral(ctx sdk.Context, totalCollateral sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(totalCollateral)
	store.Set(types.GetTotalCollateralKey(), bz)
}

func (k Keeper) GetTotalCollateral(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetTotalCollateralKey())
	if bz == nil {
		panic("total collateral is not found")
	}
	var totalCollateral sdk.Int
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &totalCollateral)
	return totalCollateral
}

func (k Keeper) SetTotalWithdrawing(ctx sdk.Context, totalWithdrawing sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(totalWithdrawing)
	store.Set(types.GetTotalWithdrawingKey(), bz)
}

func (k Keeper) GetTotalWithdrawing(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetTotalWithdrawingKey())
	if bz == nil {
		panic("total withdrawing is not found")
	}
	var totalWithdrawing sdk.Int
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &totalWithdrawing)
	return totalWithdrawing
}

func (k Keeper) SetTotalShield(ctx sdk.Context, totalShield sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(totalShield)
	store.Set(types.GetTotalShieldKey(), bz)
}

func (k Keeper) GetTotalShield(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetTotalShieldKey())
	if bz == nil {
		panic("total shield is not found")
	}
	var totalShield sdk.Int
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &totalShield)
	return totalShield
}

func (k Keeper) SetTotalLocked(ctx sdk.Context, totalLocked sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(totalLocked)
	store.Set(types.GetTotalLockedKey(), bz)
}

func (k Keeper) GetTotalLocked(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetTotalLockedKey())
	if bz == nil {
		panic("total shield is not found")
	}
	var totalLocked sdk.Int
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &totalLocked)
	return totalLocked
}

func (k Keeper) SetServiceFees(ctx sdk.Context, serviceFees types.MixedDecCoins) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(serviceFees)
	store.Set(types.GetServiceFeesKey(), bz)
}

func (k Keeper) GetServiceFees(ctx sdk.Context) types.MixedDecCoins {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetServiceFeesKey())
	if bz == nil {
		panic("service fees is not found")
	}
	var serviceFees types.MixedDecCoins
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &serviceFees)
	return serviceFees
}

func (k Keeper) SetServiceFeesPerSecond(ctx sdk.Context, serviceFees types.MixedDecCoins) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(serviceFees)
	store.Set(types.GetServiceFeesPerSecondKey(), bz)
}

func (k Keeper) GetServiceFeesPerSecond(ctx sdk.Context) types.MixedDecCoins {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetServiceFeesPerSecondKey())
	if bz == nil {
		panic("service fees per second is not found")
	}
	var serviceFees types.MixedDecCoins
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &serviceFees)
	return serviceFees
}

// SetPool sets data of a pool in kv-store.
func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(pool)
	store.Set(types.GetPoolKey(pool.ID), bz)
}

// GetPool gets data of a pool given pool ID.
func (k Keeper) GetPool(ctx sdk.Context, id uint64) (types.Pool, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPoolKey(id))
	if bz == nil {
		return types.Pool{}, false
	}
	var pool types.Pool
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &pool)
	return pool, true
}

// CreatePool creates a pool and sponsor's shield.
func (k Keeper) CreatePool(ctx sdk.Context, creator sdk.AccAddress, shield sdk.Coins, serviceFees types.MixedCoins, sponsor string, sponsorAddr sdk.AccAddress, description string, shieldLimit sdk.Coins) (uint64, error) {
	admin := k.GetAdmin(ctx)
	if !creator.Equals(admin) {
		return 0, types.ErrNotShieldAdmin
	}
	if _, found := k.GetPoolBySponsor(ctx, sponsor); found {
		return 0, types.ErrSponsorAlreadyExists
	}

	// Set the new project pool.
	poolID := k.GetNextPoolID(ctx)
	pool := types.NewPool(poolID, description, sponsor, sponsorAddr, shieldLimit.AmountOf(k.BondDenom(ctx)), sdk.ZeroInt())
	k.SetPool(ctx, pool)
	k.SetNextPoolID(ctx, poolID+1)

	// Pool is created before the purchase is made.
	// The pool will still exist even if the purchase fails.
	if _, err := k.purchaseShield(ctx, poolID, shield, "shield for sponsor", creator, serviceFees.Native); err != nil {
		return poolID, err
	}

	return poolID, nil
}

// UpdatePool updates pool info and shield for B.
func (k Keeper) UpdatePool(ctx sdk.Context, poolID uint64, description string, updater sdk.AccAddress, shield sdk.Coins, serviceFees types.MixedCoins, shieldLimit sdk.Coins) (types.Pool, error) {
	admin := k.GetAdmin(ctx)
	if !updater.Equals(admin) {
		return types.Pool{}, types.ErrNotShieldAdmin
	}

	// Update pool info.
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return types.Pool{}, types.ErrNoPoolFound
	}
	if description != "" {
		pool.Description = description
	}
	shieldLimitAmt := shieldLimit.AmountOf(k.BondDenom(ctx))
	if shieldLimitAmt.IsPositive() {
		pool.ShieldLimit = shieldLimitAmt
	}
	k.SetPool(ctx, pool)

	// Update purchase and shield.
	if !shield.IsZero() {
		if _, err := k.purchaseShield(ctx, poolID, shield, "shield for sponsor", updater, serviceFees.Native); err != nil {
			return pool, err
		}
	} else if !serviceFees.Native.IsZero() {
		// Allow adding service fees without purchasing more shield.
		totalServiceFees := k.GetServiceFees(ctx)
		totalServiceFees = totalServiceFees.Add(types.MixedDecCoins{Native: sdk.NewDecCoinsFromCoins(serviceFees.Native...)})
		k.SetServiceFees(ctx, totalServiceFees)
	}

	return pool, nil
}

// PausePool sets an active pool to be inactive.
func (k Keeper) PausePool(ctx sdk.Context, updater sdk.AccAddress, id uint64) (types.Pool, error) {
	admin := k.GetAdmin(ctx)
	if !updater.Equals(admin) {
		return types.Pool{}, types.ErrNotShieldAdmin
	}
	pool, found := k.GetPool(ctx, id)
	if !found {
		return types.Pool{}, types.ErrNoPoolFound
	}
	if !pool.Active {
		return types.Pool{}, types.ErrPoolAlreadyPaused
	}
	pool.Active = false
	k.SetPool(ctx, pool)
	return pool, nil
}

// ResumePool sets an inactive pool to be active.
func (k Keeper) ResumePool(ctx sdk.Context, updater sdk.AccAddress, id uint64) (types.Pool, error) {
	admin := k.GetAdmin(ctx)
	if !updater.Equals(admin) {
		return types.Pool{}, types.ErrNotShieldAdmin
	}
	pool, found := k.GetPool(ctx, id)
	if !found {
		return types.Pool{}, types.ErrNoPoolFound
	}
	if pool.Active {
		return types.Pool{}, types.ErrPoolAlreadyActive
	}
	pool.Active = true
	k.SetPool(ctx, pool)
	return pool, nil
}

// GetAllPools retrieves all pools in the store.
func (k Keeper) GetAllPools(ctx sdk.Context) (pools []types.Pool) {
	k.IterateAllPools(ctx, func(pool types.Pool) bool {
		pools = append(pools, pool)
		return false
	})
	return pools
}

// ClosePool closes the pool.
func (k Keeper) ClosePool(ctx sdk.Context, pool types.Pool) {
	// TODO: make sure nothing else needs to be done
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetPoolKey(pool.ID))
}

// IterateAllPools iterates over the all the stored pools and performs a callback function.
func (k Keeper) IterateAllPools(ctx sdk.Context, callback func(pool types.Pool) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PoolKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var pool types.Pool
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &pool)

		if callback(pool) {
			break
		}
	}
}
