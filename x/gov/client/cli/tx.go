package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/shentufoundation/shentu/v2/x/gov/client/utils"
)

// Proposal flags
const (
	flagVoter     = "voter"
	flagDepositor = "depositor"
	flagStatus    = "status"
)

// NewTxCmd returns the transaction commands for gov module.
func NewTxCmd(pcmds []*cobra.Command) *cobra.Command {
	govTxCmd := &cobra.Command{
		Use:                        govTypes.ModuleName,
		Short:                      "Governance transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmdSubmitProp := cli.NewCmdSubmitProposal()
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
			version.AppName,
		),
	)
	for _, pcmd := range pcmds {
		flags.AddTxFlagsToCmd(pcmd)
		cmdSubmitProp.AddCommand(pcmd)
	}

	govTxCmd.AddCommand(
		cli.NewCmdDeposit(),
		NewCmdVote(),
		cli.NewCmdWeightedVote(),
		cmdSubmitProp,
	)

	return govTxCmd
}

// NewCmdVote implements creating a new vote command.
func NewCmdVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote [proposal-id] [option]",
		Args:  cobra.ExactArgs(2),
		Short: "Vote for an active proposal, options: yes/no/no_with_veto/abstain",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a vote for an active proposal. You can
find the proposal-id by running "%[1]s query gov proposals".


Example:
$ %[1]s tx gov vote 1 yes --from mykey
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

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

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
