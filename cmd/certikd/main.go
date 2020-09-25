package main

import (
	"encoding/json"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/server"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
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
	rootCmd.AddCommand(AddGenesisAccountCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome))
	rootCmd.AddCommand(AddGenesisCertifierCmd(ctx, cdc))
	rootCmd.AddCommand(AddGenesisShieldOperatorCmd(ctx, cdc))
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

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return app.NewCertiKApp(logger, db, traceStore, true, map[int64]bool{}, uint(1),
		baseapp.SetPruning(storetypes.NewPruningOptionsFromString(viper.GetString("pruning"))),
		baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
		baseapp.SetHaltHeight(uint64(viper.GetInt(server.FlagHaltHeight))),
	)
}

func exportAppStateAndValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	if height != -1 {
		cApp := app.NewCertiKApp(logger, db, traceStore, false, map[int64]bool{}, uint(1))
		err := cApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return cApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}
	cApp := app.NewCertiKApp(logger, db, traceStore, true, map[int64]bool{}, uint(1))
	return cApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}
