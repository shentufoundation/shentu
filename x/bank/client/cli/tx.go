package cli

import (
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/bank/internal/types"
)

const (
	FlagUnlocker = "unlocker"
)

// LockedSendTxCmd sends coins to a manual vesting account
// and have them vesting.
func LockedSendTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locked-send [from_key_or_address] [to_address] [amount]",
		Short: "Send coins and have them locked (vesting).",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[0])

			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadTxCommandFlags(cliCtx, cmd.Flags())

			to, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoins(args[2])
			if err != nil {
				return err
			}

			unlocker, err := sdk.AccAddressFromBech32(viper.GetString(FlagUnlocker))
			if err != nil {
				return err
			}

			msg := types.NewMsgLockedSend(cliCtx.GetFromAddress(), to, unlocker, coins)
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(FlagUnlocker, "", "unlocker when initializing a new manual vesting account")
	return cmd
}
