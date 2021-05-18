package cli

import (
	"github.com/spf13/cobra"

	nftcli "github.com/irisnet/irismod/modules/nft/client/cli"
)

// NewTxCmd returns the transaction commands for this module
func NewTxCmd() *cobra.Command {
	txCmd := nftcli.NewTxCmd()

	txCmd.AddCommand(
	)

	return txCmd
}

func GetCmdCreateAdmin