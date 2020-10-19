package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// GetPoolCollateral retrieves collateral for a pool-provider pair.
func (k Keeper) GetCollateral(ctx sdk.Context, pool types.Pool, addr sdk.AccAddress) (types.Collateral, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetCollateralKey(pool.ID, addr))
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
	store.Set(types.GetCollateralKey(pool.ID, addr), bz)
}

// FreeCollateral frees collaterals deposited in a pool.
func (k Keeper) FreeCollaterals(ctx sdk.Context, pool types.Pool) {
	store := ctx.KVStore(k.storeKey)
	k.IteratePoolCollaterals(ctx, pool, func(collateral types.Collateral) bool {
		provider, _ := k.GetProvider(ctx, collateral.Provider)
		provider.Collateral = provider.Collateral.Sub(collateral.Amount)
		provider.Available = provider.Available.Add(collateral.Amount)
		provider.Withdrawing = provider.Withdrawing.Sub(collateral.Withdrawing)
		k.SetProvider(ctx, collateral.Provider, provider)
		store.Delete(types.GetCollateralKey(pool.ID, collateral.Provider))
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

// IteratePoolCollaterals iterates through collaterals in a pool.
func (k Keeper) IteratePoolCollaterals(ctx sdk.Context, pool types.Pool, callback func(collateral types.Collateral) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetPoolCollateralsKey(pool.ID))

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
func (k Keeper) DepositCollateral(ctx sdk.Context, from sdk.AccAddress, amount sdk.Int) error {
	totalCollateral := k.GetTotalCollateral(ctx)

	// Check eligibility.
	provider, found := k.GetProvider(ctx, from)
	if !found {
		provider = k.addProvider(ctx, from)
	}
	provider.Collateral = provider.Collateral.Add(amount)
	if amount.GT(provider.Available) {
		return types.ErrInsufficientStaking
	}
	provider.Available = provider.Available.Sub(amount)

	totalCollateral = totalCollateral.Add(amount)
	k.SetTotalCollateral(ctx, totalCollateral)
	k.SetProvider(ctx, from, provider)

	return nil
}

// WithdrawCollateral withdraws a community member's collateral for a pool.
func (k Keeper) WithdrawCollateral(ctx sdk.Context, from sdk.AccAddress, amount sdk.Int) error {
	if amount.IsZero() {
		return nil
	}

	provider, found := k.GetProvider(ctx, from)
	if !found {
		return types.ErrProviderNotFound
	}
	withdrawable := provider.Collateral.Sub(provider.Withdrawing)
	if amount.GT(withdrawable) {
		return types.ErrOverWithdraw
	}

	// Insert into withdraw queue.
	poolParams := k.GetPoolParams(ctx)
	completionTime := ctx.BlockHeader().Time.Add(poolParams.WithdrawPeriod)
	withdraw := types.NewWithdraw(from, amount, completionTime)
	k.InsertWithdrawQueue(ctx, withdraw)

	provider.Withdrawing = provider.Withdrawing.Add(amount)
	k.SetProvider(ctx, provider.Address, provider)

	return nil
}
