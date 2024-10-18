package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/shield/types"
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
		return nil, err
	}
	if err = k.SetProvider(ctx, providerAddr, provider); err != nil {
		return nil, err
	}

	// Add leftovers as service fees
	remainingServiceFees, err := k.GetRemainingServiceFees(ctx)
	if err != nil {
		return nil, err
	}
	remainingServiceFees = remainingServiceFees.Add(change...)
	if err = k.SetRemainingServiceFees(ctx, remainingServiceFees); err != nil {
		return nil, err
	}

	if err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, ctkRewards); err != nil {
		return sdk.Coins{}, err
	}
	return ctkRewards, nil
}
