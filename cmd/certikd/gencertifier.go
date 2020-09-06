package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"

	"github.com/certikfoundation/shentu/x/cert"
)

// AddGenesisCertifierCmd returns add-genesis-certifier cobra Command.
func AddGenesisCertifierCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-certifier [address]",
		Short: "Add a genesis certifier to genesis.json",
		Long:  `Add a genesis certifier to genesis.json. The provided certifier must specify the account address. `,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("failed to parse address")
			}

			genFile := config.GenesisFile()
			appState, genDoc, err := genutil.GenesisStateFromGenFile(cdc, genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}
			certGenState := cert.GetGenesisStateFromAppState(cdc, appState)

			certGenState.Certifiers = append(certGenState.Certifiers,
				cert.NewCertifier(addr, "", nil, ""))

			certGenStateBz, err := cdc.MarshalJSON(certGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
			}
			appState[cert.ModuleName] = certGenStateBz
			appStateJSON, err := cdc.MarshalJSON(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}
			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}
	return cmd
}
