package bounty

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shentufoundation/shentu/v2/x/bounty/keeper"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	// TODO: implement beginblocker
}

// check all registered invariants
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// TODO: implement endblocker
}
