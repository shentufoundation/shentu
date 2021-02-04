package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/bank/types"
)

const (
	FlagUnlocker = "unlocker"
)

// // NewTxCmd returns a root CLI command handler for all x/bank transaction commands.
// func NewTxCmd() *cobra.Command {
// 	txCmd := &cobra.Command{
// 		Use:                        banktypes.ModuleName,
// 		Short:                      "Bank transaction subcommands",
// 		DisableFlagParsing:         true,
// 		SuggestionsMinimumDistance: 2,
// 		RunE:                       client.ValidateCmd,
// 	}

// 	txCmd.AddCommand(LockedSendTxCmd())

// 	return txCmd
// }

// LockedSendTxCmd sends coins to a manual vesting account
// and have them vesting.
func LockedSendTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locked-send [from_key_or_address] [to_address] [amount]",
		Short: "Send coins and have them locked (vesting).",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[0])

			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			to, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoinsNormalized(args[2])
			if err != nil {
				return err
			}

			unlocker := viper.GetString(FlagUnlocker)
			if unlocker != "" {
				_, err = sdk.AccAddressFromBech32(unlocker)
				if err != nil {
					return err
				}
			}

			msg := types.NewMsgLockedSend(cliCtx.GetFromAddress(), to, unlocker, coins)
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(FlagUnlocker, "", "unlocker when initializing a new manual vesting account")
	return cmd
}
