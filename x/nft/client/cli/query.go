package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	nfttypes "github.com/irisnet/irismod/modules/nft/types"

	"github.com/certikfoundation/shentu/x/nft/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                nfttypes.ModuleName,
		Short:              "Querying commands for the NFT module",
		DisableFlagParsing: true,
	}

	queryCmd.AddCommand(
		GetCmdQueryAdmin(),
		GetCmdQueryAdmins(),
	)

	return queryCmd
}

func GetCmdQueryAdmin() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "admin [address]",
		Long:    "Query an address to see if it's an administrator for the NFT module.",
		Example: fmt.Sprintf("$ %s query nft admin <address>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Admin(context.Background(), &types.QueryAdminRequest{
				Address: args[0],
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryAdmins() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "admins",
		Long:    "Query all administrators for the NFT module.",
		Example: fmt.Sprintf("$ %s query nft admins", version.AppName),
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Admins(context.Background(), &types.QueryAdminsRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
