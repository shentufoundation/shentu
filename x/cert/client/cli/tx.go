package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/version"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/certikfoundation/shentu/v2/x/cert/types"
)

const (
	FlagAlias        = "alias"
	FlagCertType     = "certificate-type"
	FlagCompiler     = "compiler"
	FlagBytecodeHash = "bytecode-hash"
	FlagDescription  = "description"
	FlagCertifier    = "certifier"
	FlagPage         = "page"
	FlagLimit        = "limit"
)

// NewTxCmd returns the transaction commands for the certification module.
func NewTxCmd() *cobra.Command {
	certTxCmds := &cobra.Command{
		Use:   "cert",
		Short: "Certification transactions subcommands",
	}

	certTxCmds.AddCommand(
		GetCmdCertifyPlatform(),
		GetCmdIssueCertificate(),
		GetCmdRevokeCertificate(),
	)

	return certTxCmds
}

// GetCmdIssueCertificate returns the certificate transaction command.
func GetCmdIssueCertificate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-certificate <certificate type> <request content> [<flags>]",
		Short: "Issue a certificate",
		Args:  cobra.ExactArgs(2),
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

			compiler, bytecodeHash := "", ""
			certificateTypeString := strings.ToLower(args[0])
			if certificateTypeString == "compilation" {
				compiler, bytecodeHash, err = parseCertifyCompilationFlags()
				if err != nil {
					return err
				}
			}
			description := viper.GetString(FlagDescription)
			content := types.AssembleContent(args[0], args[1])
			msg := types.NewMsgIssueCertificate(content, compiler, bytecodeHash, description, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	cmd.Flags().String(FlagCompiler, "", "compiler version")
	cmd.Flags().String(FlagBytecodeHash, "", "bytecode hash")
	cmd.Flags().String(FlagDescription, "", "description")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// parseCertifyCompilation parses flags for compilation certificate.
func parseCertifyCompilationFlags() (string, string, error) {
	compiler := viper.GetString(FlagCompiler)
	if compiler == "" {
		return "", "", fmt.Errorf("compiler version is required to issue a compilation certificate")
	}
	bytecodeHash := viper.GetString(FlagBytecodeHash)
	if bytecodeHash == "" {
		return "", "", fmt.Errorf("bytecode hash is required to issue a compilation certificate")
	}
	return compiler, bytecodeHash, nil
}

// GetCmdCertifyPlatform returns the validator host platform certification transaction command.
func GetCmdCertifyPlatform() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certify-platform <validator pubkey> <platform>",
		Short: "Certify a validator's host platform",
		Args:  cobra.ExactArgs(2),
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

			var validator cryptotypes.PubKey
			err = cliCtx.JSONCodec.UnmarshalJSON([]byte(args[0]), validator)
			if err != nil {
				return err
			}

			msg, err := types.NewMsgCertifyPlatform(from, validator, args[1])
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

// GetCmdRevokeCertificate returns the certificate revoke command
func GetCmdRevokeCertificate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke-certificate <certificateID> [<description>]",
		Short: "revoke a certificate",
		Args:  cobra.RangeArgs(1, 2),
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

			certificateID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			var description string
			if len(args) > 1 {
				description = args[1]
			}

			msg := types.NewMsgRevokeCertificate(from, certificateID, description)
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
