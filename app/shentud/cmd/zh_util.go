package cmd

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func PubkeyToAddressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pubkey2addr <pubkey_json>",
		Short: "Convert pubkey to the address",
		Long: `sample: 
shentud pubkey2addr '{"@type": "/cosmos.crypto.ed25519.PubKey", "key": "lzn2Q4Z4DfUdrIDdxVOcXQS44gEdlLpL8T3QnO4brZk="}'`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			var pk cryptotypes.PubKey
			if err := cliCtx.Codec.UnmarshalInterfaceJSON([]byte(args[0]), &pk); err != nil {
				return err
			}
			sdk.AccAddress(pk.Address()).String()
			cliCtx.PrintString(
				sdk.AccAddress(pk.Address()).String() + "\n"+
				sdk.ValAddress(pk.Address()).String() + "\n" +
				sdk.ConsAddress(pk.Address()).String() + "\n",
			)
			return nil
		},
	}
	return cmd
}