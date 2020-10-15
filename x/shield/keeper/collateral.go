package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// GetPoolCollateral retrieves collateral for a pool-provider pair.
func (k Keeper) GetCollateral(ctx sdk.Context, pool types.Pool, addr sdk.AccAddress) (types.Collateral, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetCollateralKey(pool.PoolID, addr))
	if bz == nil {
		return types.Collateral{}, false
	}
	var collateral types.Collateral
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &collateral)
	return collateral, true
}

// GetAllCollaterals gets all collaterals.
func (k Keeper) GetAllCollaterals(ctx sdk.Context) (collaterals []types.Collateral) {
	k.IterateCollaterals(ctx, func(collateral types.Collateral) bool {
		collaterals = append(collaterals, collateral)
		return false
	})
	return collaterals
}

// SetCollateral stores collateral based on pool-provider pair.
func (k Keeper) SetCollateral(ctx sdk.Context, pool types.Pool, addr sdk.AccAddress, collateral types.Collateral) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(collateral)
	store.Set(types.GetCollateralKey(pool.PoolID, addr), bz)
}

// FreeCollateral frees collaterals deposited in a pool.
func (k Keeper) FreeCollaterals(ctx sdk.Context, pool types.Pool) {
	store := ctx.KVStore(k.storeKey)
	k.IteratePoolCollaterals(ctx, pool, func(collateral types.Collateral) bool {
		provider, _ := k.GetProvider(ctx, collateral.Provider)
		provider.Collateral = provider.Collateral.Sub(collateral.Amount)
		k.SetProvider(ctx, collateral.Provider, provider)
		store.Delete(types.GetCollateralKey(pool.PoolID, collateral.Provider))
		return false
	})
}

// IterateCollaterals iterates through all collaterals.
func (k Keeper) IterateCollaterals(ctx sdk.Context, callback func(collateral types.Collateral) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.CollateralKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var collateral types.Collateral
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &collateral)

		if callback(collateral) {
			break
		}
	}
}

// IteratePoolCollaterals iterates through collaterals in a pool
func (k Keeper) IteratePoolCollaterals(ctx sdk.Context, pool types.Pool, callback func(collateral types.Collateral) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetPoolCollateralsKey(pool.PoolID))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var collateral types.Collateral
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &collateral)

		if callback(collateral) {
			break
		}
	}
}

// GetProviderCollaterals returns a community member's all collaterals.
func (k Keeper) GetProviderCollaterals(ctx sdk.Context, address sdk.AccAddress) (collaterals []types.Collateral) {
	k.IterateAllPools(ctx, func(pool types.Pool) bool {
		collateral, found := k.GetCollateral(ctx, pool, address)
		if found {
			collaterals = append(collaterals, collateral)
		}
		return false
	})
	return collaterals
}

// GetPoolCertiKCollateral retrieves CertiK's provided collateral from a pool.
func (k Keeper) GetPoolCertiKCollateral(ctx sdk.Context, pool types.Pool) (collateral types.Collateral, found bool) {
	admin := k.GetAdmin(ctx)
	collateral, found = k.GetCollateral(ctx, pool, admin)
	return
}

// GetAllPoolCollaterals retrieves all collaterals in a pool.
func (k Keeper) GetAllPoolCollaterals(ctx sdk.Context, pool types.Pool) (collaterals []types.Collateral) {
	k.IteratePoolCollaterals(ctx, pool, func(collateral types.Collateral) bool {
		collaterals = append(collaterals, collateral)
		return false
	})
	return collaterals
}

// DepositCollateral deposits a community member's collateral for a pool.
func (k Keeper) DepositCollateral(ctx sdk.Context, from sdk.AccAddress, id uint64, amount sdk.Int) error {
	pool, err := k.GetPool(ctx, id)
	if err != nil {
		return err
	}

	// check eligibility
	provider, found := k.GetProvider(ctx, from)
	if !found {
		provider = k.addProvider(ctx, from)
	}
	provider.Collateral = provider.Collateral.Add(amount)
	if amount.GT(provider.Available) {
		return types.ErrInsufficientStaking
	}
	provider.Available = provider.Available.Sub(amount)

	// update the pool, collateral and provider
	collateral, found := k.GetCollateral(ctx, pool, from)
	if !found {
		collateral = types.NewCollateral(pool, from, amount)
	} else {
		collateral.Amount = collateral.Amount.Add(amount)
	}
	pool.TotalCollateral = pool.TotalCollateral.Add(amount)
	pool.Available = pool.Available.Add(amount)
	k.SetPool(ctx, pool)
	k.SetCollateral(ctx, pool, from, collateral)
	k.SetProvider(ctx, from, provider)

	return nil
}

// WithdrawCollateral withdraws a community member's collateral for a pool.
func (k Keeper) WithdrawCollateral(ctx sdk.Context, from sdk.AccAddress, id uint64, amount sdk.Int) error {
	if amount.IsZero() {
		return nil
	}
	pool, err := k.GetPool(ctx, id)
	if err != nil {
		return err
	}

	// retrieve the particular collateral to ensure that
	// amount is less than collateral minus collateral withdraw
	collateral, found := k.GetCollateral(ctx, pool, from)
	if !found {
		return types.ErrNoCollateralFound
	}
	withdrawable := collateral.Amount.Sub(collateral.Withdrawing)
	if amount.GT(withdrawable) {
		return types.ErrOverWithdraw
	}

	// update the pool available coins, but not pool total collateral or community which should be updated 21 days later
	pool.Available = pool.Available.Sub(amount)
	k.SetPool(ctx, pool)

	// insert into withdraw queue
	poolParams := k.GetPoolParams(ctx)
	completionTime := ctx.BlockHeader().Time.Add(poolParams.WithdrawPeriod)
	withdraw := types.NewWithdraw(id, from, amount, completionTime)
	k.InsertWithdrawQueue(ctx, withdraw)

	collateral.Withdrawing = collateral.Withdrawing.Add(amount)
	k.SetCollateral(ctx, pool, collateral.Provider, collateral)

	provider, found := k.GetProvider(ctx, from)
	if !found {
		return types.ErrProviderNotFound
	}
	provider.Withdrawing = provider.Withdrawing.Add(amount)
	k.SetProvider(ctx, provider.Address, provider)

	return nil
}
