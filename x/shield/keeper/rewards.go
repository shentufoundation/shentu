package keeper

import (
	"time"

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
	rewards.Native = sdk.DecCoins{}
	k.SetRewards(ctx, addr, rewards)

	// Add leftovers as service fees.
	serviceFees := k.GetServiceFees(ctx)
	serviceFees.Native = serviceFees.Native.Add(change...)
	k.SetServiceFees(ctx, serviceFees)

	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, ctkRewards); err != nil {
		return sdk.Coins{}, err
	}
	return ctkRewards, nil
}

// DistributeFees distributes service fees to all providers.
func (k Keeper) DistributeFees(ctx sdk.Context) {
	secondsFromLastDistribution := sdk.NewDecFromInt(sdk.NewInt(int64(ctx.BlockTime().Sub(ctx.WithBlockHeight(ctx.BlockHeight() - 1).BlockTime()).Seconds())))
	serviceFeesPerSecond := k.GetServiceFeesPerSecond(ctx)
	serviceFees := serviceFeesPerSecond.MulDec(secondsFromLastDistribution)
	remainingServiceFees := k.GetServiceFees(ctx)
	bondDenom := k.BondDenom(ctx)
	if remainingServiceFees.Native.AmountOf(bondDenom).LT(serviceFees.Native.AmountOf(bondDenom)) {
		serviceFees.Native = remainingServiceFees.Native
		serviceFeesPerSecond.Native = sdk.DecCoins{}
		k.SetServiceFeesPerSecond(ctx, serviceFeesPerSecond)
	}

	totalCollateral := k.GetTotalCollateral(ctx)

	providers := k.GetAllProviders(ctx)
	for _, provider := range providers {
		proportion := sdk.NewDecFromInt(sdk.MaxInt(provider.Collateral, sdk.ZeroInt())).QuoInt(totalCollateral)
		nativeFees := serviceFees.Native.MulDec(proportion)

		remainingServiceFees.Native = remainingServiceFees.Native.Sub(nativeFees)

		provider.Rewards = provider.Rewards.Add(types.MixedDecCoins{Native: nativeFees})
		k.SetProvider(ctx, provider.Address, provider)
	}
	k.SetServiceFees(ctx, remainingServiceFees)
}

// UpdateServiceFeesPerSecond update service fees per second based on purchases' completion time.
func (k Keeper) UpdateServiceFeesPerSecond(ctx sdk.Context) {
	serviceFeesPerSecond := k.GetServiceFeesPerSecond(ctx)
	if !serviceFeesPerSecond.Native.IsAllPositive() {
		return
	}

	deletionPeriod := k.GetPurchaseDeletionPeriod(ctx)
	startTime := ctx.WithBlockHeight(ctx.BlockHeight() - 1).BlockTime().Add(deletionPeriod)
	endTime := ctx.BlockTime().Add(deletionPeriod)
	// (startTime, endTime]
	iterator := ctx.KVStore(k.storeKey).Iterator(types.GetPurchaseCompletionTimeKey(startTime.Add(time.Nanosecond)), sdk.InclusiveEndBytes(types.GetPurchaseCompletionTimeKey(endTime)))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var timeSlice []types.PoolPurchaser
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &timeSlice)
		for _, poolPurchaser := range timeSlice {
			purchaseList, _ := k.GetPurchaseList(ctx, poolPurchaser.PoolID, poolPurchaser.Purchaser)
			for _, entry := range purchaseList.Entries {
				if entry.ProtectionEndTime.After(ctx.WithBlockHeight(ctx.BlockHeight()-1).BlockTime()) && !entry.ProtectionEndTime.After(ctx.BlockTime()) {
					if !serviceFeesPerSecond.Native.IsAllPositive() {
						serviceFeesPerSecond.Native = sdk.DecCoins{}
						k.SetServiceFeesPerSecond(ctx, serviceFeesPerSecond)
						return
					}
					serviceFeesPerSecond = serviceFeesPerSecond.Sub(entry.ServiceFeesPerSecond)
				}
			}
		}
	}

	k.SetServiceFeesPerSecond(ctx, serviceFeesPerSecond)
}
