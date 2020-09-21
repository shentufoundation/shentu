package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/CertiKProject/certik-chain/x/shield/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	shieldTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Shield transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	shieldTxCmd.AddCommand(flags.PostCommands()...)

	return shieldTxCmd
}
