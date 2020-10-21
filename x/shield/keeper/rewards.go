package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// GetRewards returns total rewards for an address.
func (k Keeper) GetRewards(ctx sdk.Context, addr sdk.AccAddress) types.MixedDecCoins {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetProviderKey(addr))
	if bz == nil {
		return types.InitMixedDecCoins()
	}
	var provider types.Provider
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &provider)
	return provider.Rewards
}

// SetRewards sets the rewards for an address.
func (k Keeper) SetRewards(ctx sdk.Context, addr sdk.AccAddress, earnings types.MixedDecCoins) {
	provider, found := k.GetProvider(ctx, addr)
	if !found {
		provider = types.NewProvider(addr)
	}
	provider.Rewards = earnings
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(provider)
	store.Set(types.GetProviderKey(addr), bz)
}

// AddRewards adds coins to earned rewards.
func (k Keeper) AddRewards(ctx sdk.Context, provider sdk.AccAddress, earnings types.MixedDecCoins) {
	rewards := k.GetRewards(ctx, provider)
	rewards = rewards.Add(earnings)
	k.SetRewards(ctx, provider, rewards)
}

// PayoutNativeRewards pays out pending CTK rewards.
func (k Keeper) PayoutNativeRewards(ctx sdk.Context, addr sdk.AccAddress) (sdk.Coins, error) {
	rewards := k.GetRewards(ctx, addr)
	ctkRewards, change := rewards.Native.TruncateDecimal()
	if ctkRewards.IsZero() {
		return nil, nil
	}
	rewards.Native = change
	k.SetRewards(ctx, addr, rewards)
	err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, ctkRewards)
	if err != nil {
		return sdk.Coins{}, err
	}
	return ctkRewards, nil
}

// DistributeFees distributes service fees to all providers.
func (k Keeper) DistributeFees(ctx sdk.Context) {
	serviceFees := k.GetServiceFees(ctx)

	totalCollateral := k.GetTotalCollateral(ctx)
	totalLocked := k.GetTotalLocked(ctx)
	totalCollateralAmount := totalCollateral.Add(totalLocked)

	providers := k.GetAllProviders(ctx)
	for _, provider := range providers {
		proportion := sdk.NewDecFromInt(sdk.MaxInt(provider.Collateral.Add(provider.Locked), sdk.ZeroInt())).QuoInt(totalCollateralAmount)
		nativeFees := serviceFees.Native.MulDecTruncate(proportion)
		foreignFees := serviceFees.Foreign.MulDecTruncate(proportion)

		serviceFees.Native = serviceFees.Native.Sub(nativeFees)
		serviceFees.Foreign = serviceFees.Foreign.Sub(foreignFees)

		rewards := types.NewMixedDecCoins(nativeFees, foreignFees)
		provider.Rewards = provider.Rewards.Add(rewards)
		k.SetProvider(ctx, provider.Address, provider)
	}

	k.SetServiceFees(ctx, serviceFees)
}
