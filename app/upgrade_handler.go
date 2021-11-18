package app

import (
	"fmt"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	sdkauthtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	"github.com/cosmos/cosmos-sdk/x/authz"
	sdkauthz "github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	sdkfeegrant "github.com/cosmos/cosmos-sdk/x/feegrant"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/modules/core/03-connection/types"

	authtypes "github.com/certikfoundation/shentu/v2/x/auth/types"
)

const upgradeName = "Shentu-v230"

func (app ShentuApp) setUpgradeHandler() {
	app.upgradeKeeper.SetUpgradeHandler(
		upgradeName,
		func(ctx sdk.Context, _ upgradetypes.Plan, _ module.VersionMap) (module.VersionMap, error) {
			app.ibcKeeper.ConnectionKeeper.SetParams(ctx, ibcconnectiontypes.DefaultParams())

			fromVM := make(map[string]uint64)
			for moduleName := range app.mm.Modules {
				fromVM[moduleName] = 1
			}
			// override versions for _new_ modules as to not skip InitGenesis
			fromVM[sdkauthz.ModuleName] = 0
			fromVM[sdkfeegrant.ModuleName] = 0
			fromVM[sdkauthtypes.ModuleName] = 2

			temp, err := app.mm.RunMigrations(ctx, app.configurator, fromVM)
			if err != nil {
				return temp, err
			}
			var iterErr error
			app.accountKeeper.IterateAccounts(ctx, func(account sdkauthtypes.AccountI) (stop bool) {
				wb, err := migrateVestingAccount(app, ctx, account)
				if err != nil {
					iterErr = err
					return true
				}

				if wb == nil {
					return false
				}
				app.accountKeeper.SetAccount(ctx, wb)
				return false
			})

			return temp, iterErr
		},
	)

	upgradeInfo, err := app.upgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if upgradeInfo.Name == upgradeName && !app.upgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{authz.ModuleName, feegrant.ModuleName},
		}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}

func migrateVestingAccount(app ShentuApp, ctx sdk.Context, account sdkauthtypes.AccountI) (sdkauthtypes.AccountI, error) {
	vacc, ok := account.(exported.VestingAccount)
	if !ok {
		return nil, nil
	}

	mvacc, ok := vacc.(*authtypes.ManualVestingAccount)
	if !ok {
		return nil, fmt.Errorf("unexpected type: %T. wanted *authtypes.ManualVestingAccount", account)
	}

	bondDenom := app.stakingKeeper.BondDenom(ctx)

	addr := mvacc.GetAddress()
	balances := app.bankKeeper.GetAllBalances(ctx, addr)
	delegations := app.stakingKeeper.GetAllDelegatorDelegations(ctx, addr)
	unbondings := app.stakingKeeper.GetAllUnbondingDelegations(ctx, addr)

	delegationsSum := sdk.NewCoins()
	for _, d := range delegations {
		delResp, err := stakingkeeper.DelegationToDelegationResponse(ctx, app.stakingKeeper.Keeper, d)
		if err != nil {
			panic(err)
		}

		delegationsSum = delegationsSum.Add(delResp.GetBalance())
	}

	unbondingsSum := sdk.NewCoins()
	for _, u := range unbondings {
		for _, e := range u.Entries {
			unbondingsSum = unbondingsSum.Add(sdk.NewCoin(bondDenom, e.Balance))
		}
	}

	delegationsSum = delegationsSum.Add(unbondingsSum...)

	mvacc.DelegatedFree = sdk.NewCoins()
	mvacc.DelegatedVesting = sdk.NewCoins()

	for _, coin := range delegationsSum {
		balances = balances.Add(coin)
	}

	mvacc.TrackDelegation(ctx.BlockTime(), balances, delegationsSum)
	return mvacc, nil
}
