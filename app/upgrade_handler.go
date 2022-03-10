package app

import (
	"fmt"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	sdkauthz "github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	sdkfeegrant "github.com/cosmos/cosmos-sdk/x/feegrant"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	ibcconnectiontypes "github.com/cosmos/ibc-go/v2/modules/core/03-connection/types"
)

const upgradeName = "Shentu-v230"

func (app ShentuApp) setUpgradeHandler() {
	app.UpgradeKeeper.SetUpgradeHandler(
		upgradeName,
		func(ctx sdk.Context, _ upgradetypes.Plan, _ module.VersionMap) (module.VersionMap, error) {
			app.IbcKeeper.ConnectionKeeper.SetParams(ctx, ibcconnectiontypes.DefaultParams())

			fromVM := make(map[string]uint64)
			for moduleName := range app.mm.Modules {
				fromVM[moduleName] = 1
			}
			// override versions for _new_ modules as to not skip InitGenesis
			fromVM[sdkauthz.ModuleName] = 0
			fromVM[sdkfeegrant.ModuleName] = 0

			fromVM[authtypes.ModuleName] = 2
			newVM, err := app.mm.RunMigrations(ctx, app.configurator, fromVM)

			if err != nil {
				return newVM, err
			}

			newVM[authtypes.ModuleName] = 1

			return app.mm.RunMigrations(ctx, app.configurator, newVM)
		},
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if upgradeInfo.Name == upgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{authz.ModuleName, feegrant.ModuleName},
		}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}
