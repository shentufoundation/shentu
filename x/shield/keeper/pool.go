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

func (k Keeper) SetTotalClaimed(ctx sdk.Context, totalClaimed sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(totalClaimed)
	store.Set(types.GetTotalClaimedKey(), bz)
}

func (k Keeper) GetTotalClaimed(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetTotalClaimedKey())
	if bz == nil {
		panic("total shield is not found")
	}
	var totalClaimed sdk.Int
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &totalClaimed)
	return totalClaimed
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
		panic("service fees are not found")
	}
	var serviceFees types.MixedDecCoins
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &serviceFees)
	return serviceFees
}

func (k Keeper) SetBlockServiceFees(ctx sdk.Context, serviceFees types.MixedDecCoins) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(serviceFees)
	store.Set(types.GetBlockServiceFeesKey(), bz)
}

func (k Keeper) GetBlockServiceFees(ctx sdk.Context) types.MixedDecCoins {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetBlockServiceFeesKey())
	if bz == nil {
		return types.InitMixedDecCoins()
	}
	var serviceFees types.MixedDecCoins
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &serviceFees)
	return serviceFees
}

func (k Keeper) DeleteBlockServiceFees(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetBlockServiceFeesKey())
}

func (k Keeper) SetRemainingServiceFees(ctx sdk.Context, serviceFees types.MixedDecCoins) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(serviceFees)
	store.Set(types.GetRemainingServiceFeesKey(), bz)
}

func (k Keeper) GetRemainingServiceFees(ctx sdk.Context) types.MixedDecCoins {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetRemainingServiceFeesKey())
	if bz == nil {
		panic("remaining service fees are not found")
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
func (k Keeper) CreatePool(ctx sdk.Context, creator sdk.AccAddress, shield sdk.Coins, serviceFees types.MixedCoins, sponsor string, sponsorAddr sdk.AccAddress, description string, shieldLimit sdk.Int) (uint64, error) {
	admin := k.GetAdmin(ctx)
	if !creator.Equals(admin) {
		return 0, types.ErrNotShieldAdmin
	}
	if _, found := k.GetPoolBySponsor(ctx, sponsor); found {
		return 0, types.ErrSponsorAlreadyExists
	}

	// Set the new project pool.
	poolID := k.GetNextPoolID(ctx)
	pool := types.NewPool(poolID, description, sponsor, sponsorAddr, shieldLimit, sdk.ZeroInt())
	k.SetPool(ctx, pool)
	k.SetNextPoolID(ctx, poolID+1)

	// Purchase shield for the pool.
	if _, err := k.purchaseShield(ctx, poolID, shield, "shield for sponsor", creator, serviceFees.Native, sdk.NewCoins()); err != nil {
		return poolID, err
	}

	return poolID, nil
}

// UpdatePool updates pool info and shield for B.
func (k Keeper) UpdatePool(ctx sdk.Context, poolID uint64, description string, updater sdk.AccAddress, shield sdk.Coins, serviceFees types.MixedCoins, shieldLimit sdk.Int) (types.Pool, error) {
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
	if !shieldLimit.IsZero() {
		pool.ShieldLimit = shieldLimit
	}
	k.SetPool(ctx, pool)

	// Update purchase and shield.
	if !shield.IsZero() {
		if _, err := k.purchaseShield(ctx, poolID, shield, "shield for sponsor", updater, serviceFees.Native, sdk.NewCoins()); err != nil {
			return pool, err
		}
	} else if !serviceFees.Native.IsZero() {
		// Allow adding service fees without purchasing more shield.
		totalServiceFees := k.GetServiceFees(ctx)
		totalServiceFees = totalServiceFees.Add(types.MixedDecCoins{Native: sdk.NewDecCoinsFromCoins(serviceFees.Native...)})
		k.SetServiceFees(ctx, totalServiceFees)
		totalRemainingServiceFees := k.GetRemainingServiceFees(ctx)
		totalRemainingServiceFees = totalRemainingServiceFees.Add(types.MixedDecCoins{Native: sdk.NewDecCoinsFromCoins(serviceFees.Native...)})
		k.SetRemainingServiceFees(ctx, totalRemainingServiceFees)
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
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetPoolKey(pool.ID))
}

// ClosePools closes pools when both of the pool's shield and shield limit is non-positive.
func (k Keeper) ClosePools(ctx sdk.Context) {
	k.IterateAllPools(ctx, func(pool types.Pool) bool {
		if !pool.Shield.IsPositive() && !pool.ShieldLimit.IsPositive() {
			k.ClosePool(ctx, pool)
		}
		return false
	})
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

// UpdateSponsor updates the sponsor information of a given pool.
func (k Keeper) UpdateSponsor(ctx sdk.Context, poolID uint64, newSponsor string, newSponsorAddr, updater sdk.AccAddress) (types.Pool, error) {
	// Check admin status of the updater.
	admin := k.GetAdmin(ctx)
	if !updater.Equals(admin) {
		return types.Pool{}, types.ErrNotShieldAdmin
	}

	// Retrieve the pool and update its sponsor information.
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return types.Pool{}, types.ErrNoPoolFound
	}
	pool.Sponsor = newSponsor
	pool.SponsorAddress = newSponsorAddr
	k.SetPool(ctx, pool)

	return pool, nil
}
