package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/gov"

	"github.com/certikfoundation/shentu/x/cert/internal/types"
)

const (
	FlagAlias        = "alias"
	FlagContentType  = "content-type"
	FlagContent      = "content"
	FlagCompiler     = "compiler"
	FlagBytecodeHash = "bytecode-hash"
	FlagDescription  = "description"
	FlagCertifier    = "certifier"
	FlagPage         = "page"
	FlagLimit        = "limit"
)

// GetTxCmd returns the transaction commands for the certification module.
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	certTxCmds := &cobra.Command{
		Use:   "cert",
		Short: "Certification transactions subcommands",
	}

	certTxCmds.AddCommand(flags.PostCommands(
		GetCmdCertifyValidator(cdc),
		GetCmdDecertifyValidator(cdc),
		GetCmdCertifyPlatform(cdc),
		GetCmdIssueCertificate(cdc),
		GetCmdRevokeCertificate(cdc),
	)...)

	return certTxCmds
}

// GetCmdCertifyValidator returns the validator certification transaction command.
func GetCmdCertifyValidator(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "certify-validator <validator pubkey>",
		Short: "Certify a validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := authtxb.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			accGetter := authtxb.NewAccountRetriever(cliCtx)

			if _, err := accGetter.GetAccount(cliCtx.GetFromAddress()); err != nil {
				return err
			}

			validator, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, args[0])
			if err != nil {
				return err
			}
			msg := types.NewMsgCertifyValidator(cliCtx.GetFromAddress(), validator)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdDecertifyValidator returns the validator de-certification tx command.
func GetCmdDecertifyValidator(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "decertify-validator <validator pubkey>",
		Short: "De-certify a validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := authtxb.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			accGetter := authtxb.NewAccountRetriever(cliCtx)

			if _, err := accGetter.GetAccount(cliCtx.GetFromAddress()); err != nil {
				return err
			}

			validator, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, args[0])
			if err != nil {
				return err
			}
			msg := types.NewMsgDecertifyValidator(cliCtx.GetFromAddress(), validator)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdIssueCertificate returns the certificate transaction command.
func GetCmdIssueCertificate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-certificate <certificate type> <request content type> <request content> [<flags>]",
		Short: "Issue a certificate",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := authtxb.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			from := cliCtx.GetFromAddress()
			accGetter := authtxb.NewAccountRetriever(cliCtx)
			if err := accGetter.EnsureExists(from); err != nil {
				return err
			}

			certificateTypeString := strings.ToLower(args[0])
			switch certificateTypeString {
			case "compilation":
				contentType := types.RequestContentTypeFromString(args[1])
				if contentType != types.RequestContentTypeSourceCodeHash {
					return types.ErrInvalidRequestContentType
				}
				compiler, bytecodeHash, description, err := parseCertifyCompilationFlags()
				if err != nil {
					return err
				}
				msg := types.NewMsgCertifyCompilation(args[2], compiler, bytecodeHash, description, from)
				if err := msg.ValidateBasic(); err != nil {
					return err
				}
				return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})

			default:
				description := viper.GetString(FlagDescription)
				msg := types.NewMsgCertifyGeneral(certificateTypeString, args[1], args[2], description, from)
				if err := msg.ValidateBasic(); err != nil {
					return err
				}
				return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
			}
		},
	}

	cmd.Flags().String(FlagCompiler, "", "compiler version")
	cmd.Flags().String(FlagBytecodeHash, "", "bytecode hash")
	cmd.Flags().String(FlagDescription, "", "description")

	return cmd
}

// parseCertifyCompilation parses flags for compilation certificate.
func parseCertifyCompilationFlags() (string, string, string, error) {
	compiler := viper.GetString(FlagCompiler)
	if compiler == "" {
		return "", "", "", fmt.Errorf("compiler version is required to issue a compilation certificate")
	}
	bytecodeHash := viper.GetString(FlagBytecodeHash)
	if bytecodeHash == "" {
		return "", "", "", fmt.Errorf("bytecode hash is required to issue a compilation certificate")
	}
	description := viper.GetString(FlagDescription)
	return compiler, bytecodeHash, description, nil
}

// GetCmdCertifyPlatform returns the validator host platform certification transaction command.
func GetCmdCertifyPlatform(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "certify-platform <validator pubkey> <platform>",
		Short: "Certify a validator's host platform",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := authtxb.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			accGetter := authtxb.NewAccountRetriever(cliCtx)
			if _, err := accGetter.GetAccount(cliCtx.GetFromAddress()); err != nil {
				return err
			}

			validator, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgCertifyPlatform(cliCtx.GetFromAddress(), validator, args[1])
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdRevokeCertificate returns the certificate revoke command
func GetCmdRevokeCertificate(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "revoke-certificate <certificateID> [<description>]",
		Short: "revoke a certificate",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := authtxb.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			accGetter := authtxb.NewAccountRetriever(cliCtx)
			description := ""

			if len(args) > 1 {
				description = args[1]
			}

			if _, err := accGetter.GetAccount(cliCtx.GetFromAddress()); err != nil {
				return err
			}

			msg := types.NewMsgRevokeCertificate(cliCtx.GetFromAddress(), types.CertificateID(args[0]), description)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdSubmitProposal implements the command to submit a certifier-update proposal
func GetCmdSubmitProposal(cdc *codec.Codec) *cobra.Command {
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
  "certifier": "certik1s5afhd6gxevu37mkqcvvsj8qeylhn0rz46zdlq",
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
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			proposal, err := ParseCertifierUpdateProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}

			from := cliCtx.GetFromAddress()
			content := types.NewCertifierUpdateProposal(
				proposal.Title,
				proposal.Description,
				proposal.Certifier,
				proposal.Alias,
				from,
				proposal.AddOrRemove,
			)

			msg := gov.NewMsgSubmitProposal(content, proposal.Deposit, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
