package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/certikfoundation/shentu/x/shield/types"
)

var (
	flagNativeDeposit  = "native-deposit"
	flagForeignDeposit = "foreign-deposit"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	shieldTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Shield transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	shieldTxCmd.AddCommand(flags.PostCommands(
		GetCmdCreateShield(cdc),
	)...)

	return shieldTxCmd
}

// GetCmdCreateShield implements the create pool command handler.
func GetCmdCreateShield(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pool",
		Args:  cobra.ExactArgs(2),
		Short: "create new Shield pool initialized with an validator address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a shield pool. Can only be executed from the shield operator address.

Example:
$ %s tx shield create-pool <coverage> <sponsor> --native-deposit <ctk deposit> --foreign-deposit <external deposit>
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			coverage, err := sdk.ParseCoins(args[0])
			if err != nil {
				return err
			}

			sponsor := args[1]

			nativeDeposit, err := sdk.ParseCoins(viper.GetString(flagNativeDeposit))
			if err != nil {
				return err
			}

			foreignDeposit, err := sdk.ParseCoins(viper.GetString(flagForeignDeposit))
			if err != nil {
				return err
			}

			deposit := types.MixedCoins{
				Native:  nativeDeposit,
				Foreign: types.ForeignCoins(foreignDeposit),
			}

			msg, err := types.NewMsgCreatePool(fromAddr, coverage, deposit, sponsor)
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagNativeDeposit, "", "native deposit")
	cmd.Flags().String(flagForeignDeposit, "", "foreign deposit")

	return cmd
}
