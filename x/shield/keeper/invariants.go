package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// RegisterInvariants registers all shield invariants.
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "module-account", ModuleAccountInvariant(k))
	ir.RegisterRoute(types.ModuleName, "withdraw", ProviderInvariant(k))
	ir.RegisterRoute(types.ModuleName, "withdraw", ShieldInvariant(k))
}

// ModuleAccountInvariant checks that the module account coins reflects the sum of
// remaining services and rewards held on store
func ModuleAccountInvariant(keeper Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		total := keeper.GetRemainingServiceFees(ctx)
		providers := keeper.GetAllProviders(ctx)

		for _, prov := range providers {
			total = total.Add(prov.Rewards)
		}

		moduleCoins := keeper.supplyKeeper.GetModuleAccount(ctx, types.ModuleName).GetCoins()
		totalInt, change := total.Native.TruncateDecimal()
		broken := !totalInt.IsEqual(moduleCoins) || !change.Empty()

		return sdk.FormatInvariant(types.ModuleName, "module-account",
			fmt.Sprintf("\n\tshield ModuleAccount coins: %s"+
				"\n\tsum of remaining service fees & rewards amount:  %s"+
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
				"\n\tsum of providers' collateral amount: %s\bn",
				totalWithdraw, withdrawSum, totalCollateral, collateralSum)), broken
	}
}

// ShieldInvariant checks that the providers' coin amounts equal to the tracked value.
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
				"\n\tsum of providers' withdrawing amount:  %s\n",
				totalShield, shieldSum)), broken
	}
}
