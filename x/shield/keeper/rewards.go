package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// GetPendingPayouts gets pending payouts for a denom.
func (k Keeper) GetPendingPayouts(ctx sdk.Context, denom string) types.PendingPayouts {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPendingPayoutsKey(denom))
	if bz == nil {
		return nil
	}
	var pending types.PendingPayouts
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &pending)
	return pending
}

// SetPendingPayouts sets pending payouts for a denom.
func (k Keeper) SetPendingPayouts(ctx sdk.Context, denom string, payout types.PendingPayouts) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(payout)
	store.Set(types.GetPendingPayoutsKey(denom), bz)
}

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

// AddPendingPayout appends a pending payment to pending payouts for a denomination.
func (k Keeper) AddPendingPayout(ctx sdk.Context, denom string, payout types.PendingPayout) {
	payouts := k.GetPendingPayouts(ctx, denom)
	if payouts == nil {
		k.SetPendingPayouts(ctx, denom, types.PendingPayouts{payout})
		return
	}
	payouts = append(payouts, payout)
	k.SetPendingPayouts(ctx, denom, payouts)
}

// PayoutNativeRewards pays out pending CTK rewards.
func (k Keeper) PayoutNativeRewards(ctx sdk.Context, addr sdk.AccAddress) (sdk.Coins, error) {
	rewards := k.GetRewards(ctx, addr)
	ctkRewards, change := rewards.Native.TruncateDecimal()
	rewards.Native = change
	k.SetRewards(ctx, addr, rewards)
	err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, ctkRewards)
	if err != nil {
		return sdk.Coins{}, err
	}
	return ctkRewards, nil
}
