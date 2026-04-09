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

	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

const (
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

			txf, err := tx.NewFactoryCLI(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			txf = txf.WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

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
		Long: strings.TrimSpace(
			fmt.Sprintf(`Certify a validator's host platform
Example:
$ %s tx cert certify-platform '{\"@type\":\"/cosmos.crypto.secp256k1.PubKey\",\"key\":\"A/N4uA+0c7xg1nOI8eGSB5tYIiQXbLfAZif0/MnwdeJP\"}' test --from=<key_or_address>
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf, err := tx.NewFactoryCLI(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			txf = txf.WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			from := cliCtx.GetFromAddress()
			if err := txf.AccountRetriever().EnsureExists(cliCtx, from); err != nil {
				return err
			}

			var validator cryptotypes.PubKey
			err = cliCtx.Codec.UnmarshalInterfaceJSON([]byte(args[0]), &validator)
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

			txf, err := tx.NewFactoryCLI(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			txf = txf.WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

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
