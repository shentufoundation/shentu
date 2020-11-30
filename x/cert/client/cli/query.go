// Package cli defines the CLI services for the cert module.
package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/cert/internal/keeper"
	"github.com/certikfoundation/shentu/x/cert/internal/types"
)

// GetQueryCmd returns the cli query commands for the certification module.
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group cert queries under a subcommand.
	certQueryCmds := &cobra.Command{
		Use:   "cert",
		Short: "Querying commands for the certification module",
	}

	certQueryCmds.AddCommand(flags.GetCommands(
		GetCmdCertifier(queryRoute, cdc),
		GetCmdCertifiers(queryRoute, cdc),
		GetCmdValidator(queryRoute, cdc),
		GetCmdValidators(queryRoute, cdc),
		GetCmdPlatform(queryRoute, cdc),
		GetCmdCertificate(queryRoute, cdc),
		GetCmdCertificates(queryRoute, cdc),
	)...)

	return certQueryCmds
}

// GetCmdCertifier returns the certifier query command.
func GetCmdCertifier(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certifier <address>",
		Short: "Get certifier information",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var addr string
			var res []byte
			var err error
			if len(args) > 0 {
				addr = args[0]
				res, _, err = cliCtx.QueryWithData(fmt.Sprintf("custom/%s/certifier/%s", queryRoute, addr), nil)
				if err != nil {
					return err
				}
			} else {
				alias := viper.GetString(FlagAlias)
				if alias == "" {
					return fmt.Errorf("require address or alias to query certifiers")
				}
				res, _, err = cliCtx.QueryWithData(fmt.Sprintf("custom/%s/certifieralias/%s", queryRoute, alias), nil)
				if err != nil {
					return err
				}
			}

			var out types.Certifier
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
	cmd.Flags().String(FlagAlias, "", "use alias to query the certifier info")
	return cmd
}

// GetCmdCertifiers returns all certifier query command
func GetCmdCertifiers(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "certifiers",
		Short: "Get certifiers information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/certifiers", queryRoute), nil)
			if err != nil {
				return err
			}
			var out types.QueryResCertifiers
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdValidator returns the validator certification query command.
func GetCmdValidator(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "validator <pubkey>",
		Short: "Get validator certification information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			key := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/validator/%s", queryRoute, key), nil)
			if err != nil {
				return err
			}
			var out types.QueryResValidator
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdValidators returns all validators certification query command
func GetCmdValidators(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "validators",
		Short: "Get validators certification information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/validators", queryRoute), nil)
			if err != nil {
				return err
			}
			var out types.QueryResValidators
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdCertificate returns the certificate query command.
func GetCmdCertificate(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "certificate <certificate id>",
		Short: "Get certificate information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			certificateID := args[0]
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/certificate/%s", queryRoute, certificateID), nil)
			if err != nil {
				return err
			}
			var certificate keeper.QueryResCertificate
			cdc.MustUnmarshalJSON(res, &certificate)
			return cliCtx.PrintOutput(certificate)
		},
	}
}

// GetCmdCertificates returns certificates query command
func GetCmdCertificates(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certificates [<flags>]",
		Short: "Get certificates information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var (
				err              error
				certifierAddress sdk.AccAddress
			)

			if certifier := viper.GetString(FlagCertifier); certifier != "" {
				certifierAddress, err = sdk.AccAddressFromBech32(certifier)
				if err != nil {
					return err
				}
			}

			contentTypeString := viper.GetString(FlagContentType)
			content := viper.GetString(FlagContent)

			page := viper.GetInt(FlagPage)
			limit := viper.GetInt(FlagLimit)
			params := types.NewQueryCertificatesParams(page, limit, certifierAddress, contentTypeString, content)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/certificates", queryRoute), bz)
			if err != nil {
				return err
			}
			var out keeper.QueryResCertificates
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
	cmd.Flags().String(FlagCertifier, "", "certificates issued by certifier")
	cmd.Flags().String(FlagContent, "", "certificates by request content")
	cmd.Flags().String(FlagContentType, "", "type of request content")
	cmd.Flags().Int(FlagPage, 1, "pagination page of certificates to to query for")
	cmd.Flags().Int(FlagLimit, 100, "pagination limit of certificates to query for")
	return cmd
}

// GetCmdPlatform returns the validator host platform certification query command.
func GetCmdPlatform(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "platform <pubkey>",
		Short: "Get validator host platform certification information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			key := args[0]
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/platform/%s", queryRoute, key), nil)
			if err != nil {
				return err
			}

			if res == nil {
				return fmt.Errorf("this platform is not certified")
			}

			var out types.QueryResPlatform
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
