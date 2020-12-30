package cmd

import (
	"encoding/json"
	"fmt"

	shieldtypes "github.com/certikfoundation/shentu/x/shield/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/spf13/cobra"
)

// AddGenesisShieldAdminCmd returns add-genesis-shield-admin cobra Command.
func AddGenesisShieldAdminCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-shield-admin [address]",
		Short: "Add a genesis shield admin to genesis.json",
		Long:  `Add a genesis shield admin to genesis.json. The provided shield admin must specify the account address. `,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.GetClientContextFromCmd(cmd)
			depCdc := ctx.JSONMarshaler
			cdc := depCdc.(codec.Marshaler)

			config := server.GetServerContextFromCmd(cmd).Config
			config.SetRoot(ctx.HomeDir)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("failed to parse address")
			}

			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}
			shieldGenState := shieldtypes.GetGenesisStateFromAppState(cdc, appState)

			shieldGenState.ShieldAdmin = addr.String()

			shieldGenStateBz, err := json.Marshal(shieldGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal shield genesis state: %w", err)
			}
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
