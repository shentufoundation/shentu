package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/x/shield/types"
)

// TODO: Wrapper for native and foreign service fee

func (k Keeper) SetTotalCollateral(ctx sdk.Context, totalCollateral sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&sdk.IntProto{Int: totalCollateral})
	store.Set(types.GetTotalCollateralKey(), bz)
}

func (k Keeper) GetTotalCollateral(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetTotalCollateralKey())
	if bz == nil {
		panic("total collateral is not found")
	}

	ip := sdk.IntProto{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &ip)
	return ip.Int
}

func (k Keeper) SetTotalWithdrawing(ctx sdk.Context, totalWithdrawing sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&sdk.IntProto{Int: totalWithdrawing})
	store.Set(types.GetTotalWithdrawingKey(), bz)
}

func (k Keeper) GetTotalWithdrawing(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetTotalWithdrawingKey())
	if bz == nil {
		panic("total withdrawing is not found")
	}

	ip := sdk.IntProto{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &ip)
	return ip.Int
}

func (k Keeper) SetTotalShield(ctx sdk.Context, totalShield sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&sdk.IntProto{Int: totalShield})
	store.Set(types.GetTotalShieldKey(), bz)
}

func (k Keeper) GetTotalShield(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetTotalShieldKey())
	if bz == nil {
		panic("total shield is not found")
	}

	ip := sdk.IntProto{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &ip)
	return ip.Int
}

func (k Keeper) SetTotalClaimed(ctx sdk.Context, totalClaimed sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&sdk.IntProto{Int: totalClaimed})
	store.Set(types.GetTotalClaimedKey(), bz)
}

func (k Keeper) GetTotalClaimed(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetTotalClaimedKey())
	if bz == nil {
		panic("total shield is not found")
	}

	ip := sdk.IntProto{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &ip)
	return ip.Int
}

func (k Keeper) SetNativeServiceFee(ctx sdk.Context, nativeServiceFee sdk.DecCoins) {
	store := ctx.KVStore(k.storeKey)
	serviceFee := types.NativeServiceFee{
		NativeServiceFee: nativeServiceFee,
	}
	bz := k.cdc.MustMarshalLengthPrefixed(&serviceFee)
	store.Set(types.GetNativeServiceFeeKey(), bz)
}

func (k Keeper) SetForeignServiceFee(ctx sdk.Context, foreignServiceFee sdk.DecCoins) {
	store := ctx.KVStore(k.storeKey)
	serviceFee := types.ForeignServiceFee{
		ForeignServiceFee: foreignServiceFee,
	}
	bz := k.cdc.MustMarshalLengthPrefixed(&serviceFee)
	store.Set(types.GetForeignServiceFeeKey(), bz)
}

func (k Keeper) GetNativeServiceFee(ctx sdk.Context) sdk.DecCoins {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetNativeServiceFeeKey())
	if bz == nil {
		panic("service fees are not found")
	}
	var nativeServiceFee types.NativeServiceFee
	k.cdc.MustUnmarshalLengthPrefixed(bz, &nativeServiceFee)
	return nativeServiceFee.NativeServiceFee
}

func (k Keeper) GetForeignServiceFee(ctx sdk.Context) sdk.DecCoins {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetForeignServiceFeeKey())
	if bz == nil {
		panic("service fees are not found")
	}
	var foreignServiceFee types.ForeignServiceFee
	k.cdc.MustUnmarshalLengthPrefixed(bz, &foreignServiceFee)
	return foreignServiceFee.ForeignServiceFee
}

func (k Keeper) SetBlockNativeServiceFee(ctx sdk.Context, nativeServiceFee sdk.DecCoins) {
	store := ctx.KVStore(k.storeKey)
	serviceFee := types.NativeServiceFee{
		NativeServiceFee: nativeServiceFee,
	}
	bz := k.cdc.MustMarshalLengthPrefixed(&serviceFee)
	store.Set(types.GetBlockNativeServiceFeeKey(), bz)
}

func (k Keeper) SetBlockForeignServiceFee(ctx sdk.Context, foreignServiceFee sdk.DecCoins) {
	store := ctx.KVStore(k.storeKey)
	serviceFee := types.ForeignServiceFee{
		ForeignServiceFee: foreignServiceFee,
	}
	bz := k.cdc.MustMarshalLengthPrefixed(&serviceFee)
	store.Set(types.GetBlockForeignServiceFeeKey(), bz)
}

func (k Keeper) GetBlockNativeServiceFee(ctx sdk.Context) sdk.DecCoins {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetBlockNativeServiceFeeKey())
	if bz == nil {
		return sdk.DecCoins{}
	}
	var blockNativeServiceFee types.NativeServiceFee
	k.cdc.MustUnmarshalLengthPrefixed(bz, &blockNativeServiceFee)
	return blockNativeServiceFee.NativeServiceFee
}

func (k Keeper) GetBlockForeignServiceFee(ctx sdk.Context) sdk.DecCoins {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetBlockForeignServiceFeeKey())
	if bz == nil {
		return sdk.DecCoins{}
	}
	var blockForeignServiceFee types.ForeignServiceFee
	k.cdc.MustUnmarshalLengthPrefixed(bz, &blockForeignServiceFee)
	return blockForeignServiceFee.ForeignServiceFee
}

func (k Keeper) DeleteBlockNativeServiceFee(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetBlockNativeServiceFeeKey())
}

func (k Keeper) DeleteForeignNativeServiceFee(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetBlockForeignServiceFeeKey())
}

func (k Keeper) SetRemainingNativeServiceFee(ctx sdk.Context, nativeServiceFee sdk.DecCoins) {
	store := ctx.KVStore(k.storeKey)
	serviceFee := types.NativeServiceFee{
		NativeServiceFee: nativeServiceFee,
	}
	bz := k.cdc.MustMarshalLengthPrefixed(&serviceFee)
	store.Set(types.GetRemainingNativeServiceFeeKey(), bz)
}

func (k Keeper) SetRemainingForeignServiceFee(ctx sdk.Context, foreignServiceFee sdk.DecCoins) {
	store := ctx.KVStore(k.storeKey)
	serviceFee := types.ForeignServiceFee{
		ForeignServiceFee: foreignServiceFee,
	}
	bz := k.cdc.MustMarshalLengthPrefixed(&serviceFee)
	store.Set(types.GetRemainingForeignServiceFeeKey(), bz)
}

func (k Keeper) GetRemainingNativeServiceFee(ctx sdk.Context) sdk.DecCoins {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetRemainingNativeServiceFeeKey())
	if bz == nil {
		panic("remaining service fee is not found")
	}
	var remainingNativeServiceFee types.NativeServiceFee
	k.cdc.MustUnmarshalLengthPrefixed(bz, &remainingNativeServiceFee)
	return remainingNativeServiceFee.NativeServiceFee
}

func (k Keeper) GetRemainingForeignServiceFee(ctx sdk.Context) sdk.DecCoins {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetRemainingForeignServiceFeeKey())
	if bz == nil {
		panic("remaining service fee is not found")
	}
	var remainingForeignServiceFee types.ForeignServiceFee
	k.cdc.MustUnmarshalLengthPrefixed(bz, &remainingForeignServiceFee)
	return remainingForeignServiceFee.ForeignServiceFee
}

// SetPool sets data of a pool in kv-store.
func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&pool)
	store.Set(types.GetPoolKey(pool.Id), bz)
}

// GetPool gets data of a pool given pool ID.
func (k Keeper) GetPool(ctx sdk.Context, id uint64) (types.Pool, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPoolKey(id))
	if bz == nil {
		return types.Pool{}, false
	}
	var pool types.Pool
	k.cdc.MustUnmarshalLengthPrefixed(bz, &pool)
	return pool, true
}

// CreatePool creates a pool and sponsor's shield.
func (k Keeper) CreatePool(ctx sdk.Context, creator sdk.AccAddress, shield sdk.Coins, nativeServiceFee sdk.Coins, sponsor string, sponsorAddr sdk.AccAddress, description string, shieldLimit sdk.Int) (uint64, error) {
	admin := k.GetAdmin(ctx)
	if !creator.Equals(admin) {
		return 0, types.ErrNotShieldAdmin
	}
	if _, found := k.GetPoolsBySponsor(ctx, sponsor); found {
		return 0, types.ErrSponsorAlreadyExists
	}

	// Set the new project pool.
	poolID := k.GetNextPoolID(ctx)
	pool := types.NewPool(poolID, description, sponsor, sponsorAddr, shieldLimit, sdk.ZeroInt())
	k.SetPool(ctx, pool)
	k.SetNextPoolID(ctx, poolID+1)

	// Purchase shield for the pool.
	if _, err := k.purchaseShield(ctx, poolID, shield, "shield for sponsor", creator, nativeServiceFee, sdk.NewCoins()); err != nil {
		return poolID, err
	}

	return poolID, nil
}

// UpdatePool updates pool info and shield for B.
func (k Keeper) UpdatePool(ctx sdk.Context, poolID uint64, description string, updater sdk.AccAddress, shield sdk.Coins, nativeServiceFee sdk.Coins, shieldLimit sdk.Int) (types.Pool, error) {
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
		if _, err := k.purchaseShield(ctx, poolID, shield, "shield for sponsor", updater, nativeServiceFee, sdk.NewCoins()); err != nil {
			return pool, err
		}
	} else if !nativeServiceFee.IsZero() {
		// Allow adding service fees without purchasing more shield.
		totalNativeServiceFee := k.GetNativeServiceFee(ctx)
		totalNativeServiceFee = totalNativeServiceFee.Add(totalNativeServiceFee...)
		k.SetNativeServiceFee(ctx, totalNativeServiceFee)
		totalRemainingNativeServiceFee := k.GetRemainingNativeServiceFee(ctx)
		totalRemainingNativeServiceFee = totalRemainingNativeServiceFee.Add(totalNativeServiceFee...)
		k.SetRemainingNativeServiceFee(ctx, totalRemainingNativeServiceFee)
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
	store.Delete(types.GetPoolKey(pool.Id))
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
		k.cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &pool)

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
	pool.SponsorAddr = newSponsorAddr.String()
	k.SetPool(ctx, pool)

	return pool, nil
}
