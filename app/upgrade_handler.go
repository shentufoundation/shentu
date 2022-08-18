package app

import (
	"fmt"
	shieldtypes "github.com/certikfoundation/shentu/v2/x/shield/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
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
			ctx.Logger().Info("Start to run module migrations...")
			ctx.Logger().Info("Start to run module migrations...")
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
	var rewards shieldtypes.MixedDecCoins
	for _, provider := range sk.GetAllProviders(ctx) {
		rewards = rewards.Add(provider.Rewards)
	}

	totalInt, _ := remainingServiceFees.Add(rewards).Native.TruncateDecimal()

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
	blockServiceFees := sk.GetBlockServiceFees(ctx).Native.AmountOf(bondDenom).TruncateInt()

	totalInt = totalInt.Add(sdk.NewCoin(bondDenom, shieldStake)).Add(sdk.NewCoin(bondDenom, reimbursement)).Add(sdk.NewCoin(bondDenom, blockServiceFees))

	if moduleCoins.IsAllGTE(totalInt) {
		blockServiceFees = blockServiceFees.Add(moduleCoins.Sub(totalInt).AmountOf(bondDenom))
		newBlockServiceFees := shieldtypes.NewMixedDecCoins(sdk.NewDecCoinsFromCoins(sdk.NewCoin(bondDenom, blockServiceFees)), sdk.NewDecCoins())
		sk.SetBlockServiceFees(ctx, newBlockServiceFees)
	} else {
		diff := totalInt.Sub(moduleCoins)
		fmt.Println("diff: ", diff)
		// first try to take away from remaining services
		rSFInt, decimals := remainingServiceFees.Native.TruncateDecimal()
		if !rSFInt.IsAllGTE(diff) {
			additionalFunds := diff.Sub(rSFInt)
			app.BankKeeper.SendCoinsFromModuleToModule(ctx, distrtypes.ModuleName, shieldtypes.ModuleName, additionalFunds)
			fp := app.DistrKeeper.GetFeePool(ctx)
			fp.CommunityPool = fp.CommunityPool.Sub(sdk.NewDecCoinsFromCoins(additionalFunds...))
			app.DistrKeeper.SetFeePool(ctx, fp)
			moduleCoins = app.BankKeeper.GetAllBalances(ctx, app.AccountKeeper.GetModuleAccount(ctx, shieldtypes.ModuleName).GetAddress())
		}
		rSFInt = rSFInt.Sub(diff)
		remainingServiceFees.Native = sdk.NewDecCoinsFromCoins(rSFInt...)
		remainingServiceFees.Native = remainingServiceFees.Native.Add(decimals...)
		app.ShieldKeeper.SetRemainingServiceFees(ctx, remainingServiceFees)
		fmt.Printf("remainingServiceFees: %s\n", remainingServiceFees.Native.String())
		fmt.Printf("rewards: %s\n", rewards.Native.String())
		fmt.Printf("reimbursement: %s\n", reimbursement.String())
		fmt.Printf("blockServiceFees: %s\n", blockServiceFees.String())
		fmt.Printf("shieldStake: %s\n", shieldStake.String())
		fmt.Printf("TotalInt: %s\n", totalInt.String())
		fmt.Printf("ModuleCoins: %s\n", moduleCoins.String())
	}
}
