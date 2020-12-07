package main

import (
	"encoding/json"
	"io"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/snapshots"
	"github.com/cosmos/cosmos-sdk/store"

	"github.com/cosmos/cosmos-sdk/client/flags"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"

	"github.com/certikfoundation/shentu/app"
	certikinit "github.com/certikfoundation/shentu/cmd/init"
	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/x/staking"
)

func main() {
	cobra.EnableCommandSorting = false

	cdc := app.MakeCodec()

	// read in the configuration file for the sdk
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(common.Bech32PrefixAccAddr, common.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(common.Bech32PrefixValAddr, common.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(common.Bech32PrefixConsAddr, common.Bech32PrefixConsPub)
	config.Seal()

	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "certikd",
		Short:             "CertiK App Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
		Run: func(cmd *cobra.Command, args []string) {
			docDir, err := cmd.Flags().GetString(certikinit.DocFlag)
			if err == nil && docDir != "" {
				certikinit.GenDoc(cmd, docDir)
			} else if err = cmd.Help(); err != nil {
				panic(err)
			}
		},
	}

	rootCmd.Flags().StringP(certikinit.DocFlag, certikinit.DocFlagAbbr, "", certikinit.DocFlagUsage)

	rootCmd.AddCommand(genutilcli.InitCmd(ctx, cdc, app.ModuleBasics, app.DefaultNodeHome))
	rootCmd.AddCommand(genutilcli.CollectGenTxsCmd(ctx, cdc, auth.GenesisAccountIterator{}, app.DefaultNodeHome))
	rootCmd.AddCommand(AddGenesisAccountCmd(app.DefaultNodeHome))
	rootCmd.AddCommand(AddGenesisCertifierCmd(ctx, cdc))
	rootCmd.AddCommand(AddGenesisShieldAdminCmd(ctx, cdc))
	rootCmd.AddCommand(MigrateGenesisCmd(ctx, cdc))
	rootCmd.AddCommand(certikinit.TestnetFilesCmd(ctx, cdc, app.ModuleBasics, auth.GenesisAccountIterator{}))
	rootCmd.AddCommand(genutilcli.GenTxCmd(
		ctx,
		cdc,
		app.ModuleBasics,
		staking.AppModuleBasic{},
		auth.GenesisAccountIterator{},
		app.DefaultNodeHome,
		app.DefaultCLIHome,
	))
	rootCmd.AddCommand(genutilcli.ValidateGenesisCmd(ctx, cdc, app.ModuleBasics))
	rootCmd.AddCommand(version.Cmd)
	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndValidators)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "NS", app.DefaultNodeHome)
	err := executor.Execute()
	if err != nil {
		// handle with #870
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer, appOpts servertypes.AppOptions) servertypes.Application {
	var cache sdk.MultiStorePersistentCache

	if cast.ToBool(appOpts.Get(server.FlagInterBlockCache)) {
		cache = store.NewCommitKVStoreCacheManager()
	}

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	pruningOpts, err := server.GetPruningOptionsFromFlags(appOpts)
	if err != nil {
		panic(err)
	}

	snapshotDir := filepath.Join(cast.ToString(appOpts.Get(flags.FlagHome)), "data", "snapshots")
	snapshotDB, err := sdk.NewLevelDB("metadata", snapshotDir)
	if err != nil {
		panic(err)
	}
	snapshotStore, err := snapshots.NewStore(snapshotDB, snapshotDir)
	if err != nil {
		panic(err)
	}

	return app.NewCertiKApp(logger, db, traceStore, true, map[int64]bool{}, cast.ToString(appOpts.Get(flags.FlagHome)), uint(10000),
		app.MakeEncodingConfig(), // Ideally, we would reuse the one created by NewRootCmd.
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
		baseapp.SetHaltHeight(uint64(viper.GetInt(server.FlagHaltHeight))),
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(cast.ToString(appOpts.Get(server.FlagMinGasPrices))),
		baseapp.SetMinRetainBlocks(cast.ToUint64(appOpts.Get(server.FlagMinRetainBlocks))),
		baseapp.SetHaltHeight(cast.ToUint64(appOpts.Get(server.FlagHaltHeight))),
		baseapp.SetHaltTime(cast.ToUint64(appOpts.Get(server.FlagHaltTime))),
		baseapp.SetInterBlockCache(cache),
		baseapp.SetTrace(cast.ToBool(appOpts.Get(server.FlagTrace))),
		baseapp.SetIndexEvents(cast.ToStringSlice(appOpts.Get(server.FlagIndexEvents))),
		baseapp.SetSnapshotStore(snapshotStore),
		baseapp.SetSnapshotInterval(cast.ToUint64(appOpts.Get(server.FlagStateSyncSnapshotInterval))),
		baseapp.SetSnapshotKeepRecent(cast.ToUint32(appOpts.Get(server.FlagStateSyncSnapshotKeepRecent))),
	)
}

func exportAppStateAndValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	if height != -1 {
		cApp := app.NewCertiKApp(logger, db, traceStore, false, map[int64]bool{}, uint(1), app.MakeEncodingConfig())
		err := cApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return cApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}
	cApp := app.NewCertiKApp(logger, db, traceStore, true, map[int64]bool{}, uint(1), app.MakeEncodingConfig())
	return cApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}
