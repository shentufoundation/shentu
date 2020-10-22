package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// RegisterInvariants registers all shield invariants.
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "module-account", ModuleAccountInvariant(k))
}

// ModuleAccountInvariant checks that the module account coins reflects the sum of
// deposit amounts held on store
func ModuleAccountInvariant(keeper Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		total := keeper.GetServiceFees(ctx)
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
