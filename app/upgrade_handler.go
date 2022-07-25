package app

import (
	"fmt"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	sdkauthz "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	sdkfeegrant "github.com/cosmos/cosmos-sdk/x/feegrant"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	ibcconnectiontypes "github.com/cosmos/ibc-go/v2/modules/core/03-connection/types"

	shieldtypes "github.com/certikfoundation/shentu/v2/x/shield/types"
)

const (
	v230Upgrade = "Shentu-v230"
	shieldv2    = "Shield-V2"
	tmp         = "tmp"
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
			fromVM := make(map[string]uint64)
			for moduleName := range app.mm.Modules {
				fromVM[moduleName] = 2
			}

			fromVM[shieldtypes.ModuleName] = 1
			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if upgradeInfo.Name == shieldv2 && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}

// TODO: rename upgrade title
func (app ShentuApp) setTmpUpgradeHandler() {
	app.UpgradeKeeper.SetUpgradeHandler(
		tmp,
		func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			ctx.Logger().Info("Initializing Crisis module...")
			constantFee := crisistypes.DefaultGenesisState()
			constantFee.ConstantFee.Denom = app.StakingKeeper.BondDenom(ctx)
			app.CrisisKeeper.InitGenesis(ctx, constantFee)

			fromVM[icatypes.ModuleName] = icaModule.ConsensusVersion()
			// create ICS27 Controller submodule params
			controllerParams := icacontrollertypes.Params{}
			// create ICS27 Host submodule params
			hostParams := icahosttypes.Params{
				HostEnabled: true,
				AllowMessages: []string{
					sdk.MsgTypeURL(&authz.MsgExec{}),
					sdk.MsgTypeURL(&authz.MsgGrant{}),
					sdk.MsgTypeURL(&authz.MsgRevoke{}),
					sdk.MsgTypeURL(&banktypes.MsgSend{}),
					sdk.MsgTypeURL(&distrtypes.MsgWithdrawDelegatorReward{}),
					sdk.MsgTypeURL(&distrtypes.MsgSetWithdrawAddress{}),
					sdk.MsgTypeURL(&distrtypes.MsgWithdrawValidatorCommission{}),
					sdk.MsgTypeURL(&distrtypes.MsgFundCommunityPool{}),
					sdk.MsgTypeURL(&govtypes.MsgVote{}),
					sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}),
					sdk.MsgTypeURL(&stakingtypes.MsgBeginRedelegate{}),
					sdk.MsgTypeURL(&stakingtypes.MsgCreateValidator{}),
					sdk.MsgTypeURL(&stakingtypes.MsgEditValidator{}),
					sdk.MsgTypeURL(&shieldtypes.MsgCreatePool{}),
					sdk.MsgTypeURL(&shieldtypes.MsgDepositCollateral{}),
					sdk.MsgTypeURL(&shieldtypes.MsgWithdrawCollateral{}),
					sdk.MsgTypeURL(&shieldtypes.MsgWithdrawRewards{}),
					sdk.MsgTypeURL(&shieldtypes.MsgWithdrawReimbursement{}),
					sdk.MsgTypeURL(&shieldtypes.MsgPurchaseShield{}),
				},
			}

			ctx.Logger().Info("start to init interchainaccount module...")
			// initialize ICS27 module
			icaModule.InitModule(ctx, controllerParams, hostParams)
			ctx.Logger().Info("Start to run module migrations...")
			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if upgradeInfo.Name == tmp && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}
