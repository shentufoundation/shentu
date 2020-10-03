package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

func (k Keeper) SetProvider(ctx sdk.Context, delAddr sdk.AccAddress, provider types.Provider) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(provider)
	store.Set(types.GetProviderKey(delAddr), bz)
}

func (k Keeper) GetProvider(ctx sdk.Context, delegator sdk.AccAddress) (dt types.Provider, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetProviderKey(delegator))
	if bz == nil {
		return types.Provider{}, false
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &dt)
	return dt, true
}

// addProvider adds a new provider into shield module. Should only be called from DepositCollateral.
func (k Keeper) addProvider(ctx sdk.Context, addr sdk.AccAddress) types.Provider {
	delegations := k.sk.GetAllDelegatorDelegations(ctx, addr)

	totalStaked := sdk.Coins{}
	for _, del := range delegations {
		val, found := k.sk.GetValidator(ctx, del.GetValidatorAddr())
		if !found {
			panic("expected validator, not found")
		}
		totalStaked = totalStaked.Add(sdk.NewCoin(k.sk.BondDenom(ctx), val.TokensFromShares(del.GetShares()).TruncateInt()))
	}

	provider := types.NewProvider(addr)
	provider.DelegationBonded = totalStaked

	k.SetProvider(ctx, addr, provider)
	return provider
}

func (k Keeper) UpdateDelegationAmount(ctx sdk.Context, delAddr sdk.AccAddress) {
	// go through delAddr's delegations to recompute total amount of bonded delegation
	// update or create a new entry
	provider, found := k.GetProvider(ctx, delAddr)
	if !found {
		return // ignore non-participating addr
	}

	// update delegations
	totalStaked := sdk.Coins{}
	delegations := k.sk.GetAllDelegatorDelegations(ctx, delAddr)
	for _, del := range delegations {
		val, found := k.sk.GetValidator(ctx, del.GetValidatorAddr())
		if !found {
			panic("expected validator, not found")
		}
		totalStaked = totalStaked.Add(sdk.NewCoin(k.sk.BondDenom(ctx), val.TokensFromShares(del.GetShares()).TruncateInt()))
	}

	provider.DelegationBonded = totalStaked

	if provider.DelegationBonded.IsAllLT(provider.Collateral) {
		withdrawAmount := provider.Collateral.Sub(provider.DelegationBonded)
		k.WithdrawFromPools(ctx, delAddr, withdrawAmount)
	}
	k.SetProvider(ctx, delAddr, provider)
}

// IterateProviders iterates through all providers
func (k Keeper) IterateProviders(ctx sdk.Context, callback func(provider types.Provider) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ProviderKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var provider types.Provider
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &provider)

		if callback(provider) {
			break
		}
	}
}

// GetAllProviders retrieves all providres.
func (k Keeper) GetAllProviders(ctx sdk.Context) (providers []types.Provider) {
	k.IterateProviders(ctx, func(provider types.Provider) bool {
		providers = append(providers, provider)
		return false
	})
	return
}
