package app

import (
	"fmt"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	ica "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts"
	icacontrollertypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"

	shieldtypes "github.com/certikfoundation/shentu/v2/x/shield/types"
)

const (
	tmp = "tmp"
)

// TODO: rename upgrade title
func (app ShentuApp) setTmpUpgradeHandler() {
	app.UpgradeKeeper.SetUpgradeHandler(
		tmp,
		func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			migrationOrder := make([]string, len(fromVM))
			i := 0
			for moduleName := range fromVM {
				migrationOrder[i] = moduleName
				i++
			}
			order := module.DefaultMigrationsOrder(migrationOrder)
			// need to run crisis module last to avoid it being run before shield which has broken invariant before migration
			order = append(order, crisistypes.ModuleName)
			app.mm.SetOrderMigrations(order...)

			// create ICS27 Controller submodule params, controller module not enabled.
			controllerParams := icacontrollertypes.Params{}

			// create ICS27 Host submodule params
			hostParams := icahosttypes.Params{
				HostEnabled: true,
				AllowMessages: []string{
					sdk.MsgTypeURL(&banktypes.MsgSend{}),
					sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}),
					sdk.MsgTypeURL(&stakingtypes.MsgBeginRedelegate{}),
					sdk.MsgTypeURL(&stakingtypes.MsgCreateValidator{}),
					sdk.MsgTypeURL(&stakingtypes.MsgEditValidator{}),
					sdk.MsgTypeURL(&distrtypes.MsgWithdrawDelegatorReward{}),
					sdk.MsgTypeURL(&distrtypes.MsgSetWithdrawAddress{}),
					sdk.MsgTypeURL(&distrtypes.MsgWithdrawValidatorCommission{}),
					sdk.MsgTypeURL(&distrtypes.MsgFundCommunityPool{}),
					sdk.MsgTypeURL(&govtypes.MsgVote{}),
					sdk.MsgTypeURL(&authz.MsgExec{}),
					sdk.MsgTypeURL(&authz.MsgGrant{}),
					sdk.MsgTypeURL(&authz.MsgRevoke{}),
					sdk.MsgTypeURL(&shieldtypes.MsgPurchaseShield{}),
					sdk.MsgTypeURL(&shieldtypes.MsgWithdrawRewards{}),
					sdk.MsgTypeURL(&shieldtypes.MsgWithdrawReimbursement{}),
					sdk.MsgTypeURL(&shieldtypes.MsgDepositCollateral{}),
					sdk.MsgTypeURL(&shieldtypes.MsgWithdrawCollateral{}),
				},
			}

			// initialize ICS27 module
			icamodule, correctTypecast := app.ModuleManager().Modules[icatypes.ModuleName].(ica.AppModule)
			if !correctTypecast {
				panic("mm.Modules[icatypes.ModuleName] is not of type ica.AppModule")
			}

			icamodule.InitModule(ctx, controllerParams, hostParams)

			ctx.Logger().Info("Start to run module migrations...")
			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if upgradeInfo.Name == tmp && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{crisistypes.ModuleName},
		}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}
