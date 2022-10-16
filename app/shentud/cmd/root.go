package cmd

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	tmcmds "github.com/tendermint/tendermint/cmd/tendermint/commands"
	tmcfg "github.com/tendermint/tendermint/config"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/snapshots"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkauthcli "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingcli "github.com/cosmos/cosmos-sdk/x/auth/vesting/client/cli"
	sdkbankcli "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"

	"github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/app/params"
	shentuinit "github.com/shentufoundation/shentu/v2/app/shentud/init"
	"github.com/shentufoundation/shentu/v2/common"
	authcli "github.com/shentufoundation/shentu/v2/x/auth/client/cli"
	bankcli "github.com/shentufoundation/shentu/v2/x/bank/client/cli"
	"github.com/shentufoundation/shentu/v2/x/crisis"
	cvmcli "github.com/shentufoundation/shentu/v2/x/cvm/client/cli"
)

// NewRootCmd creates a new root command for simd. It is called once in the
// main function.
func NewRootCmd() (*cobra.Command, params.EncodingConfig) {
	encodingConfig := app.MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(app.DefaultNodeHome).
		WithViper("SHENTU")

	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(common.Bech32PrefixAccAddr, common.Bech32PrefixAccPub)
	cfg.SetBech32PrefixForValidator(common.Bech32PrefixValAddr, common.Bech32PrefixValPub)
	cfg.SetBech32PrefixForConsensusNode(common.Bech32PrefixConsAddr, common.Bech32PrefixConsPub)
	cfg.Seal()

	rootCmd := &cobra.Command{
		Use:   "shentud",
		Short: "Stargate Shentu Chain App",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}
			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			template, customAppConfig := initAppConfig()
			return server.InterceptConfigsPreRunHandler(cmd, template, customAppConfig)
		},
		Run: func(cmd *cobra.Command, args []string) {
			docDir, err := cmd.Flags().GetString(shentuinit.DocFlag)
			if err == nil && docDir != "" {
				shentuinit.GenDoc(cmd, docDir)
			} else if err = cmd.Help(); err != nil {
				panic(err)
			}
		},
	}

	initRootCmd(rootCmd, encodingConfig)

	return rootCmd, encodingConfig
}

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, interface{}) {
	type CustomAppConfig struct {
		serverconfig.Config
	}

	// Optionally allow the chain developer to overwrite the SDK's default
	// server config.
	srvCfg := serverconfig.DefaultConfig()
	srvCfg.API.Enable = true
	srvCfg.StateSync.SnapshotInterval = 1500
	srvCfg.StateSync.SnapshotKeepRecent = 2
	srvCfg.MinGasPrices = "0.025" + common.MicroCTKDenom
	srvCfg.IAVLDisableFastNode = false

	AppTemplate := serverconfig.DefaultConfigTemplate

	return AppTemplate, CustomAppConfig{*srvCfg}
}

// Execute executes the root command of an application. It handles creating a
// server context object with the appropriate server and client objects injected
// into the underlying stdlib Context. It also handles adding core CLI flags,
// specifically the logging flags. It returns an error upon execution failure.
func Execute(rootCmd *cobra.Command, defaultHome string) error {
	// Create and set a client.Context on the command's Context. During the pre-run
	// of the root command, a default initialized client.Context is provided to
	// seed child command execution with values such as AccountRetriver, Keyring,
	// and a Tendermint RPC. This requires the use of a pointer reference when
	// getting and setting the client.Context. Ideally, we utilize
	// https://github.com/spf13/cobra/pull/1118.
	srvCtx := server.NewDefaultContext()
	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &client.Context{})
	ctx = context.WithValue(ctx, server.ServerContextKey, srvCtx)

	rootCmd.PersistentFlags().String(flags.FlagLogLevel, zerolog.InfoLevel.String(), "The logging level (trace|debug|info|warn|error|fatal|panic)")
	rootCmd.PersistentFlags().String(flags.FlagLogFormat, tmcfg.LogFormatPlain, "The logging format (json|plain)")

	executor := tmcli.PrepareBaseCmd(rootCmd, "", defaultHome)
	return executor.ExecuteContext(ctx)
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig params.EncodingConfig) {
	rootCmd.AddCommand(
		genutilcli.InitCmd(app.ModuleBasics, app.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome),
		genutilcli.MigrateGenesisCmd(),
		genutilcli.GenTxCmd(app.ModuleBasics, encodingConfig.TxConfig, banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome),
		genutilcli.ValidateGenesisCmd(app.ModuleBasics),
		AddGenesisAccountCmd(app.DefaultNodeHome),
		AddGenesisCertifierCmd(app.DefaultNodeHome),
		AddGenesisShieldAdminCmd(app.DefaultNodeHome),
		RotateValKeysCmd(),
		tmcli.NewCompletionCmd(rootCmd, true),
		testnetCmd(app.ModuleBasics, banktypes.GenesisBalancesIterator{}),
		tmcmds.RollbackStateCmd,
		debug.Cmd(),
		config.Cmd(),
	)

	rootCmd.Flags().StringP(shentuinit.DocFlag, shentuinit.DocFlagAbbr, "", shentuinit.DocFlagUsage)

	server.AddCommands(rootCmd, app.DefaultNodeHome, newApp, createSimappAndExport, addModuleInitFlags)

	// add keybase, auxiliary RPC, query, and tx child commands
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		queryCommand(),
		txCommand(),
		keys.Commands(app.DefaultNodeHome),
	)

	// add rosetta
	rootCmd.AddCommand(server.RosettaCommand(encodingConfig.InterfaceRegistry, encodingConfig.Marshaler))
}

func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		sdkauthcli.GetAccountCmd(),
		rpc.ValidatorCommand(),
		rpc.BlockCommand(),
		sdkauthcli.QueryTxsByEventsCmd(),
		sdkauthcli.QueryTxCmd(),
	)

	app.ModuleBasics.AddQueryCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		bankcli.LockedSendTxCmd(),
		sdkbankcli.NewSendTxCmd(),
		flags.LineBreak,
		sdkauthcli.GetSignCommand(),
		sdkauthcli.GetSignBatchCommand(),
		sdkauthcli.GetMultiSignCommand(),
		sdkauthcli.GetValidateSignaturesCommand(),
		sdkauthcli.GetBroadcastCommand(),
		sdkauthcli.GetEncodeCommand(),
		sdkauthcli.GetDecodeCommand(),
		authcli.GetCmdUnlock(),
		flags.LineBreak,
		vestingcli.GetTxCmd(),
		flags.LineBreak,
		cvmcli.GetCmdCall(),
		cvmcli.GetCmdDeploy(),
		flags.LineBreak,
	)

	app.ModuleBasics.AddTxCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
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

	return app.NewShentuApp(
		logger, db, traceStore, true, skipUpgradeHeights,
		cast.ToString(appOpts.Get(flags.FlagHome)),
		cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
		app.MakeEncodingConfig(), // Ideally, we would reuse the one created by NewRootCmd.
		appOpts,
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(cast.ToString(appOpts.Get(server.FlagMinGasPrices))),
		baseapp.SetHaltHeight(cast.ToUint64(appOpts.Get(server.FlagHaltHeight))),
		baseapp.SetHaltTime(cast.ToUint64(appOpts.Get(server.FlagHaltTime))),
		baseapp.SetMinRetainBlocks(cast.ToUint64(appOpts.Get(server.FlagMinRetainBlocks))),
		baseapp.SetInterBlockCache(cache),
		baseapp.SetTrace(cast.ToBool(appOpts.Get(server.FlagTrace))),
		baseapp.SetIndexEvents(cast.ToStringSlice(appOpts.Get(server.FlagIndexEvents))),
		baseapp.SetSnapshotStore(snapshotStore),
		baseapp.SetSnapshotInterval(cast.ToUint64(appOpts.Get(server.FlagStateSyncSnapshotInterval))),
		baseapp.SetSnapshotKeepRecent(cast.ToUint32(appOpts.Get(server.FlagStateSyncSnapshotKeepRecent))),
		baseapp.SetIAVLCacheSize(cast.ToInt(appOpts.Get(server.FlagIAVLCacheSize))),
		baseapp.SetIAVLDisableFastNode(cast.ToBool(appOpts.Get(server.FlagIAVLFastNode))),
	)
}

func createSimappAndExport(logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailAllowedAddrs []string,
	appOpts servertypes.AppOptions) (servertypes.ExportedApp, error) {
	encCfg := app.MakeEncodingConfig() // Ideally, we would reuse the one created by NewRootCmd.
	encCfg.Marshaler = codec.NewProtoCodec(encCfg.InterfaceRegistry)
	var gaiaApp *app.ShentuApp
	if height != -1 {
		gaiaApp = app.NewShentuApp(logger, db, traceStore, false, map[int64]bool{}, "", uint(1), encCfg, appOpts)

		if err := gaiaApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	} else {
		gaiaApp = app.NewShentuApp(logger, db, traceStore, true, map[int64]bool{}, "", uint(1), encCfg, appOpts)
	}

	return gaiaApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs)
}
