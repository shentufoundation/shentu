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

	shieldtypes "github.com/shentufoundation/shentu/v2/x/shield/types"
)

const (
	upgradeName = "v2.6.0"
)

func (app ShentuApp) setUpgradeHandler() {
	app.UpgradeKeeper.SetUpgradeHandler(
		upgradeName,
		func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			// don't run Initgenesis since it'll be set with a wrong denom
			fromVM[crisistypes.ModuleName] = app.mm.Modules[crisistypes.ModuleName].ConsensusVersion()
			// don't run icamodule's Initgenesis since it'll overwrite the icahost params that be set here
			// the InitModule will be called later on to set params.
			// this assumes it's the first time ica module go into fromVM
			fromVM[icatypes.ModuleName] = app.mm.Modules[icatypes.ModuleName].ConsensusVersion()

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

			crisisGenesis := crisistypes.DefaultGenesisState()
			crisisGenesis.ConstantFee.Denom = app.StakingKeeper.BondDenom(ctx)
			app.CrisisKeeper.InitGenesis(ctx, crisisGenesis)

			ctx.Logger().Info("Start to run module migrations...")
			newVersionMap, err := app.mm.RunMigrations(ctx, app.configurator, fromVM)

			ctx.Logger().Info("Fixing Shield invariant...")
			RunShieldMigration(app, ctx)

			return newVersionMap, err
		},
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if upgradeInfo.Name == upgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{crisistypes.ModuleName, icahosttypes.SubModuleName},
		}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}

func RunShieldMigration(app ShentuApp, ctx sdk.Context) {
	sk := app.ShieldKeeper
	bondDenom := sk.BondDenom(ctx)

	// remaining service fees
	remainingServiceFees := sk.GetRemainingServiceFees(ctx)

	// rewards
	var rewards sdk.DecCoins
	for _, provider := range sk.GetAllProviders(ctx) {
		rewards = rewards.Add(provider.Rewards...)
	}

	totalInt, remainder := remainingServiceFees.Add(rewards...).TruncateDecimal()
	if !remainder.Empty() {
		panic("remaining coins in the shield module is not an sdk.Int")
	}

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

	// sum of total coins tracked by the shield module
	totalInt = totalInt.Add(sdk.NewCoin(bondDenom, shieldStake)).Add(sdk.NewCoin(bondDenom, reimbursement)).Add(blockServiceFees...)
	// actual balance of shield module
	moduleCoins := app.BankKeeper.GetAllBalances(ctx, app.AccountKeeper.GetModuleAccount(ctx, shieldtypes.ModuleName).GetAddress())

	if moduleCoins.IsAllGTE(totalInt) {
		// if actual balance is greater, send remainder to the next block reward
		blockServiceFees = blockServiceFees.Add(moduleCoins.Sub(totalInt)...)
		newBlockServiceFees := sdk.NewDecCoinsFromCoins(blockServiceFees...)
		sk.SetBlockServiceFees(ctx, newBlockServiceFees)
	} else {
		diff := totalInt.Sub(moduleCoins) // assuming there is only CTK in shield module.
		ctx.Logger().Info("Shield Module Account Coin diff: ", diff)
		// first try to take away from remaining service fees
		rSFInt, decimals := remainingServiceFees.TruncateDecimal()
		if !rSFInt.IsAllGTE(diff) {
			// if the remaining service fees is not enough, take the diff from the community pool.
			additionalFunds := diff.Sub(rSFInt)
			app.BankKeeper.SendCoinsFromModuleToModule(ctx, distrtypes.ModuleName, shieldtypes.ModuleName, additionalFunds)
			fp := app.DistrKeeper.GetFeePool(ctx)
			fp.CommunityPool = fp.CommunityPool.Sub(sdk.NewDecCoinsFromCoins(additionalFunds...))
			app.DistrKeeper.SetFeePool(ctx, fp)
			moduleCoins = app.BankKeeper.GetAllBalances(ctx, app.AccountKeeper.GetModuleAccount(ctx, shieldtypes.ModuleName).GetAddress())
			diff = totalInt.Sub(moduleCoins)
		}
		rSFInt = rSFInt.Sub(diff)
		if rSFInt.IsAnyNegative() {
			panic(fmt.Sprintf("remaining service fees - module coins diff < 0.\nRSF: %s, diff: %s", rSFInt.Add(diff...), diff))
		}
		remainingServiceFees = sdk.NewDecCoinsFromCoins(rSFInt...)
		remainingServiceFees = remainingServiceFees.Add(decimals...)
		app.ShieldKeeper.SetRemainingServiceFees(ctx, remainingServiceFees)
		ctx.Logger().Info(fmt.Sprintf("remainingServiceFees: %s", remainingServiceFees.String()))
		ctx.Logger().Info(fmt.Sprintf("rewards: %s", rewards.String()))
		ctx.Logger().Info(fmt.Sprintf("reimbursement: %s", reimbursement.String()))
		ctx.Logger().Info(fmt.Sprintf("blockServiceFees: %s", blockServiceFees.String()))
		ctx.Logger().Info(fmt.Sprintf("shieldStake: %s", shieldStake.String()))
		ctx.Logger().Info(fmt.Sprintf("TotalInt: %s", totalInt.String()))
		ctx.Logger().Info(fmt.Sprintf("ModuleCoins: %s", moduleCoins.String()))
	}
}
