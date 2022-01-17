package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/certikfoundation/shentu/v2/x/bank/types"
)

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

			unlocker, _ := cmd.Flags().GetString(FlagUnlocker)
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
	cmd.Flags().AddFlagSet(FlagAddUnlocker())
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}