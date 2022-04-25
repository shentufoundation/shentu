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
	remainingServiceFees := k.GetServiceFees(ctx)
	remainingServiceFees = remainingServiceFees.Add(change...)
	k.SetServiceFees(ctx, remainingServiceFees)

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
	if leverage.Add(targetLeverage).IsZero() {
		return sdk.ZeroDec()
	} else {
		return leverage.Quo(leverage.Add(targetLeverage)).Mul(modelParamB.Sub(modelParamA).MulInt64(2)).Add(modelParamA)
	}
}
