package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/shentufoundation/shentu/v2/x/blacklist/types"
)

// GetQueryCmd returns the root query command for the blacklist module.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the blacklist module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetEntriesCmd(),
		GetEntryCmd(),
	)

	return cmd
}

// GetEntriesCmd returns the command to list all blacklisted addresses.
func GetEntriesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "entries",
		Short: "List all blacklisted addresses",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.Entries(cmd.Context(), &types.QueryEntriesRequest{Pagination: pageReq})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "entries")

	return cmd
}

// GetEntryCmd returns the command to look up a single blacklisted address.
func GetEntryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "entry [address]",
		Short: "Look up the blacklist entry for an address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Entry(cmd.Context(), &types.QueryEntryRequest{Address: args[0]})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
