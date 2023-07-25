package app

import (
	"fmt"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
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

			runSlashingMigration(app, ctx)
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

func runSlashingMigration(app ShentuApp, ctx sdk.Context) {
	k := app.SlashingKeeper
	store := ctx.KVStore(app.keys[slashingtypes.StoreKey])

	// get data
	signingInfos := make([]slashingtypes.SigningInfo, 0)
	missedBlocks := make([]slashingtypes.ValidatorMissedBlocks, 0)
	k.IterateValidatorSigningInfos(ctx, func(address sdk.ConsAddress, info slashingtypes.ValidatorSigningInfo) (stop bool) {
		bechAddr := address.String()
		signingInfos = append(signingInfos, slashingtypes.SigningInfo{
			Address:              bechAddr,
			ValidatorSigningInfo: info,
		})

		localMissedBlocks := k.GetValidatorMissedBlocks(ctx, address)

		missedBlocks = append(missedBlocks, slashingtypes.ValidatorMissedBlocks{
			Address:      bechAddr,
			MissedBlocks: localMissedBlocks,
		})

		return false
	})

	// migration data
	for _, info := range signingInfos {
		oldAddr, err := sdk.ConsAddressFromBech32(info.Address)
		if err != nil {
			panic(err)
		}
		newAddr, err := common.PrefixToShentu(info.Address)
		if err != nil {
			panic(err)
		}
		newConsAddress, err := sdk.ConsAddressFromBech32(newAddr)
		if err != nil {
			panic(err)
		}
		info.ValidatorSigningInfo.Address, err = common.PrefixToShentu(info.ValidatorSigningInfo.Address)
		if err != nil {
			panic(err)
		}

		k.SetValidatorSigningInfo(ctx, newConsAddress, info.ValidatorSigningInfo)
		store.Delete(slashingtypes.ValidatorSigningInfoKey(oldAddr))
	}

	for _, array := range missedBlocks {
		oldAddr, err := sdk.ConsAddressFromBech32(array.Address)
		if err != nil {
			panic(err)
		}
		newAddr, err := common.PrefixToShentu(array.Address)
		if err != nil {
			panic(err)
		}
		newConsAddress, err := sdk.ConsAddressFromBech32(newAddr)
		if err != nil {
			panic(err)
		}

		for _, missed := range array.MissedBlocks {
			k.SetValidatorMissedBlockBitArray(ctx, newConsAddress, missed.Index, missed.Missed)
			store.Delete(slashingtypes.ValidatorMissedBlockBitArrayKey(oldAddr, missed.Index))
		}
	}
}
