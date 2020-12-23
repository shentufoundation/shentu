// Package cli defines the CLI services for the cert module.
package cli

import (
	"context"

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
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadQueryCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			res, err := queryClient.Certifier(
				context.Background(),
				&types.QueryCertifierRequest{Address: args[0], Alias: viper.GetString(FlagAlias)},
			)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(res)
		},
	}
	cmd.Flags().String(FlagAlias, "", "use alias to query the certifier info")
	return cmd
}

// GetCmdCertifiers returns all certifier query command
func GetCmdCertifiers() *cobra.Command {
	return &cobra.Command{
		Use:   "certifiers",
		Short: "Get certifiers information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadQueryCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			res, err := queryClient.Certifiers(context.Background(), &types.QueryCertifiersRequest{})
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(res)
		},
	}
}

// GetCmdValidator returns the validator certification query command.
func GetCmdValidator() *cobra.Command {
	return &cobra.Command{
		Use:   "validator <pubkey>",
		Short: "Get validator certification information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadQueryCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			pk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, args[0])
			if err != nil {
				return err
			}
			pkAny, err := codectypes.PackAny(pk)
			if err != nil {
				return err
			}

			res, err := queryClient.Validator(context.Background(), &types.QueryValidatorRequest{Pubkey: pkAny})
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(res)
		},
	}
}

// GetCmdValidators returns all validators certification query command
func GetCmdValidators() *cobra.Command {
	return &cobra.Command{
		Use:   "validators",
		Short: "Get validators certification information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadQueryCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			res, err := queryClient.Validators(context.Background(), &types.QueryValidatorsRequest{})
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(res)
		},
	}
}

// GetCmdCertificate returns the certificate query command.
func GetCmdCertificate() *cobra.Command {
	return &cobra.Command{
		Use:   "certificate <certificate id>",
		Short: "Get certificate information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadQueryCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			res, err := queryClient.Certificate(context.Background(), &types.QueryCertificateRequest{CertificateId: args[0]})
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(res)
		},
	}
}

// GetCmdCertificates returns certificates query command
func GetCmdCertificates() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certificates [<flags>]",
		Short: "Get certificates information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadQueryCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.Certificates(
				context.Background(),
				&types.QueryCertificatesRequest{
					Certifier:   viper.GetString(FlagCertifier),
					Content:     viper.GetString(FlagContent),
					ContentType: viper.GetString(FlagContentType),
					Pagination:  pageReq,
				})
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(res)
		},
	}
	cmd.Flags().String(FlagCertifier, "", "certificates issued by certifier")
	cmd.Flags().String(FlagContent, "", "certificates by request content")
	cmd.Flags().String(FlagContentType, "", "type of request content")
	flags.AddPaginationFlagsToCmd(cmd, "votes")
	return cmd
}

// GetCmdPlatform returns the validator host platform certification query command.
func GetCmdPlatform() *cobra.Command {
	return &cobra.Command{
		Use:   "platform <pubkey>",
		Short: "Get validator host platform certification information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadQueryCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			pk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, args[0])
			if err != nil {
				return err
			}
			pkAny, err := codectypes.PackAny(pk)
			if err != nil {
				return err
			}

			res, err := queryClient.Platform(context.Background(), &types.QueryPlatformRequest{Pubkey: pkAny})
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(res)
		},
	}
}
