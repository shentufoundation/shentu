package keeper

import (
	types "github.com/certikfoundation/shentu/v2/x/shield/types/v1beta1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PayoutNativeRewards pays out pending CTK rewards.
func (k Keeper) PayoutNativeRewards(ctx sdk.Context, addr sdk.AccAddress) (sdk.Coins, error) {
	provider, found := k.GetProvider(ctx, addr)
	if !found {
		return sdk.Coins{}, types.ErrProviderNotFound
	}

	ctkRewards, change := provider.Rewards.TruncateDecimal()
	if ctkRewards.IsZero() {
		return nil, nil
	}
	provider.Rewards = sdk.DecCoins{}
	providerAddr, err := sdk.AccAddressFromBech32(provider.Address)
	if err != nil {
		panic(err)
	}
	k.SetProvider(ctx, providerAddr, provider)

	// Add leftovers as fees
	remainingFees := k.GetRemainingFees(ctx)
	remainingFees = remainingFees.Add(change...)
	k.SetRemainingFees(ctx, remainingFees)

	if err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, ctkRewards); err != nil {
		return sdk.Coins{}, err
	}
	return ctkRewards, nil
}
