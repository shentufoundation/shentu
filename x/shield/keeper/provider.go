package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/x/shield/types"
	"github.com/certikfoundation/shentu/v2/x/shield/types/v1beta1"
)

// SetProvider sets data of a provider in the kv-store.
func (k Keeper) SetProvider(ctx sdk.Context, delAddr sdk.AccAddress, provider v1beta1.Provider) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&provider)
	store.Set(types.GetProviderKey(delAddr), bz)
}

// GetProvider returns data of a provider given its address.
func (k Keeper) GetProvider(ctx sdk.Context, delegator sdk.AccAddress) (dt v1beta1.Provider, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetProviderKey(delegator))
	if bz == nil {
		return v1beta1.Provider{}, false
	}
	k.cdc.MustUnmarshal(bz, &dt)
	return dt, true
}

// addProvider adds a new provider into shield module.
// Should only be called from CreatePool or DepositCollateral.
func (k Keeper) addProvider(ctx sdk.Context, addr sdk.AccAddress) v1beta1.Provider {
	delegations := k.sk.GetAllDelegatorDelegations(ctx, addr)

	// Track provider's total stakings.
	totalStaked := sdk.ZeroInt()
	for _, del := range delegations {
		val, found := k.sk.GetValidator(ctx, del.GetValidatorAddr())
		if !found {
			panic("expected validator, not found")
		}
		totalStaked = totalStaked.Add(val.TokensFromShares(del.GetShares()).TruncateInt())
	}

	provider := v1beta1.NewProvider(addr)
	provider.DelegationBonded = totalStaked
	k.SetProvider(ctx, addr, provider)
	return provider
}

// UpdateDelegationAmount updates the provider based on tha changes of its delegations.
func (k Keeper) UpdateDelegationAmount(ctx sdk.Context, delAddr sdk.AccAddress) {
	// Go through delAddr's delegations to recompute total amount of bonded delegation
	// update or create a new entry.
	if _, found := k.GetProvider(ctx, delAddr); !found {
		return // ignore non-participating addr
	}

	// Calculate the amount of its total delegations.
	totalStakedAmount := sdk.ZeroInt()
	delegations := k.sk.GetAllDelegatorDelegations(ctx, delAddr)
	for _, del := range delegations {
		val, found := k.sk.GetValidator(ctx, del.GetValidatorAddr())
		if !found {
			panic("expected validator, not found")
		}
		totalStakedAmount = totalStakedAmount.Add(val.TokensFromShares(del.GetShares()).TruncateInt())
	}

	k.updateProviderForDelegationChanges(ctx, delAddr, totalStakedAmount)
}

// RemoveDelegation updates the provider when its delegation is removed.
func (k Keeper) RemoveDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	provider, found := k.GetProvider(ctx, delAddr)
	if !found {
		return
	}

	delegation, found := k.sk.GetDelegation(ctx, delAddr, valAddr)
	if !found {
		panic("delegation is not found")
	}
	validator, found := k.sk.GetValidator(ctx, valAddr)
	if !found {
		panic("validator is not found")
	}
	deltaAmount := validator.TokensFromShares(delegation.Shares).TruncateInt()

	k.updateProviderForDelegationChanges(ctx, delAddr, provider.DelegationBonded.Sub(deltaAmount))
}

// updateProviderForDelegationChanges updates provider based on delegation changes.
func (k Keeper) updateProviderForDelegationChanges(ctx sdk.Context, delAddr sdk.AccAddress, stakedAmt sdk.Int) {
	provider, found := k.GetProvider(ctx, delAddr)
	if !found {
		return
	}

	// Update the provider.
	provider.DelegationBonded = stakedAmt
	k.SetProvider(ctx, delAddr, provider)

	// Withdraw collaterals when the delegations are not enough to back collaterals.
	withdrawAmount := provider.Collateral.Sub(provider.Withdrawing).Sub(stakedAmt)
	if withdrawAmount.IsPositive() {
		if err := k.WithdrawCollateral(ctx, delAddr, withdrawAmount); err != nil {
			panic("failed to withdraw collateral from the shield global pool")
		}
	}
}

// IterateProviders iterates through all providers.
func (k Keeper) IterateProviders(ctx sdk.Context, callback func(provider v1beta1.Provider) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ProviderKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var provider v1beta1.Provider
		k.cdc.MustUnmarshal(iterator.Value(), &provider)

		if callback(provider) {
			break
		}
	}
}

// GetAllProviders retrieves all providers.
func (k Keeper) GetAllProviders(ctx sdk.Context) (providers []v1beta1.Provider) {
	k.IterateProviders(ctx, func(provider v1beta1.Provider) bool {
		providers = append(providers, provider)
		return false
	})
	return
}

// GetProvidersPaginated performs paginated query of providers.
func (k Keeper) GetProvidersPaginated(ctx sdk.Context, page, limit uint) (providers []v1beta1.Provider) {
	k.IterateProvidersPaginated(ctx, page, limit, func(provider v1beta1.Provider) bool {
		providers = append(providers, provider)
		return false
	})
	return
}

// IterateProvidersPaginated iterates over providers based on
// pagination parameters and performs a callback function.
func (k Keeper) IterateProvidersPaginated(ctx sdk.Context, page, limit uint, cb func(vote v1beta1.Provider) (stop bool)) {
	iterator := k.GetProvidersIteratorPaginated(ctx, page, limit)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var provider v1beta1.Provider
		k.cdc.MustUnmarshal(iterator.Value(), &provider)

		if cb(provider) {
			break
		}
	}
}

// GetProvidersIteratorPaginated returns an iterator to go over
// providers based on pagination parameters.
func (k Keeper) GetProvidersIteratorPaginated(ctx sdk.Context, page, limit uint) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIteratorPaginated(store, types.ProviderKey, page, limit)
}
