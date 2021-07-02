package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
		GetCmdQueryCertificate(),
		GetCmdQueryCertificates(),
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

func GetCmdQueryCertificate() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "certificate [denom-id] [token-id]",
		Long:    "Query a certificate by the specific denom-id and token-id.",
		Example: fmt.Sprintf("$ %s query nft certificate <denom-id> <token-id>", version.AppName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if err := nfttypes.ValidateDenomID(args[0]); err != nil {
				return err
			}
			if err := nfttypes.ValidateTokenID(args[1]); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Certificate(context.Background(), &types.QueryCertificateRequest{
				DenomId: args[0],
				TokenId: args[1],
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

func GetCmdQueryCertificates() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "certificates [denom-id] [<flags>]",
		Long:    "Query all certificates of given denom.",
		Example: fmt.Sprintf("$ %s query nft certificates <denom-id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if err := nfttypes.ValidateDenomID(args[0]); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			resp, err := queryClient.Certificates(context.Background(),
				&types.QueryCertificatesRequest{
					Certifier:  viper.GetString(FlagCertifier),
					DenomId:    args[0],
					Pagination: pageReq,
				})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	cmd.Flags().String(FlagCertifier, "", "certificates issued by certifier")
	flags.AddPaginationFlagsToCmd(cmd, "certificates")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
