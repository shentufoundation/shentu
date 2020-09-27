package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	shieldQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the shield module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	shieldQueryCmd.AddCommand(flags.GetCommands(
		GetCmdPool(queryRoute, cdc),
	)...)

	return shieldQueryCmd
}

// GetCmdPool returns the command for querying the pool.
func GetCmdPool(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool [pool-id]",
		Short: "Get pool information",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			var res []byte
			var err error
			if len(args) == 1 {
				res, _, err = cliCtx.QueryWithData(fmt.Sprintf("custom/%s/pool/id/%s", queryRoute, args[0]), nil)
				if err != nil {
					return err
				}
			} else {
				sponsor := viper.GetString(flagSponsor)
				if sponsor == "" {
					return fmt.Errorf("either poolID or sponsor is required to query pool")
				}
				res, _, err = cliCtx.QueryWithData(fmt.Sprintf("custom/%s/pool/sponsor/%s", queryRoute, sponsor), nil)
				if err != nil {
					return err
				}
			}
			var out types.Pool
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
	cmd.Flags().String(flagSponsor, "", "use sponsor to query the pool info")

	return cmd
}
