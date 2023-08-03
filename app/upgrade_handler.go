package app

import (
	"errors"
	"fmt"
	"time"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	sdkauthtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/shentufoundation/shentu/v2/common"
	authtypes "github.com/shentufoundation/shentu/v2/x/auth/types"
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
func transAddrPrefix(ctx sdk.Context, app ShentuApp) (err error) {
	//if err = transAddrPrefixForStaking(ctx, app); err != nil {
	//	return err
	//}
	//if err = transAddrPrefixForFeegrant(ctx, app); err != nil {
	//	return err
	//}
	//if err = transAddrPrefixForGov(ctx, app); err != nil {
	//	return err
	//}
	//if err = runSlashingMigration(ctx, app); err != nil {
	//	return err
	//}
	//if err = runAuthMigration(ctx, app); err != nil {
	//	return err
	//}
	err = runAuthzMigration(ctx, app)
	return err
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
	skKeeper.IterateAllDelegations(ctx, func(delg stakingtypes.Delegation) bool {
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
	skKeeper.IterateRedelegations(ctx, func(idx int64, red stakingtypes.Redelegation) bool {
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
			for i := range dvvts {
				dvvts[i].DelegatorAddress, err = common.PrefixToShentu(dvvts[i].DelegatorAddress)
				if err != nil {
					return true
				}
				dvvts[i].ValidatorDstAddress, err = common.PrefixToShentu(dvvts[i].ValidatorDstAddress)
				if err != nil {
					return true
				}
				dvvts[i].ValidatorSrcAddress, err = common.PrefixToShentu(dvvts[i].ValidatorSrcAddress)
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
	skKeeper.IterateUnbondingDelegations(ctx, func(idx int64, ubd stakingtypes.UnbondingDelegation) bool {
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
			for i := range dvps {
				dvps[i].DelegatorAddress, err = common.PrefixToShentu(dvps[i].DelegatorAddress)
				if err != nil {
					return true
				}
				dvps[i].ValidatorAddress, err = common.PrefixToShentu(dvps[i].ValidatorAddress)
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
	skKeeper.IterateHistoricalInfo(ctx, func(hi stakingtypes.HistoricalInfo) bool {
		for i := range hi.Valset {
			hi.Valset[i].OperatorAddress, err = common.PrefixToShentu(hi.Valset[i].OperatorAddress)
			if err != nil {
				return true
			}
		}
		skKeeper.SetHistoricalInfo(ctx, hi.Header.Height, &hi)
		return false
	})
	return err
}

func transAddrPrefixForFeegrant(ctx sdk.Context, app ShentuApp) (err error) {
	fgKeeper := app.FeegrantKeeper
	fgKeeper.IterateAllFeeAllowances(ctx, func(grant feegrant.Grant) bool {
		grant.Grantee, err = common.PrefixToShentu(grant.Grantee)
		if err != nil {
			return true
		}
		grant.Granter, err = common.PrefixToShentu(grant.Granter)
		if err != nil {
			return true
		}
		var granteeAcc, granterAcc sdk.AccAddress
		var allowance feegrant.FeeAllowanceI
		granteeAcc, err = sdk.AccAddressFromBech32(grant.Grantee)
		if err != nil {
			return true
		}
		granterAcc, err = sdk.AccAddressFromBech32(grant.Granter)
		if err != nil {
			return true
		}
		allowance, err = grant.GetGrant()
		err = fgKeeper.GrantAllowance(ctx, granterAcc, granteeAcc, allowance)
		return err != nil
	})
	return err
}

func transAddrPrefixForGov(ctx sdk.Context, app ShentuApp) (err error) {
	govKeeper := app.GovKeeper.Keeper
	govKeeper.IterateAllDeposits(ctx, func(deposit govtypes.Deposit) (stop bool) {
		deposit.Depositor, err = common.PrefixToShentu(deposit.Depositor)
		if err != nil {
			return true
		}
		govKeeper.SetDeposit(ctx, deposit)
		return false
	})
	if err != nil {
		return err
	}
	govKeeper.IterateAllVotes(ctx, func(vote govtypes.Vote) (stop bool) {
		vote.Voter, err = common.PrefixToShentu(vote.Voter)
		if err != nil {
			return true
		}
		govKeeper.SetVote(ctx, vote)
		return false
	})
	return err
}

func runSlashingMigration(ctx sdk.Context, app ShentuApp) (err error) {
	sk := app.SlashingKeeper
	sk.IterateValidatorSigningInfos(ctx, func(address sdk.ConsAddress, info slashingtypes.ValidatorSigningInfo) (stop bool) {
		info.Address, err = common.PrefixToShentu(info.Address)
		if err != nil {
			return true
		}
		sk.SetValidatorSigningInfo(ctx, address, info)

		return false
	})

	return err
}

func runAuthMigration(ctx sdk.Context, app ShentuApp) (err error) {
	ak := app.AccountKeeper
	ak.IterateAccounts(ctx, func(acc sdkauthtypes.AccountI) (stop bool) {
		var newAddr string
		newAddr, err = common.PrefixToShentu(acc.GetAddress().String())
		if err != nil {
			return true
		}

		switch account := acc.(type) {
		case *sdkauthtypes.BaseAccount:
			account.Address = newAddr
			ak.SetAccount(ctx, account)
		case *sdkauthtypes.ModuleAccount:
			account.Address = newAddr
			ak.SetAccount(ctx, account)
		case *authtypes.ManualVestingAccount:
			var newUnlocker string
			newUnlocker, err = common.PrefixToShentu(account.Unlocker)
			if err != nil {
				return true
			}
			account.Address = newAddr
			account.Unlocker = newUnlocker
			ak.SetAccount(ctx, account)
		default:
			err = errors.New("unknown account type")
			return true
		}
		return false
	})
	return err
}

func runAuthzMigration(ctx sdk.Context, app ShentuApp) (err error) {
	ak := app.AuthzKeeper
	ak.IterateGrants(ctx, func(granterAddr sdk.AccAddress, granteeAddr sdk.AccAddress, grant authz.Grant) bool {
		authorization := grant.Authorization
		value := authorization.GetValue()

		switch authorization.GetTypeUrl() {
		case "/cosmos.authz.v1beta1.GenericAuthorization":
		case "/cosmos.staking.v1beta1.StakeAuthorization":
			stakeAuthorization := &stakingtypes.StakeAuthorization{}
			if err = stakeAuthorization.Unmarshal(value); err != nil {
				return true
			}
			if err = processStakeAuthorization(stakeAuthorization); err != nil {
				return true
			}
			if err := ak.SaveGrant(ctx, granterAddr, granteeAddr, stakeAuthorization, grant.Expiration); err != nil {
				return true
			}
		default:
			err = errors.New("unknown authorization types")
			return true
		}
		return false
	})
	return err
}

func processStakeAuthorization(stakeAuthorization *stakingtypes.StakeAuthorization) error {
	denyList := stakeAuthorization.GetDenyList()
	allowList := stakeAuthorization.GetAllowList()
	if denyList.Size() > 0 {
		newList, err := prefixToShentuAddrs(denyList.GetAddress())
		if err != nil {
			return err
		}
		stakeAuthorization.Validators = &stakingtypes.StakeAuthorization_DenyList{DenyList: &stakingtypes.StakeAuthorization_Validators{Address: newList}}
	}
	if allowList.Size() > 0 {
		newList, err := prefixToShentuAddrs(denyList.GetAddress())
		if err != nil {
			return err
		}
		stakeAuthorization.Validators = &stakingtypes.StakeAuthorization_AllowList{AllowList: &stakingtypes.StakeAuthorization_Validators{Address: newList}}
	}
	return nil
}

func prefixToShentuAddrs(addrs []string) (newAddrs []string, err error) {
	for _, addr := range addrs {
		var newAddr string
		newAddr, err = common.PrefixToShentu(addr)
		if err != nil {
			return newAddrs, err
		}
		newAddrs = append(newAddrs, newAddr)
	}
	return newAddrs, err
}
