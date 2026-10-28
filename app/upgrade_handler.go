package app

import (
	"context"
	"fmt"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	blacklisttypes "github.com/shentufoundation/shentu/v2/x/blacklist/types"
)

const upgradeName = "v2.20.0"

// v220BootstrapBlacklist is seeded into the blacklist store during the
// v2.20.0 upgrade, so the listed addresses are blocked the moment the upgrade
// goes live. Additional bans and releases afterwards are managed via
// governance proposals (MsgBanAddress / MsgUnbanAddress).
var v220BootstrapBlacklist = []blacklisttypes.Entry{
	{Address: "shentu1nypcgyxxjdrey6h32pxm84uuvsjulszevzl4yx", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu14ykck33qhfku7t9cny3ns6md966anxy22g2y6n", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1pkquacg05j5gcqxwlguyc0qhjhulr0y059ks3e", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1t25cls37d54xh3us8pfdkvttp3mcce2yn5rnpl", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu17e6xj8nt3hetlp7ef88yyhvk24d80stgrjnsnd", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu13gcp0t266cewgp9c0008csp6ukex9p4fy44m26", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu16x8n3deh5sl0hp9pegymmeyg0xz685z0mkgayp", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1487ra2xdmsj97yz68hzvxr0kqyxst90wgx59r4", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1l3zvcvlvfl6wvdqu6tpzdexk9e8vcv5q3xslgw", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu18gpxe2wpy7vuqgdwdtfpwzr45vvud5le5nge8a", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1euq50cwl3nnjndf46yvetgspfvnx7nvukumqp5", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu12r95v4dy8phmadme2uqxh77pdrfuh54jzaht8m", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1x4zl9n2j5495rpfxjcqhmzfnmmyq23pk3anhpx", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1tajz5g5u5zzazj5m4ysgl6cuf9sylre2v6qhp8", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1yqh5cjr4zvxg7dv0hwq64dw5xskswq78pft203", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu15aq7npcvmvdh3ce534x4vv3jn95pvcrtmwc9qc", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1cqtt5zq4s3gc7w6zzvtux8mgzd676nr093l8kg", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1lw6amz84hkf0qns5pqzj3z5j08vc24m7uy98rp", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1vxvf8xfn7v76c7s6yx8s8agsasfygk6n287jdr", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu164hw0w02qpddpwf2hdv4puhgw98yz0h7x26mw5", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu138zwxxv9q6zdtwu5w5ct2lg8nhk6pmxpjwwm2t", Category: blacklisttypes.CATEGORY_COMPROMISED},
}

func (app ShentuApp) setUpgradeHandler(_ codec.BinaryCodec) {
	app.UpgradeKeeper.SetUpgradeHandler(
		upgradeName,
		func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			fromVM, err := app.mm.RunMigrations(ctx, app.configurator, fromVM)
			if err != nil {
				return nil, err
			}

			// seed the bootstrap blacklist: the entries take effect
			// immediately at the upgrade height.
			sdkCtx := sdk.UnwrapSDKContext(ctx)
			if err := app.SeedV220Blacklist(sdkCtx); err != nil {
				return nil, err
			}

			return fromVM, nil
		},
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if upgradeInfo.Name == upgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{
				blacklisttypes.StoreKey,
			},
		}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}

// SeedV220Blacklist writes the v2.20.0 bootstrap blacklist entries.
func (app ShentuApp) SeedV220Blacklist(ctx sdk.Context) error {
	for i, entry := range v220BootstrapBlacklist {
		if err := app.BlacklistKeeper.SetEntry(ctx, entry); err != nil {
			return fmt.Errorf("failed to seed blacklist entry %d (%s): %w", i, entry.Address, err)
		}
	}
	return nil
}
