package app

import (
	"fmt"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	sdkauthz "github.com/cosmos/cosmos-sdk/x/authz"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	sdkfeegrant "github.com/cosmos/cosmos-sdk/x/feegrant"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	ibcconnectiontypes "github.com/cosmos/ibc-go/v2/modules/core/03-connection/types"

	v231 "github.com/certikfoundation/shentu/v2/app/v231"
	shieldtypes "github.com/certikfoundation/shentu/v2/x/shield/types"
)

const (
	v230Upgrade = "Shentu-v230"
	shieldv2    = "Shield-V2"
)

func (app ShentuApp) setv230UpgradeHandler() {
	app.UpgradeKeeper.SetUpgradeHandler(
		v230Upgrade,
		func(ctx sdk.Context, _ upgradetypes.Plan, _ module.VersionMap) (module.VersionMap, error) {
			app.IBCKeeper.ConnectionKeeper.SetParams(ctx, ibcconnectiontypes.DefaultParams())

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

	if upgradeInfo.Name == v230Upgrade && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{authz.ModuleName, feegrant.ModuleName},
		}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}

func (app ShentuApp) setShieldV2UpgradeHandler() {
	app.UpgradeKeeper.SetUpgradeHandler(
		shieldv2,
		func(ctx sdk.Context, _ upgradetypes.Plan, _ module.VersionMap) (module.VersionMap, error) {
			// Refund v1 purchases
			v231.RefundPurchasers(ctx, app.appCodec, app.BankKeeper, &app.StakingKeeper, app.ShieldKeeper, app.keys[shieldtypes.StoreKey])

			// Payout reimbursements
			v231.PayoutReimbursements(ctx, app.appCodec, app.BankKeeper, app.ShieldKeeper, app.keys[shieldtypes.StoreKey])

			fromVM := make(map[string]uint64)
			for moduleName := range app.mm.Modules {
				fromVM[moduleName] = 2
			}

			fromVM[shieldtypes.ModuleName] = 1
			fromVM[govtypes.ModuleName] = 2

			// Delete crisis from fromVM to trigger crisis genesis init
			delete(fromVM, crisistypes.ModuleName)

			// Modify migration order. Put crisis behind shield so that shield gets migrated before crisis invariant check
			allOtherModules := make([]string, len(fromVM))
			i := 0
			for moduleName := range fromVM {
				allOtherModules[i] = moduleName
				i++
			}
			order := module.DefaultMigrationsOrder(allOtherModules)
			order = append(order, crisistypes.ModuleName)
			app.mm.SetOrderMigrations(order...)

			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if upgradeInfo.Name == shieldv2 && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{crisistypes.ModuleName},
		}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}
