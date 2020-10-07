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
func (k Keeper) addProvider(ctx sdk.Context, addr sdk.AccAddress) {
	delegations := k.sk.GetAllDelegatorDelegations(ctx, addr)

	bondDenom := k.sk.BondDenom(ctx)
	totalStaked := sdk.Coins{}
	for _, del := range delegations {
		val, found := k.sk.GetValidator(ctx, del.GetValidatorAddr())
		if !found {
			panic("expected validator, not found")
		}
		totalStaked = totalStaked.Add(sdk.NewCoin(bondDenom, val.TokensFromShares(del.GetShares()).TruncateInt()))
	}

	provider := types.NewProvider(addr)
	provider.DelegationBonded = totalStaked
	provider.Available = totalStaked.AmountOf(bondDenom)

	k.SetProvider(ctx, addr, provider)
}

func (k Keeper) UpdateDelegationAmount(ctx sdk.Context, delAddr sdk.AccAddress) {
	// go through delAddr's delegations to recompute total amount of bonded delegation
	// update or create a new entry
	provider, found := k.GetProvider(ctx, delAddr)
	if !found {
		return // ignore non-participating addr
	}

	// update delegations
	totalStakedAmount := sdk.NewInt(0)
	delegations := k.sk.GetAllDelegatorDelegations(ctx, delAddr)
	for _, del := range delegations {
		val, found := k.sk.GetValidator(ctx, del.GetValidatorAddr())
		if !found {
			panic("expected validator, not found")
		}
		totalStakedAmount = totalStakedAmount.Add(val.TokensFromShares(del.GetShares()).TruncateInt())
	}

	bondDenom := k.sk.BondDenom(ctx)
	deltaAmount := totalStakedAmount.Sub(provider.DelegationBonded.AmountOf(bondDenom))
	provider.DelegationBonded = sdk.NewCoins(sdk.NewCoin(bondDenom, totalStakedAmount))
	withdrawalAmount := sdk.NewInt(0)
	if deltaAmount.IsNegative() {
		if totalStakedAmount.LT(provider.Collateral.AmountOf(bondDenom).Sub(provider.Withdrawal)) {
			withdrawalAmount = provider.Collateral.AmountOf(bondDenom).Sub(provider.Withdrawal).Sub(totalStakedAmount)
		}
		provider.Available = provider.Available.Sub(deltaAmount.Neg())
	} else {
		provider.Available = provider.Available.Add(deltaAmount)
	}
	k.SetProvider(ctx, delAddr, provider)

	// save the change of provider before this because withdrawal also updates the provider
	if withdrawalAmount.IsPositive() {
		k.WithdrawFromPools(ctx, delAddr, sdk.NewCoins(sdk.NewCoin(bondDenom, withdrawalAmount)))
	}
}

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

	bondDenom := k.sk.BondDenom(ctx)
	provider.DelegationBonded = provider.DelegationBonded.Sub(sdk.NewCoins(sdk.NewCoin(bondDenom, deltaAmount)))
	withdrawalAmount := sdk.NewInt(0)
	if deltaAmount.IsNegative() {
		if provider.DelegationBonded.AmountOf(bondDenom).LT(
			provider.Collateral.AmountOf(bondDenom).Sub(provider.Withdrawal),
		) {
			withdrawalAmount = provider.Collateral.AmountOf(bondDenom).Sub(
				provider.Withdrawal).Sub(provider.DelegationBonded.AmountOf(bondDenom))
		}
		provider.Available = provider.Available.Sub(deltaAmount.Neg())
	} else {
		provider.Available = provider.Available.Add(deltaAmount)
	}
	k.SetProvider(ctx, delAddr, provider)

	// note that this will also be triggered by redelegations
	if withdrawalAmount.IsPositive() {
		k.WithdrawFromPools(ctx, delAddr, sdk.NewCoins(sdk.NewCoin(bondDenom, withdrawalAmount)))
	}
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
