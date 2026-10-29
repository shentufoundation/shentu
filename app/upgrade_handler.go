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
	{Address: "shentu1sc4fnz8usty7z4leqgrcfyjx2qyf5vrdw7geu4", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1rmpduv5jva5l4u4uh3c6lhx8r3zax6xl24jrj8", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu14ys54ka0ggenylw9rlcrrc0wjq8qrl48hddypx", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1t52aa68sv4jhy7xtyssqhrzy3u227czgk525cr", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1gkg9rapa0pd05azk6hgyddplamvvfedvs5c8xy", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1yq8a2ksa7dz8wd8wlks3k8nqdmht76xdttshd2", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1xn45plrz5dpeuk997tkkhyu9fv0a658tpzhhwz", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1n6jktj6x9hfv62ck7a6nrnzpq6fwwvgyfdrnld", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1487p6e4fn34fk5v65g6f58rwnp6dgeu9xnns8d", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1ts0t9kszwvqf3d7kc50wq25hpe3haztqvyjftx", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1hl5rksmf5uqkxyje3dkr9enufy0g8swc2ugh6j", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1e5wrmafqyt3ssmtjxhrp8v83lh02x5zmdqjaw8", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu15xc9lkxyvwf252gvsvgauaazf0jer8zyz9lsyw", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1fmc2vs87qyfa94q08rp0wtj7vc6y8hex3nd40h", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1exlzdrpgn7lza4yp29rqrwtqk54ul3atavz36m", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1avc5scdxva26xews6ea26xntjwwpzu0frpjqs0", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu130q9ca8kwqdz42lpuj30vc56lk5z62yunkyvm5", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu17umyc5vr74m0u5806zce09klvc6dum6vqtq9uv", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu17wzc3j092hkydz7pw6znukkc0dm3668rxhzhsg", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1pd46zly352umqa6ucgwzvfmaz3t9cz50l6yssx", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu19yz9xk0xxuu2w4rs3pktq0273v9hqqq8a07urc", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1yl23kgrjjfc8yg3mpemnkg8nvfm847jstnu4tm", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1frs30f8thm0adepn963kjre4kap4uzygq8xdl8", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1wgjq800n0ng822wjuydk2xdty3m9d6cxtsfjjf", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1kcz8jpzkmpfw9a4kgagvf9t0l22vffzaandp0s", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1rypp5zu4ra5xcce042z3z3wyfwsezwd7zyjzg2", Category: blacklisttypes.CATEGORY_COMPROMISED},
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
	{Address: "shentu1k0he7cu3u9s5vesqxhur7qweavfu6rt790vwju", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1cxdw08md6ejpk9gz0pkdz4dq3xjewsadt9ntuz", Category: blacklisttypes.CATEGORY_COMPROMISED},
	{Address: "shentu1uesxh2u9qy0zdmqzf2qw25krth7e5j6nsr230a", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1yfq94ymzevxrq493y0r0jkymajrevux2sg6ul3", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu150cq0ddqem93p6fnv2fxdrkr5sg8xx9gh7xzl9", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1qurjn08pve9v3fpjkvjrmfxj3aka7qay0upwgf", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1kjd3t9xq8frv9frgezw9j4ufrqu83ts2aykdxa", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1xmzqnsfknjawzxxyu64ljyuq6wvgyhd753wh8a", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu17tyhzsyw6zjea8qepa8wkl3njg9jg4h7229y2r", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1v74pf75p3qk9n28qt87yhp4cl42x3knt9kq2dk", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1msvxu0f8mjex26m25z3hqywp2netse6uc8e2h6", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1faazq6pmqe8r8zwum75z4lcjmmfdjql5es988s", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1e6wtkw4yq8rz6qh3dr9wa8wtkukfepzagsw2jz", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1vx6t5872zg7p4nwnnnrrv9p2rw8q49e0j5kntq", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1w0z0k24jvv9vgca45u2numhk34re3vxu3zfmj0", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu15y86ggl94zr6zrfzkm4c6revqcmhc83lgqxdtk", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1r7pq7r90aqjsknxn2mqezq64624na54d0dpy0m", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1xtej4gp0x88hslgw8gpfg22r4e4z66t0d39822", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1mgthz9sde2u5zv7ns0tfhacwtx4eeaxhv8tw9l", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1jkxtzfxpx7e96dma6dlvmny3nfd2epndeuyc0n", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1zaxq0asajadqy8lhgl395jxvzmeraqy6ckq0f7", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu19rtkyzmucvtk3u4cjkzapf42uf3336pf8newxv", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1596pmtswzheuxtpvsnyp5ehu9sn2e7kpfl0z8s", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1606yz0ese8wn0fe0fnhsnv5dr3np73uqamalzc", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu17pckvnex4uqjjhayc9nl4j2kwlqmfx4txusswu", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1fd8hmql6gc6rdw84zgtxdn8m6ly5sulfjamaks", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu165wd4nn3yspzdjhskpgr2fqp3vxrnnh2tkuqcv", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1ae565tkn2gsmn4hw6hedgcdjrtfd62glr6su9c", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu15qc5kam392yn7m6py8pkty8wk7ughw7g62k2qj", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1fl3phjkdg6rd2wkk6acc40hfff03p228ly0c99", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1ngcgl6ksqqp8vwevglxtk48wc2uc0dl40g73h6", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1w78tmrugchfzn2ye5lrcp7t4033069jsw53q37", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu19lvzrdreqmatn83cgr96uj0jwr0avj8976eh3n", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1ek55c2gu3ujrw9luk7gqlw0ny3gqshf7h0mpv4", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1zgqknzzecqwmt0c85pg978ujuekwxax6pycjz0", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1hv4d58gq02z4l7ruuj76yxp72qarh8e0u660cg", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1azxdh422fnx4df78xsk3j4vsz2wl224wygn8am", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1d5dn3uzgk2geazfj8y4xdc0737x2p8ulyxz2pc", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1pcw07rrjad4prqg8pnv2xtvadj54657sz0zyj4", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu127gsr567nqa7le0y7zn8anje0jmsvmktc03dgl", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1gxamt8ar49v675epuhv7se4v89v8dgpmel88jv", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1460zm6z4y27dnm3nt77l02447fhal5e29vtt87", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1dkr0fl6ufwfvqwzyv8dm7dl46qw05plxy6z2ha", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1k23xmfxweh2z8kf37elf7l0wqlrpuy0h89sfd3", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1su0pmg6fnl96m8rsf87gsqgz6rxxkg5unattpz", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu14klwnax76jjmk7ylxsg2vrxfjluwddwrh52ur4", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1m3xlemd7krm06u7ryxxglcrpu7ned4gpkc2dea", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1e7xmy8dk57d6wss276mntghffft4xl49jvk2x8", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1rhf94zrhhapm5wm73yjd9y9jrxj9kcly2zdxgu", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1xdumvqcmr8sur59nrjwql4a8w8vnpp0437k0t4", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1d0ce2t5zfwaa6aw49uyadfsdcv0wvxpkmzcmpw", Category: blacklisttypes.CATEGORY_SUSPECTED},
	{Address: "shentu1d2xl5vuh32qpyxkwzgp6u6u2jd4ps6kfgcxysz", Category: blacklisttypes.CATEGORY_SUSPECTED},
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
