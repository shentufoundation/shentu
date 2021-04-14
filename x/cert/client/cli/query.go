// Package cli defines the CLI services for the cert module.
package cli

import (
	"context"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/cert/types"
)

// GetQueryCmd returns the cli query commands for the certification module.
func GetQueryCmd() *cobra.Command {
	// Group cert queries under a subcommand.
	certQueryCmds := &cobra.Command{
		Use:   "cert",
		Short: "Querying commands for the certification module",
	}

	certQueryCmds.AddCommand(
		GetCmdCertifier(),
		GetCmdCertifiers(),
		GetCmdValidator(),
		GetCmdValidators(),
		GetCmdPlatform(),
		GetCmdCertificate(),
		GetCmdCertificates(),
	)

	return certQueryCmds
}

// GetCmdCertifier returns the certifier query command.
func GetCmdCertifier() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certifier <address>",
		Short: "Get certifier information",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			var req types.QueryCertifierRequest
			req.Alias = viper.GetString(FlagAlias)
			if len(args) > 0 {
				req.Address = args[0]
			}

			queryClient := types.NewQueryClient(cliCtx)
			res, err := queryClient.Certifier(
				context.Background(),
				&req,
			)
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}

	cmd.Flags().String(FlagAlias, "", "use alias to query the certifier info")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdCertifiers returns all certifier query command
func GetCmdCertifiers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certifiers",
		Short: "Get certifiers information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(cliCtx)
			res, err := queryClient.Certifiers(context.Background(), &types.QueryCertifiersRequest{})
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdValidator returns the validator certification query command.
func GetCmdValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator <pubkey>",
		Short: "Get validator certification information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			_, err = sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.Validator(context.Background(), &types.QueryValidatorRequest{Pubkey: args[0]})
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdValidators returns all validators certification query command
func GetCmdValidators() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validators",
		Short: "Get validators certification information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			res, err := queryClient.Validators(context.Background(), &types.QueryValidatorsRequest{})
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdCertificate returns the certificate query command.
func GetCmdCertificate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certificate <certificate id>",
		Short: "Get certificate information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			certificateID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			res, err := queryClient.Certificate(context.Background(), &types.QueryCertificateRequest{CertificateId: certificateID})
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdCertificates returns certificates query command
func GetCmdCertificates() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certificates [<flags>]",
		Short: "Get certificates information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.Certificates(
				cmd.Context(),
				&types.QueryCertificatesRequest{
					Certifier:       viper.GetString(FlagCertifier),
					CertificateType: viper.GetString(FlagCertType),
					Pagination:      pageReq,
				})
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}

	cmd.Flags().String(FlagCertifier, "", "certificates issued by certifier")
	cmd.Flags().String(FlagCertType, "", "certificates by type")
	flags.AddPaginationFlagsToCmd(cmd, "certificates")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdPlatform returns the validator host platform certification query command.
func GetCmdPlatform() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "platform <pubkey>",
		Short: "Get validator host platform certification information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			pk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, args[0])
			if err != nil {
				return err
			}
			var pkAny *codectypes.Any
			if pk != nil {
				var err error
				if pkAny, err = codectypes.NewAnyWithValue(pk); err != nil {
					return err
				}
			}

			res, err := queryClient.Platform(context.Background(), &types.QueryPlatformRequest{Pubkey: pkAny})
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
