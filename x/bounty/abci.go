package bounty

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/crisis/keeper"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	// TODO: implement beginblocker
}

// check all registered invariants
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// TODO: implement endblocker
}
