package main

import (
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	authcli "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankcli "github.com/cosmos/cosmos-sdk/x/bank/client/cli"

	"github.com/certikfoundation/shentu/app"
	"github.com/certikfoundation/shentu/client/lcd"
	certikinit "github.com/certikfoundation/shentu/cmd/init"
	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/toolsets/oracle-operator"
	cvmcli "github.com/certikfoundation/shentu/x/cvm/client/cli"
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

	rootCmd := &cobra.Command{
		Use:   "certikcli",
		Short: "CertiK Client",
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

	// add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(flags.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		client.ConfigCmd(app.DefaultCLIHome),
		queryCmd(cdc),
		txCmd(cdc),
		flags.LineBreak,
		lcd.ServeCommand(cdc, registerRoutes),
		oracle.ServeCommand(cdc),
		flags.LineBreak,
		keys.Commands(),
		flags.LineBreak,
		version.Cmd,
	)

	executor := cli.PrepareMainCmd(rootCmd, "NS", app.DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func registerRoutes(rs *lcd.RestServer) {
	rpc.RegisterRPCRoutes(rs.CliCtx, rs.Mux)
	client.RegisterRoutes(rs.CliCtx, rs.Mux)
	authrest.RegisterTxRoutes(rs.CliCtx, rs.Mux)
	app.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
}

func queryCmd(cdc *amino.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		rpc.ValidatorCommand(cdc),
		rpc.BlockCommand(),
		authcli.QueryTxsByEventsCmd(cdc),
		cvmcli.QueryTxCmd(cdc),
		flags.LineBreak,
		authcli.GetAccountCmd(cdc),
	)

	// add modules' query commands
	app.ModuleBasics.AddQueryCommands(queryCmd, cdc)

	return queryCmd
}

func txCmd(cdc *amino.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(
		bankcli.SendTxCmd(cdc),
		flags.LineBreak,
		authcli.GetSignCommand(cdc),
		authcli.GetEncodeCommand(cdc),
		authcli.GetDecodeCommand(cdc),
		authcli.GetMultiSignCommand(cdc),
		authcli.GetBroadcastCommand(cdc),
		flags.LineBreak,
		cvmcli.GetCmdCall(cdc),
		cvmcli.GetCmdDeploy(cdc),
		flags.LineBreak,
	)

	// add modules' tx commands
	app.ModuleBasics.AddTxCommands(txCmd, cdc)

	// remove auth and bank commands as they're mounted under the root tx command
	var cmdsToRemove []*cobra.Command

	for _, cmd := range txCmd.Commands() {
		if cmd.Use == bank.ModuleName {
			cmdsToRemove = append(cmdsToRemove, cmd)
		}
	}

	txCmd.RemoveCommand(cmdsToRemove...)

	return txCmd
}

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}

	cfgFile := path.Join(home, "config", "config.toml")
	if _, err := os.Stat(cfgFile); err == nil {
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}
	if err := viper.BindPFlag(flags.FlagChainID, cmd.PersistentFlags().Lookup(flags.FlagChainID)); err != nil {
		return err
	}
	if err := viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}
