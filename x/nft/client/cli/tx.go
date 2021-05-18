package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"

	nftcli "github.com/irisnet/irismod/modules/nft/client/cli"

	"github.com/certikfoundation/shentu/x/nft/types"
)

// NewTxCmd returns the transaction commands for this module
func NewTxCmd() *cobra.Command {
	txCmd := nftcli.NewTxCmd()

	txCmd.AddCommand(
		GetCmdCreateAdmin(),
		GetCmdRevokeAdmin(),
	)

	return txCmd
}

func GetCmdCreateAdmin() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create-admin [address]",
		Long: "Create an NFT administrator account.",
		Example: fmt.Sprintf(
			"$ %s tx nft create-admin certik1... ",
			version.AppName,
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			creator := clientCtx.GetFromAddress().String()
			msg := &types.MsgCreateAdmin{
				Creator: creator,
				Address: args[0],
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdRevokeAdmin() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revoke-admin [address]",
		Long: "Revoke an NFT administrator account.",
		Example: fmt.Sprintf(
			"$ %s tx nft revoke-admin certik1... ",
			version.AppName,
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			creator := clientCtx.GetFromAddress().String()
			msg := &types.MsgRevokeAdmin{
				Revoker: creator,
				Address: args[0],
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
