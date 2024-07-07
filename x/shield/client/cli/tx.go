package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/shentufoundation/shentu/v2/x/shield/types"
)

// NewTxCmd returns the transaction commands for this module.
func NewTxCmd() *cobra.Command {
	shieldTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Shield transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	shieldTxCmd.AddCommand(
		GetCmdWithdrawRewards(),
	)

	return shieldTxCmd
}

// GetCmdWithdrawRewards implements command for requesting to withdraw native tokens rewards.
func GetCmdWithdrawRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-rewards",
		Short: "withdraw CTK rewards",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			fromAddr := cliCtx.GetFromAddress()

			msg := types.NewMsgWithdrawRewards(fromAddr)

			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
