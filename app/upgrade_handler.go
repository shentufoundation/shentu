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
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	ica "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts"
	icacontrollertypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"

	shieldtypes "github.com/shentufoundation/shentu/v2/x/shield/types"
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
			order = append(order, crisistypes.ModuleName, paramtypes.ModuleName)
			app.mm.SetOrderMigrations(order...)

			ctx.Logger().Info("Fixing Shield invariant...")
			RunShieldMigration(app, ctx)

			ctx.Logger().Info("Start to run module migrations...")
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

func RunShieldMigration(app ShentuApp, ctx sdk.Context) {
	sk := app.ShieldKeeper
	bondDenom := sk.BondDenom(ctx)

	moduleCoins := app.BankKeeper.GetAllBalances(ctx, app.AccountKeeper.GetModuleAccount(ctx, shieldtypes.ModuleName).GetAddress())

	// remaining service fees
	remainingServiceFees := sk.GetRemainingServiceFees(ctx)

	// rewards
	var rewards sdk.DecCoins
	for _, provider := range sk.GetAllProviders(ctx) {
		rewards = rewards.Add(provider.Rewards...)
	}

	totalInt, _ := remainingServiceFees.Add(rewards...).TruncateDecimal()

	// shield stake
	shieldStake := sdk.ZeroInt()
	for _, stake := range sk.GetAllStakeForShields(ctx) {
		shieldStake = shieldStake.Add(stake.Amount)
	}

	// reimbursement
	reimbursement := sdk.ZeroInt()
	for _, rmb := range sk.GetAllReimbursements(ctx) {
		reimbursement = reimbursement.Add(rmb.Amount.AmountOf(bondDenom))
	}

	// block service fees
	blockServiceFees, _ := sk.GetBlockServiceFees(ctx).TruncateDecimal()

	totalInt = totalInt.Add(sdk.NewCoin(bondDenom, shieldStake)).Add(sdk.NewCoin(bondDenom, reimbursement)).Add(blockServiceFees...)

	if moduleCoins.IsAllGTE(totalInt) {
		blockServiceFees = blockServiceFees.Add(moduleCoins.Sub(totalInt)...)
		newBlockServiceFees := sdk.NewDecCoinsFromCoins(blockServiceFees...)
		sk.SetBlockServiceFees(ctx, newBlockServiceFees)
	} else {
		diff := totalInt.Sub(moduleCoins)
		fmt.Println("diff: ", diff)
		// first try to take away from remaining services
		rSFInt, decimals := remainingServiceFees.TruncateDecimal()
		if !rSFInt.IsAllGTE(diff) {
			additionalFunds := diff.Sub(rSFInt)
			app.BankKeeper.SendCoinsFromModuleToModule(ctx, distrtypes.ModuleName, shieldtypes.ModuleName, additionalFunds)
			fp := app.DistrKeeper.GetFeePool(ctx)
			fp.CommunityPool = fp.CommunityPool.Sub(sdk.NewDecCoinsFromCoins(additionalFunds...))
			app.DistrKeeper.SetFeePool(ctx, fp)
			moduleCoins = app.BankKeeper.GetAllBalances(ctx, app.AccountKeeper.GetModuleAccount(ctx, shieldtypes.ModuleName).GetAddress())
		}
		rSFInt = rSFInt.Sub(diff)
		remainingServiceFees = sdk.NewDecCoinsFromCoins(rSFInt...)
		remainingServiceFees = remainingServiceFees.Add(decimals...)
		app.ShieldKeeper.SetRemainingServiceFees(ctx, remainingServiceFees)
		fmt.Printf("remainingServiceFees: %s\n", remainingServiceFees.String())
		fmt.Printf("rewards: %s\n", rewards.String())
		fmt.Printf("reimbursement: %s\n", reimbursement.String())
		fmt.Printf("blockServiceFees: %s\n", blockServiceFees.String())
		fmt.Printf("shieldStake: %s\n", shieldStake.String())
		fmt.Printf("TotalInt: %s\n", totalInt.String())
		fmt.Printf("ModuleCoins: %s\n", moduleCoins.String())
	}
}
