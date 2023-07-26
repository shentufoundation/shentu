package app

import (
	"fmt"
	"time"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	stakingTypes  "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/shentufoundation/shentu/v2/common"
)

const (
	upgradeName = "v2.8.0"
)

func (app ShentuApp) setUpgradeHandler() {
	app.UpgradeKeeper.SetUpgradeHandler(
		upgradeName,
		func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			// transfer module consensus version has been bumped to 2
			ctx.Logger().Info("Start to run module migrations...")
			newVersionMap, err := app.mm.RunMigrations(ctx, app.configurator, fromVM)

			if err == nil {
				ctx.Logger().Info("Transite address prefix to shentu for modules - staking, bank ...")
				err = transAddrPrefix(ctx, app)
			}
			return newVersionMap, err
		},
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if upgradeInfo.Name == upgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}

// the function transite bech32 address prefix from 'certik' to 'shentu' for the values stored in related modules
// this function is supposed to be called when chain upgraded from v2.7.1 to v2.8.0
func transAddrPrefix(ctx sdk.Context, app ShentuApp) error {
	if err := transAddrPrefixForStaking(ctx, app); err != nil {
		return err
	}
	return nil
}

func transAddrPrefixForStaking(ctx sdk.Context, app ShentuApp) (err error) {
	skKeeper := app.StakingKeeper.Keeper
	//transite prefix for Validators and UnbondingValidatorQueue
	allValidators := skKeeper.GetAllValidators(ctx)
	for _, v := range allValidators {
		if v.IsUnbonding() {
			skKeeper.DeleteValidatorQueue(ctx, v)
		}
		v.OperatorAddress, err = common.PrefixToShentu(v.OperatorAddress)
		if err != nil {
			return err
		}
		skKeeper.SetValidator(ctx, v)
		if v.IsUnbonding() {
			skKeeper.InsertUnbondingValidatorQueue(ctx, v)
		}
	}
	//transite prefix for delegations
	skKeeper.IterateAllDelegations(ctx, func(delg stakingTypes.Delegation) bool {
		delg.DelegatorAddress, err = common.PrefixToShentu(delg.DelegatorAddress)
		if err != nil {
			return true
		}
		delg.ValidatorAddress, err = common.PrefixToShentu(delg.ValidatorAddress)
		if err != nil {
			return true
		}
		skKeeper.SetDelegation(ctx, delg)
		return false
	})
	if err != nil {
		return err
	}
	//transite prefix for redelegation and redelegationQueue
	touchedTimes := make(map[time.Time]bool)
	skKeeper.IterateRedelegations(ctx, func(idx int64, red stakingTypes.Redelegation) bool {
		red.DelegatorAddress, err = common.PrefixToShentu(red.DelegatorAddress)
		if err != nil {
			return true
		}
		red.ValidatorSrcAddress, err = common.PrefixToShentu(red.ValidatorSrcAddress)
		if err != nil {
			return true
		}
		red.ValidatorDstAddress, err = common.PrefixToShentu(red.ValidatorDstAddress)
		if err != nil {
			return true
		}
		skKeeper.SetRedelegation(ctx, red)
		for _, e := range red.Entries {
			if touchedTimes[e.CompletionTime] {
				continue
			}
			dvvts := skKeeper.GetRedelegationQueueTimeSlice(ctx, e.CompletionTime)
			if len(dvvts) == 0 {
				continue
			}
			for _, t := range dvvts {
				t.DelegatorAddress, err = common.PrefixToShentu(t.DelegatorAddress)
				if err != nil {
					return true
				}
				t.ValidatorDstAddress, err = common.PrefixToShentu(t.ValidatorDstAddress)
				if err != nil {
					return true
				}
				t.ValidatorSrcAddress, err = common.PrefixToShentu(t.ValidatorSrcAddress)
				if err != nil {
					return true
				}
			}
			skKeeper.SetRedelegationQueueTimeSlice(ctx, e.CompletionTime, dvvts)
			touchedTimes[e.CompletionTime] = true
		}
		return false
	})
	if err != nil {
		return err
	}
	//transite prefix for UnbondingDelegation and UnbondingQueue
	skKeeper.IterateUnbondingDelegations(ctx, func(idx int64, ubd stakingTypes.UnbondingDelegation) bool {
		ubd.DelegatorAddress, err = common.PrefixToShentu(ubd.DelegatorAddress)
		if err != nil {
			return true
		}
		ubd.ValidatorAddress, err = common.PrefixToShentu(ubd.ValidatorAddress)
		if err != nil {
			return true
		}
		skKeeper.SetUnbondingDelegation(ctx, ubd)
		touchedTimes := make(map[time.Time]bool)
		for _, e := range ubd.Entries {
			if touchedTimes[e.CompletionTime] {
				continue
			}
			dvps := skKeeper.GetUBDQueueTimeSlice(ctx, e.CompletionTime)
			for _, p := range dvps {
				p.DelegatorAddress, err = common.PrefixToShentu(p.DelegatorAddress)
				if err != nil {
					return true
				}
				p.ValidatorAddress, err = common.PrefixToShentu(p.ValidatorAddress)
				if err != nil {
					return true
				}
			}
			skKeeper.SetUBDQueueTimeSlice(ctx, e.CompletionTime, dvps)
			touchedTimes[e.CompletionTime] = true
		}
		return false
	})
	if err != nil {
		return err
	}
	//transite prefix for HistoricalInfo
	skKeeper.IterateHistoricalInfo(ctx, func(hi stakingTypes.HistoricalInfo) bool {
		for _, v := range hi.Valset {
			v.OperatorAddress, err = common.PrefixToShentu(v.OperatorAddress)
			if err != nil {
				return true
			}
		}
		skKeeper.SetHistoricalInfo(ctx, hi.Header.Height, &hi)
		return false
	})

	return err
}
