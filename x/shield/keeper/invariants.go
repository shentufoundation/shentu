package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// RegisterInvariants registers all shield invariants.
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "module-account", ModuleAccountInvariant(k))
	ir.RegisterRoute(types.ModuleName, "provider", ProviderInvariant(k))
	ir.RegisterRoute(types.ModuleName, "shield", ShieldInvariant(k))
	ir.RegisterRoute(types.ModuleName, "global-staking-pool", GlobalStakingPoolInvariant(k))
	ir.RegisterRoute(types.ModuleName, "original-global-staking", StakingForShieldPurchaseInvariant(k))
}

// ModuleAccountInvariant checks that the module account coins reflects the sum of
// remaining services and rewards held on store
func ModuleAccountInvariant(keeper Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		moduleCoins := keeper.supplyKeeper.GetModuleAccount(ctx, types.ModuleName).GetCoins()
		total := keeper.GetRemainingServiceFees(ctx)
		providers := keeper.GetAllProviders(ctx)

		for _, prov := range providers {
			total = total.Add(prov.Rewards)
		}

		bondDenom := keeper.BondDenom(ctx)
		totalInt, change := total.Native.TruncateDecimal()
		stakedCoin := sdk.NewCoin(bondDenom, sdk.ZeroInt())
		for _, staked := range keeper.GetAllStakeForShields(ctx) {
			stakedCoin = stakedCoin.Add(sdk.NewCoin(bondDenom, staked.Amount))
		}
		totalInt = totalInt.Add(stakedCoin)

		blockServiceFees := keeper.GetBlockServiceFees(ctx).Native.AmountOf(bondDenom).TruncateInt()
		blockFeesCoin := sdk.NewCoin(bondDenom, blockServiceFees)
		totalInt = totalInt.Add(blockFeesCoin)

		for _, rmb := range keeper.GetAllReimbursements(ctx) {
			totalInt = totalInt.Add(sdk.NewCoin(bondDenom, rmb.Amount.AmountOf(bondDenom)))
		}

		broken := !totalInt.IsEqual(moduleCoins) || !change.Empty()

		return sdk.FormatInvariant(types.ModuleName, "module-account",
			fmt.Sprintf("\n\tshield ModuleAccount coins: %s"+
				"\n\tsum of remaining service fees & rewards & staked & reimbursement amount:  %s"+
				"\n\tremaining change amount: %s\n",
				moduleCoins, totalInt, change)), broken
	}
}

// ProviderInvariant checks that the providers' coin amounts equal to the tracked value.
func ProviderInvariant(keeper Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		providers := keeper.GetAllProviders(ctx)
		withdrawSum := sdk.NewInt(0)
		collateralSum := sdk.NewInt(0)
		for _, prov := range providers {
			withdrawSum = withdrawSum.Add(prov.Withdrawing)
			collateralSum = collateralSum.Add(prov.Collateral)
		}

		totalWithdraw := keeper.GetTotalWithdrawing(ctx)
		totalCollateral := keeper.GetTotalCollateral(ctx)
		broken := !totalWithdraw.Equal(withdrawSum) || !totalCollateral.Equal(collateralSum)

		return sdk.FormatInvariant(types.ModuleName, "provider",
			fmt.Sprintf("\n\ttotal withdraw amount: %s"+
				"\n\tsum of providers' withdrawing amount:  %s"+
				"\n\ttotal collateral amount: %s"+
				"\n\tsum of providers' collateral amount: %s\n",
				totalWithdraw, withdrawSum, totalCollateral, collateralSum)), broken
	}
}

// ShieldInvariant checks that the sum of individual pools' shield is
// equal to the total shield.
func ShieldInvariant(keeper Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		pools := keeper.GetAllPools(ctx)
		shieldSum := sdk.NewInt(0)
		for _, pool := range pools {
			shieldSum = shieldSum.Add(pool.Shield)
		}

		totalShield := keeper.GetTotalShield(ctx)
		broken := !totalShield.Equal(shieldSum)

		return sdk.FormatInvariant(types.ModuleName, "shield",
			fmt.Sprintf("\n\ttotal shield amount: %s"+
				"\n\tsum of pools' shield amount:  %s\n",
				totalShield, shieldSum)), broken
	}
}

// GlobalStakingPoolInvariant checks the total staked sum equals to the global staking pool amount.
func GlobalStakingPoolInvariant(keeper Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		stakedCoin := sdk.NewCoin(keeper.BondDenom(ctx), sdk.ZeroInt())
		for _, staked := range keeper.GetAllStakeForShields(ctx) {
			stakedCoin = stakedCoin.Add(sdk.NewCoin(keeper.BondDenom(ctx), staked.Amount))
		}
		stakedInt := stakedCoin.Amount
		globalStakingPool := keeper.GetGlobalShieldStakingPool(ctx)
		broken := !stakedInt.Equal(globalStakingPool)

		return sdk.FormatInvariant(types.ModuleName, "global-staking-pool",
			fmt.Sprintf("\n\tsum of staked amount:  %s"+
				"\n\tglobal staking pool amount: %s\n",
				stakedInt, globalStakingPool.String())), broken
	}
}

// SFSPurchaseInvariant checks that sum of original staked shield equals to the total.
func StakingForShieldPurchaseInvariant(keeper Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		globalStakingPool := keeper.GetGlobalShieldStakingPool(ctx)

		sum := sdk.ZeroInt()
		for _, os := range keeper.GetAllOriginalStakings(ctx) {
			sum = sum.Add(os.Amount)
		}

		broken := !globalStakingPool.Equal(sum)

		return sdk.FormatInvariant(types.ModuleName, "global-staking-pool",
			fmt.Sprintf("\n\tsum of originally staked amount:  %s"+
				"\n\tglobal staking pool amount: %s\n",
				sum, globalStakingPool.String())), broken
	}
}
