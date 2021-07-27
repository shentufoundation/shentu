package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/certikfoundation/shentu/x/cert/types"
)

const (
	FlagAlias     = "alias"
	FlagCertifier = "certifier"
)

// NewTxCmd returns the transaction commands for the certification module.
func NewTxCmd() *cobra.Command {
	certTxCmds := &cobra.Command{
		Use:   "cert",
		Short: "Certification transactions subcommands",
	}

	certTxCmds.AddCommand(
		GetCmdCertifyValidator(),
		GetCmdDecertifyValidator(),
		GetCmdSubmitProposal(),
	)

	return certTxCmds
}

// GetCmdCertifyValidator returns the validator certification transaction command.
func GetCmdCertifyValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certify-validator <validator pubkey>",
		Short: "Certify a validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			from := cliCtx.GetFromAddress()
			if err := txf.AccountRetriever().EnsureExists(cliCtx, from); err != nil {
				return err
			}

			validator, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, args[0])
			if err != nil {
				return err
			}
			msg, err := types.NewMsgCertifyValidator(from, validator)
			if err != nil {
				return err
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdDecertifyValidator returns the validator de-certification tx command.
func GetCmdDecertifyValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decertify-validator <validator pubkey>",
		Short: "De-certify a validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			from := cliCtx.GetFromAddress()
			if err := txf.AccountRetriever().EnsureExists(cliCtx, from); err != nil {
				return err
			}

			validator, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, args[0])
			if err != nil {
				return err
			}
			msg, err := types.NewMsgDecertifyValidator(from, validator)
			if err != nil {
				return err
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdSubmitProposal implements the command to submit a certifier-update proposal
func GetCmdSubmitProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certifier-update [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a certifier update proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a certifier update proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.
Example:
$ %s tx gov submit-proposal certifier-update <path/to/proposal.json> --from=<key_or_address>
Where proposal.json contains:
{
  "title": "New Certifier, Joe Shmoe",
  "description": "Why we should make Joe Shmoe a certifier",
  "certifier": "certik1fdyv6hpukqj6kqdtwc42qacq9lpxm0pn85w6l9",
  "add_or_remove": "add",
  "alias": "joe",
  "deposit": [
    {
      "denom": "ctk",
      "amount": "100"
    }
  ]
}
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			from := cliCtx.GetFromAddress()
			if err := txf.AccountRetriever().EnsureExists(cliCtx, from); err != nil {
				return err
			}

			proposal, err := ParseCertifierUpdateProposalJSON(cliCtx.LegacyAmino, args[0])
			if err != nil {
				return err
			}

			content := types.NewCertifierUpdateProposal(
				proposal.Title,
				proposal.Description,
				proposal.Certifier,
				proposal.Alias,
				from,
				proposal.AddOrRemove,
			)

			msg, err := govtypes.NewMsgSubmitProposal(content, proposal.Deposit, from)
			if err != nil {
				return err
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	return cmd
}
