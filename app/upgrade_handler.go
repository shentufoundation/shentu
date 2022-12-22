package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

const (
	upgradeName = "v2.7.0"
)

func (app ShentuApp) setUpgradeHandler() {
	app.UpgradeKeeper.SetUpgradeHandler(
		upgradeName,
		func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			// don't run Initgenesis since it'll be set with a wrong denom
			fromVM[gov.ModuleName] = app.mm.Modules[gov.ModuleName].ConsensusVersion()
			// don't run icamodule's Initgenesis since it'll overwrite the icahost params that be set here
			// the InitModule will be called later on to set params.
			// this assumes it's the first time ica module go into fromV

			ctx.Logger().Info("Start to run module migrations...")
			newVersionMap, err := app.mm.RunMigrations(ctx, app.configurator, fromVM)

			ctx.Logger().Info("Fixing Shield invariant...")

			return newVersionMap, err
		},
	)
}
