package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"

	shieldtypes "github.com/shentufoundation/shentu/v2/x/shield/types"
)

// AddGenesisShieldAdminCmd returns add-genesis-shield-admin cobra Command.
func AddGenesisShieldAdminCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-shield-admin [address]",
		Short: "Add a genesis shield admin to genesis.json",
		Long:  `Add a genesis shield admin to genesis.json. The provided shield admin must specify the account address. `,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := client.GetClientContextFromCmd(cmd)
			depCdc := ctx.Codec
			cdc := depCdc

			config := server.GetServerContextFromCmd(cmd).Config
			config.SetRoot(ctx.HomeDir)

			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}
			shieldGenState := shieldtypes.GetGenesisStateFromAppState(cdc, appState)

			shieldGenStateBz := cdc.MustMarshalJSON(&shieldGenState)

			appState[shieldtypes.ModuleName] = shieldGenStateBz
			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}
			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	return cmd
}
