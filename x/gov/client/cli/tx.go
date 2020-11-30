package cli

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authUtils "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/certikfoundation/shentu/x/gov/client/utils"
)

// Proposal flags
const (
	flagVoter     = "voter"
	flagDepositor = "depositor"
	flagStatus    = "status"
)

// GetTxCmd returns the transaction commands for this module
// governance ModuleClient is slightly different from other ModuleClients in that
// it contains a slice of "proposal" child commands. These commands are respective
// to proposal type handlers that are implemented in other modules but are mounted
// under the governance CLI (eg. parameter change proposals).
func GetTxCmd(storeKey string, cdc *codec.Codec, pcmds []*cobra.Command) *cobra.Command {
	govTxCmd := &cobra.Command{
		Use:                        govTypes.ModuleName,
		Short:                      "Governance transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmdSubmitProp := cli.GetCmdSubmitProposal(cdc)
	cmdSubmitProp.Long = strings.TrimSpace(
		fmt.Sprintf(`Submit a proposal along with an initial deposit.
Proposal title, description, type and deposit can be given directly or through a proposal JSON file.

Example:
$ %[1]s tx gov submit-proposal --proposal="path/to/proposal.json" --from mykey

Where proposal.json contains:

{
  "title": "Test Proposal",
  "description": "My awesome proposal",
  "type": "Text",
  "deposit": "10ctk"
}

Which is equivalent to:

$ %[1]s tx gov submit-proposal --title="Test Proposal" --description="My awesome proposal" --type="Text" --deposit="10ctk" --from mykey
`,
			version.ClientName,
		),
	)
	for _, pcmd := range pcmds {
		cmdSubmitProp.AddCommand(flags.PostCommands(pcmd)[0])
	}

	govTxCmd.AddCommand(flags.PostCommands(
		cli.GetCmdDeposit(cdc),
		GetCmdVote(cdc),
		cmdSubmitProp,
	)...)

	return govTxCmd
}

// GetCmdVote implements creating a new vote command.
func GetCmdVote(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "vote [proposal-id] [option]",
		Args:  cobra.ExactArgs(2),
		Short: "Vote for an active proposal, options: yes/no/no_with_veto/abstain",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a vote for an active proposal. You can
find the proposal-id by running "%[1]s query gov proposals".


Example:
$ %[1]s tx gov vote 1 yes --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authUtils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			// get voting address
			from := cliCtx.GetFromAddress()

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s is not a valid int, please input a valid proposal-id", args[0])
			}

			// find out which vote option user chose
			byteVoteOption, err := govTypes.VoteOptionFromString(utils.NormalizeVoteOption(args[1]))
			if err != nil {
				return err
			}

			// build vote message and run basic validation
			msg := govTypes.NewMsgVote(from, proposalID, byteVoteOption)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return authUtils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
