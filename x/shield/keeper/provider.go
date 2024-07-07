package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/shield/types"
)

// SetProvider sets data of a provider in the kv-store.
func (k Keeper) SetProvider(ctx sdk.Context, delAddr sdk.AccAddress, provider types.Provider) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&provider)
	store.Set(types.GetProviderKey(delAddr), bz)
}

// GetProvider returns data of a provider given its address.
func (k Keeper) GetProvider(ctx sdk.Context, delegator sdk.AccAddress) (dt types.Provider, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetProviderKey(delegator))
	if bz == nil {
		return types.Provider{}, false
	}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &dt)
	return dt, true
}

// IterateProviders iterates through all providers.
func (k Keeper) IterateProviders(ctx sdk.Context, callback func(provider types.Provider) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ProviderKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var provider types.Provider
		k.cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &provider)

		if callback(provider) {
			break
		}
	}
}

// GetAllProviders retrieves all providers.
func (k Keeper) GetAllProviders(ctx sdk.Context) (providers []types.Provider) {
	k.IterateProviders(ctx, func(provider types.Provider) bool {
		providers = append(providers, provider)
		return false
	})
	return
}

// GetProvidersPaginated performs paginated query of providers.
func (k Keeper) GetProvidersPaginated(ctx sdk.Context, page, limit uint) (providers []types.Provider) {
	k.IterateProvidersPaginated(ctx, page, limit, func(provider types.Provider) bool {
		providers = append(providers, provider)
		return false
	})
	return
}

// IterateProvidersPaginated iterates over providers based on
// pagination parameters and performs a callback function.
func (k Keeper) IterateProvidersPaginated(ctx sdk.Context, page, limit uint, cb func(vote types.Provider) (stop bool)) {
	iterator := k.GetProvidersIteratorPaginated(ctx, page, limit)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var provider types.Provider
		k.cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &provider)

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
