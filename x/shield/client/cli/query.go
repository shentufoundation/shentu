package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/shentufoundation/shentu/v2/x/shield/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	shieldQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the shield module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	shieldQueryCmd.AddCommand(
		GetCmdProvider(),
		GetCmdProviders(),
		GetCmdStatus(),
	)

	return shieldQueryCmd
}

// GetCmdProvider returns the command for querying a provider.
func GetCmdProvider() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "provider [provider_address]",
		Short: "get provider information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.Provider(
				cmd.Context(),
				&types.QueryProviderRequest{Address: address.String()},
			)
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdProviders returns the command for querying all providers.
func GetCmdProviders() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "providers",
		Args:  cobra.ExactArgs(0),
		Short: "query all providers",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query providers with pagination parameters

Example:
$ %[1]s query shield providers
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			res, err := queryClient.Providers(cmd.Context(), &types.QueryProvidersRequest{})
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdStatus returns the command for querying shield status.
func GetCmdStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "get shield status",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCtx)

			res, err := queryClient.ShieldStatus(cmd.Context(), &types.QueryShieldStatusRequest{})
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
