package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/x/shield/types"
)

// DistributeShieldRewards distributes Shield Rewards to
// collateral providers.
func (k Keeper) DistributeShieldRewards(ctx sdk.Context) {
	// Add block service fees that need to be distributed for this block
	blockServiceFees := k.GetBlockServiceFees(ctx)

	// FIX DIS
	// serviceFees = serviceFees.Add(blockServiceFees)
	serviceFees := blockServiceFees
	remainingServiceFees := blockServiceFees
	k.DeleteBlockServiceFees(ctx)

	// Distribute service fees.
	totalCollateral := k.GetTotalCollateral(ctx)
	providers := k.GetAllProviders(ctx)
	bondDenom := k.sk.BondDenom(ctx)
	for _, provider := range providers {
		providerAddr, err := sdk.AccAddressFromBech32(provider.Address)
		if err != nil {
			panic(err)
		}

		// fees * providerCollateral / totalCollateral
		nativeFees := serviceFees.MulDec(sdk.NewDecFromInt(provider.Collateral).QuoInt(totalCollateral))
		if nativeFees.AmountOf(bondDenom).GT(remainingServiceFees.AmountOf(bondDenom)) {
			nativeFees = remainingServiceFees
		}
		provider.Rewards = provider.Rewards.Add(nativeFees...)
		k.SetProvider(ctx, providerAddr, provider)

		remainingServiceFees = remainingServiceFees.Sub(nativeFees)
	}
	// add back block service fees
	remainingServiceFees = remainingServiceFees.Add(blockServiceFees...)
	k.SetRemainingServiceFees(ctx, remainingServiceFees)
}

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

	// Add leftovers as service fees.
	remainingServiceFees := k.GetRemainingServiceFees(ctx)
	remainingServiceFees = remainingServiceFees.Add(change...)
	k.SetRemainingServiceFees(ctx, remainingServiceFees)

	if err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, ctkRewards); err != nil {
		return sdk.Coins{}, err
	}
	return ctkRewards, nil
}
