// Package oracle defines oracle-operator
package oracle

import (
	"bufio"
	"os"

	"github.com/spf13/cobra"

	tmconfig "github.com/tendermint/tendermint/config"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/certikfoundation/shentu/toolsets/oracle-operator/runner"
	"github.com/certikfoundation/shentu/toolsets/oracle-operator/types"
)

// start starts the service.
func start(ctx types.Context) {
	for err := range runner.Start(ctx.WithLoggerLabels("module", runner.Name())) {
		ctx.Logger().Error(err.Error())
		os.Exit(1)
	}
}

// ServeCommand will start the oracle operator as a blocking process.
func ServeCommand(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oracle-operator",
		Short: "Start oracle operator",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			cliCtx.SkipConfirm = true // TODO: new cosmos version
			accGetter := authtxb.NewAccountRetriever(cliCtx)
			if _, err := accGetter.GetAccount(cliCtx.GetFromAddress()); err != nil {
				return err
			}

			ctx, err := types.NewContextWithDefaultConfigAndLogger()
			if err != nil {
				return err
			}
			ctx = ctx.WithClientContext(&cliCtx).WithTxBuilder(&txBldr)

			if err := serve(ctx); err != nil {
				return err
			}
			return nil
		},
	}

	return registerFlags(cmd)
}

// registerFlags registers additional flags to the command.
func registerFlags(cmd *cobra.Command) *cobra.Command {
	cmd = flags.PostCommands(cmd)[0]
	cmd.Flags().Uint(flags.FlagRPCReadTimeout, 10, "the RPC read timeout (in seconds)")
	cmd.Flags().Uint(flags.FlagRPCWriteTimeout, 10, "the RPC write timeout (in seconds)")
	cmd.Flags().String(types.FlagLogLevel, tmconfig.DefaultLogLevel(), "log level")
	return cmd
}

// serve sets up operator runner running environment.
func serve(ctx types.Context) error {
	done := make(chan struct{})
	panicChan := make(chan interface{}, 1)

	server.TrapSignal(func() {
		done <- struct{}{}
		ctx.Logger().Error("Stopping oracle-operator...")
	})

	ctx.Logger().Info("Starting oracle-operator...")
	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		start(ctx)
	}()

	defer close(done)
	select {
	case p := <-panicChan:
		panic(p)
	case <-done:
		ctx.Logger().Info("Stopping oracle-operator...")
	case <-ctx.Context().Done():
		ctx.Logger().Info("Stopping oracle-operator...")
	}
	return nil
}
