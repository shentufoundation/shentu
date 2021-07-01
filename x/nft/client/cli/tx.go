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

const (
	FlagCertifier   = "certifier"
	FlagContent     = "content"
	FlagDescription = "description"
	FlagName        = "name"
	FlagURI         = "uri"
)

// NewTxCmd returns the transaction commands for this module
func NewTxCmd() *cobra.Command {
	txCmd := nftcli.NewTxCmd()

	txCmd.AddCommand(
		GetCmdCreateAdmin(),
		GetCmdRevokeAdmin(),
		GetCmdIssueCertificate(),
		GetCmdEditCertificate(),
		GetCmdRevokeCertificate(),
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

func GetCmdIssueCertificate() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "issue-certificate [denom-id] [token-id] [<flags>]",
		Long: "Issue a certificate NFT.",
		Example: fmt.Sprintf(
			"$ %s tx nft issue-certificate <denom-id> <token-id> <content> "+
				"--uri=<uri> "+
				"--name=<name> "+
				"--content=<content> "+
				"--description=<description>",
			version.AppName,
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			from := cliCtx.GetFromAddress()
			if err := txf.AccountRetriever().EnsureExists(cliCtx, from); err != nil {
				return err
			}

			name, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return err
			}
			uri, err := cmd.Flags().GetString(FlagURI)
			if err != nil {
				return err
			}
			content, err := cmd.Flags().GetString(FlagContent)
			if err != nil {
				return err
			}
			description, err := cmd.Flags().GetString(FlagDescription)
			if err != nil {
				return err
			}

			msg := types.NewMsgIssueCertificate(args[0], args[1], name, uri, content, description, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	cmd.Flags().String(FlagName, "", "name")
	cmd.Flags().String(FlagURI, "", "uri")
	cmd.Flags().String(FlagContent, "", "content")
	cmd.Flags().String(FlagDescription, "", "description")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdEditCertificate() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "edit-certificate [denom-id] [token-id] [<flags>]",
		Long: "Edit an existing certificate NFT.",
		Example: fmt.Sprintf(
			"$ %s tx nft edit-certificate <denom-id> <token-id> "+
				"--uri=<uri> "+
				"--name=<name> "+
				"--content=<content> "+
				"--description=<description>",
			version.AppName,
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			from := cliCtx.GetFromAddress()
			if err := txf.AccountRetriever().EnsureExists(cliCtx, from); err != nil {
				return err
			}

			name, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return err
			}
			uri, err := cmd.Flags().GetString(FlagURI)
			if err != nil {
				return err
			}
			content, err := cmd.Flags().GetString(FlagContent)
			if err != nil {
				return err
			}
			description, err := cmd.Flags().GetString(FlagDescription)
			if err != nil {
				return err
			}

			msg := types.NewMsgEditCertificate(args[0], args[1], name, uri, content, description, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	cmd.Flags().String(FlagName, "", "name")
	cmd.Flags().String(FlagURI, "", "uri")
	cmd.Flags().String(FlagContent, "", "content")
	cmd.Flags().String(FlagDescription, "", "description")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdRevokeCertificate() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revoke-certificate [denom-id] [token-id] [<description>]",
		Long: "Revoke a certificate.",
		Example: fmt.Sprintf(
			"$ %s tx nft revoke-certificate <denom-id> <token-id>",
			version.AppName,
		),
		Args: cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig).WithAccountRetriever(cliCtx.AccountRetriever)

			from := cliCtx.GetFromAddress()
			if err := txf.AccountRetriever().EnsureExists(cliCtx, from); err != nil {
				return err
			}

			var description string
			if len(args) > 2 {
				description = args[2]
			}

			msg := types.NewMsgRevokeCertificate(args[0], args[1], description, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(cliCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
