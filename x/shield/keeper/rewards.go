package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/x/shield/types"
)

// PayoutNativeRewards pays out pending CTK rewards.
func (k Keeper) PayoutNativeRewards(ctx sdk.Context, addr sdk.AccAddress) (sdk.Coins, error) {
	provider, found := k.GetProvider(ctx, addr)
	if !found {
		return sdk.Coins{}, types.ErrProviderNotFound
	}

	ctkRewards, change := provider.NativeReward.TruncateDecimal()
	if ctkRewards.IsZero() {
		return nil, nil
	}
	provider.NativeReward = sdk.DecCoins{}
	providerAddr, err := sdk.AccAddressFromBech32(provider.Address)
	if err != nil {
		panic(err)
	}
	k.SetProvider(ctx, providerAddr, provider)

	// Add leftovers as service fees.
	remainingNativeServiceFee := k.GetRemainingNativeServiceFee(ctx)
	remainingNativeServiceFee = remainingNativeServiceFee.Add(change...)
	k.SetRemainingNativeServiceFee(ctx, remainingNativeServiceFee)

	if err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, ctkRewards); err != nil {
		return sdk.Coins{}, err
	}
	return ctkRewards, nil
}
