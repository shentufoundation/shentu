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
	remainingServiceFees := blockServiceFees
	k.DeleteBlockServiceFees(ctx)

	// TODO: Add support for any denoms.

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
		fees := blockServiceFees.MulDec(sdk.NewDecFromInt(provider.Collateral).QuoInt(totalCollateral))
		if fees.AmountOf(bondDenom).GT(remainingServiceFees.AmountOf(bondDenom)) {
			fees = remainingServiceFees
		}
		provider.Rewards = provider.Rewards.Add(fees...)
		k.SetProvider(ctx, providerAddr, provider)

		remainingServiceFees = remainingServiceFees.Sub(fees)
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

// GetShieldBlockRewardRatio calculates the dynamic ratio for block rewards to shield module, based on total shield and total collateral.
func (k Keeper) GetShieldBlockRewardRatio(ctx sdk.Context) sdk.Dec {
	totalShield := k.GetTotalShield(ctx)
	totalCollateral := k.GetTotalCollateral(ctx)

	var leverage sdk.Dec // l = (total shield) / (total collateral)
	if totalCollateral.IsZero() {
		leverage = sdk.ZeroDec()
	} else {
		leverage = totalShield.ToDec().Quo(totalCollateral.ToDec())
	}

	blockRewardParams := k.GetBlockRewardParams(ctx)
	modelParamA := blockRewardParams.ModelParamA       // a
	modelParamB := blockRewardParams.ModelParamB       // b
	targetLeverage := blockRewardParams.TargetLeverage // L

	/* The non-linear model:
	 *                         l
	 *   r = a + 2(b - a) * -------
	 *                       l + L
	 */
	return leverage.Quo(leverage.Add(targetLeverage)).Mul(modelParamB.Sub(modelParamA).MulInt64(2)).Add(modelParamA)
}
