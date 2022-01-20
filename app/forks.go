package app

import (
	"github.com/certikfoundation/shentu/v2/app/mva"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlockForks is intended to be ran in
func BeginBlockForks(ctx sdk.Context, app *ShentuApp) {
	switch ctx.BlockHeight() {
	case mva.UpgradeHeight:
		mva.RunForkLogic(ctx, &app.accountKeeper, app.bankKeeper, &app.stakingKeeper)
	default:
		// do nothing
		return
	}
}
