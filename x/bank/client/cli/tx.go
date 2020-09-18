package cli

import (
	"bufio"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/certikfoundation/shentu/x/bank/internal/types"
)

// SendLockTxCmd sends coins to a manual vesting account
// and have them vesting.
func SendLockTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-lock [from_key_or_address] [to_address] [amount]",
		Short: "Send coins and have them locked (vesting).",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := authtxb.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			accGetter := authtxb.NewAccountRetriever(cliCtx)

			if _, err := accGetter.GetAccount(cliCtx.GetFromAddress()); err != nil {
				return err
			}

			to, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoin(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgSendLock(cliCtx.GetFromAddress(), to, amount)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd = flags.PostCommands(cmd)[0]
	return cmd
}
