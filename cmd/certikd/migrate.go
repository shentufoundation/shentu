package main

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	extypes "github.com/cosmos/cosmos-sdk/x/genutil"

	"github.com/certikfoundation/shentu/x/genutil"
	"github.com/certikfoundation/shentu/x/slashing"
)

const (
	flagGenesisTime = "genesis-time"
	flagChainID     = "chain-id"
)

// MigrateGenesisCmd returns a command to execute genesis state migration.
func MigrateGenesisCmd(_ *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate [genesis-file]",
		Short: "Migrate genesis",
		Long: fmt.Sprintf(`Migrate the source genesis into the target version and print to STDOUT.
Example:
$ %s migrate /path/to/genesis.json --chain-id=shentu-incentivized-2 --genesis-time=2020-07-24T17:00:00Z
`, version.ServerName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			importGenesis := args[0]

			genDoc, err := types.GenesisDocFromFile(importGenesis)
			if err != nil {
				return errors.Wrapf(err, "failed to read genesis document from file %s", importGenesis)
			}

			var initialState extypes.AppMap
			if err := cdc.UnmarshalJSON(genDoc.AppState, &initialState); err != nil {
				return errors.Wrap(err, "failed to JSON unmarshal initial genesis state")
			}

			newGenState := genutil.Migrate(initialState)

			// set slashing parameters to zero
			if initialState[slashing.ModuleName] != nil {
				var slashingGenState slashing.GenesisState
				cdc.MustUnmarshalJSON(initialState[slashing.ModuleName], &slashingGenState)
				slashingGenState.Params.SlashFractionDoubleSign = sdk.ZeroDec()
				slashingGenState.Params.SlashFractionDowntime = sdk.ZeroDec()

				// ensure resulting JSON has empty slices, not nil
				for k, v := range slashingGenState.MissedBlocks {
					if v == nil {
						slashingGenState.MissedBlocks[k] = []slashing.MissedBlock{}
					}
				}
				initialState[slashing.ModuleName] = cdc.MustMarshalJSON(slashingGenState)
				genDoc.AppState = cdc.MustMarshalJSON(initialState)
			}

			genDoc.AppState, err = cdc.MarshalJSON(newGenState)
			if err != nil {
				return errors.Wrap(err, "failed to JSON marshal migrated genesis state")
			}

			genesisTime := cmd.Flag(flagGenesisTime).Value.String()
			if genesisTime != "" {
				var t time.Time

				err := t.UnmarshalText([]byte(genesisTime))
				if err != nil {
					return errors.Wrap(err, "failed to unmarshal genesis time")
				}

				genDoc.GenesisTime = t
			}

			chainID := cmd.Flag(flagChainID).Value.String()
			if chainID != "" {
				genDoc.ChainID = chainID
			}

			bz, err := cdc.MarshalJSONIndent(genDoc, "", "  ")
			if err != nil {
				return errors.Wrap(err, "failed to marshal genesis doc")
			}

			sortedBz, err := sdk.SortJSON(bz)
			if err != nil {
				return errors.Wrap(err, "failed to sort JSON genesis doc")
			}

			fmt.Println(string(sortedBz))
			return nil
		},
	}

	cmd.Flags().String(flagGenesisTime, "", "override genesis_time with this flag")
	cmd.Flags().String(flagChainID, "", "override chain_id with this flag")

	return cmd
}
