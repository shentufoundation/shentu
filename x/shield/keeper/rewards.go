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

	ctkRewards, change := provider.Rewards.Native.TruncateDecimal()
	if ctkRewards.IsZero() {
		return nil, nil
	}
	provider.Rewards.Native = sdk.DecCoins{}
	providerAddr, err := sdk.AccAddressFromBech32(provider.Address)
	if err != nil {
		panic(err)
	}
	k.SetProvider(ctx, providerAddr, provider)

	// Add leftovers as service fees.
	remainingServiceFees := k.GetRemainingServiceFees(ctx)
	remainingServiceFees.Native = remainingServiceFees.Native.Add(change...)
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
	totalBondedTokens := k.bk.GetAllBalances(ctx, k.sk.GetBondedPool(ctx).GetAddress()).AmountOf(k.BondDenom(ctx)).ToDec() // c + n
	totalShieldDeposit := k.GetGlobalShieldStakingPool(ctx).ToDec()                                                        // d

	var leverage sdk.Dec // l = (total shield) / (total collateral)
	if totalCollateral.IsZero() {
		leverage = sdk.ZeroDec()
	} else {
		leverage = totalShield.ToDec().Quo(totalCollateral.ToDec())
	}

	blockRewardParams := k.GetDistributionParams(ctx)
	modelParamA := blockRewardParams.ModelParamA    // a
	modelParamB := blockRewardParams.ModelParamB    // b
	targetLeverage := blockRewardParams.MaxLeverage // L

	/* The non-linear model:
	 *         c+n                        l        d
	 *   r = -------- ( a + 2(b - a) * ------- + -----)
	 *        c+n+d                     l + L     c+n
	 */
	if leverage.Add(targetLeverage).IsZero() || totalBondedTokens.IsZero() {
		return sdk.ZeroDec()
	} else {
		leading := totalBondedTokens.Quo(totalBondedTokens.Add(totalShieldDeposit))                        // (c+n)/(c+n+d)
		first := modelParamA                                                                               // a
		second := leverage.Quo(leverage.Add(targetLeverage)).Mul(modelParamB.Sub(modelParamA).MulInt64(2)) // 2(b-a)(l/l+L)
		third := totalShieldDeposit.Quo(totalBondedTokens)                                                 // d/(c+n)
		inner := first.Add(second).Add(third)
		return leading.Mul(inner)
	}
}
