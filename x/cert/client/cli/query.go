// Package cli defines the CLI services for the cert module.
package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
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
