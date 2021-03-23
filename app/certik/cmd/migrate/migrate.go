package migrate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	oracletypes "github.com/certikfoundation/shentu/x/oracle/types"

	shieldtypes "github.com/certikfoundation/shentu/x/shield/types"

	"github.com/spf13/cobra"

	tmjson "github.com/tendermint/tendermint/libs/json"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"
	v039auth "github.com/cosmos/cosmos-sdk/x/auth/legacy/v039"
	v040auth "github.com/cosmos/cosmos-sdk/x/auth/legacy/v040"
	v036supply "github.com/cosmos/cosmos-sdk/x/bank/legacy/v036"
	v038bank "github.com/cosmos/cosmos-sdk/x/bank/legacy/v038"
	v040bank "github.com/cosmos/cosmos-sdk/x/bank/legacy/v040"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	captypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	v039crisis "github.com/cosmos/cosmos-sdk/x/crisis/legacy/v039"
	v040crisis "github.com/cosmos/cosmos-sdk/x/crisis/legacy/v040"
	v036distr "github.com/cosmos/cosmos-sdk/x/distribution/legacy/v036"
	v038distr "github.com/cosmos/cosmos-sdk/x/distribution/legacy/v038"
	v040distr "github.com/cosmos/cosmos-sdk/x/distribution/legacy/v040"
	v038evidence "github.com/cosmos/cosmos-sdk/x/evidence/legacy/v038"
	v040evidence "github.com/cosmos/cosmos-sdk/x/evidence/legacy/v040"
	evtypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	v039genutil "github.com/cosmos/cosmos-sdk/x/genutil/legacy/v039"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	v036gov "github.com/cosmos/cosmos-sdk/x/gov/legacy/v036"
	v040gov "github.com/cosmos/cosmos-sdk/x/gov/legacy/v040"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ibcxfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	host "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
	"github.com/cosmos/cosmos-sdk/x/ibc/core/exported"
	ibccoretypes "github.com/cosmos/cosmos-sdk/x/ibc/core/types"
	v039mint "github.com/cosmos/cosmos-sdk/x/mint/legacy/v039"
	v040mint "github.com/cosmos/cosmos-sdk/x/mint/legacy/v040"
	v036params "github.com/cosmos/cosmos-sdk/x/params/legacy/v036"
	v039slashing "github.com/cosmos/cosmos-sdk/x/slashing/legacy/v039"
	v040slashing "github.com/cosmos/cosmos-sdk/x/slashing/legacy/v040"
	v038staking "github.com/cosmos/cosmos-sdk/x/staking/legacy/v038"
	v040staking "github.com/cosmos/cosmos-sdk/x/staking/legacy/v040"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	v038upgrade "github.com/cosmos/cosmos-sdk/x/upgrade/legacy/v038"

	certtypes "github.com/certikfoundation/shentu/x/cert/types"
	cvmtypes "github.com/certikfoundation/shentu/x/cvm/types"
)

const (
	flagGenesisTime     = "genesis-time"
	flagInitialHeight   = "initial-height"
	flagReplacementKeys = "replacement-cons-keys"
	flagNoProp29        = "no-prop-29"
)

func migrateGenutil(oldGenState v039genutil.GenesisState) *types.GenesisState {
	return &types.GenesisState{
		GenTxs: oldGenState.GenTxs,
	}
}

func Migrate(appState types.AppMap, clientCtx client.Context) types.AppMap {
	v039Codec := codec.NewLegacyAmino()
	RegisterAuthLegacyAminoCodec(v039Codec)
	v036gov.RegisterLegacyAminoCodec(v039Codec)
	v036distr.RegisterLegacyAminoCodec(v039Codec)
	v036params.RegisterLegacyAminoCodec(v039Codec)
	v038upgrade.RegisterLegacyAminoCodec(v039Codec)
	RegisterCVMLegacyAminoCodec(v039Codec)
	RegisterCertLegacyAminoCodec(v039Codec)
	RegisterShieldLegacyAminoCodec(v039Codec)

	v040Codec := clientCtx.JSONMarshaler

	if appState[v038bank.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var bankGenState v038bank.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v038bank.ModuleName], &bankGenState)

		// unmarshal x/auth genesis state to retrieve all account balances
		var authGenState v039auth.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v039auth.ModuleName], &authGenState)

		// unmarshal x/supply genesis state to retrieve total supply
		var supplyGenState v036supply.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v036supply.ModuleName], &supplyGenState)

		// delete deprecated x/bank genesis state
		delete(appState, v038bank.ModuleName)

		// delete deprecated x/supply genesis state
		delete(appState, v036supply.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v040bank.ModuleName] = v040Codec.MustMarshalJSON(v040bank.Migrate(bankGenState, authGenState, supplyGenState))
	}

	// remove balances from existing accounts
	if appState[v039auth.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var authGenState v039auth.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v039auth.ModuleName], &authGenState)

		// delete deprecated x/auth genesis state
		delete(appState, v039auth.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v040auth.ModuleName] = v040Codec.MustMarshalJSON(authMigrate(authGenState))
	}

	// Migrate x/crisis.
	if appState[v039crisis.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var crisisGenState v039crisis.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v039crisis.ModuleName], &crisisGenState)

		// delete deprecated x/crisis genesis state
		delete(appState, v039crisis.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v040crisis.ModuleName] = v040Codec.MustMarshalJSON(v040crisis.Migrate(crisisGenState))
	}

	// Migrate x/distribution.
	if appState[v038distr.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var distributionGenState v038distr.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v038distr.ModuleName], &distributionGenState)

		// delete deprecated x/distribution genesis state
		delete(appState, v038distr.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v040distr.ModuleName] = v040Codec.MustMarshalJSON(v040distr.Migrate(distributionGenState))
	}

	// Migrate x/evidence.
	if appState[v038evidence.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var evidenceGenState v038evidence.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v038bank.ModuleName], &evidenceGenState)

		// delete deprecated x/evidence genesis state
		delete(appState, v038evidence.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v040evidence.ModuleName] = v040Codec.MustMarshalJSON(v040evidence.Migrate(evidenceGenState))
	}

	// Migrate x/gov.
	if appState[govtypes.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var govGenState v131GovGenesisState
		v039Codec.MustUnmarshalJSON(appState[govtypes.ModuleName], &govGenState)

		// delete deprecated x/gov genesis state
		delete(appState, v036gov.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v040gov.ModuleName] = v040Codec.MustMarshalJSON(govMigrate(govGenState))
	}

	// Migrate x/mint.
	if appState[v039mint.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var mintGenState v039mint.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v039mint.ModuleName], &mintGenState)

		// delete deprecated x/mint genesis state
		delete(appState, v039mint.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v040mint.ModuleName] = v040Codec.MustMarshalJSON(v040mint.Migrate(mintGenState))
	}

	// Migrate x/slashing.
	if appState[v039slashing.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var slashingGenState v039slashing.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v039slashing.ModuleName], &slashingGenState)

		// delete deprecated x/slashing genesis state
		delete(appState, v039slashing.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v040slashing.ModuleName] = v040Codec.MustMarshalJSON(v040slashing.Migrate(slashingGenState))
	}

	// Migrate x/staking.
	if appState[v038staking.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var stakingGenState v038staking.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v038staking.ModuleName], &stakingGenState)

		// delete deprecated x/staking genesis state
		delete(appState, v038staking.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v040staking.ModuleName] = v040Codec.MustMarshalJSON(v040staking.Migrate(stakingGenState))
	}

	// Migrate x/genutil
	if appState[v039genutil.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var genutilGenState v039genutil.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v039genutil.ModuleName], &genutilGenState)

		// delete deprecated x/staking genesis state
		delete(appState, v039genutil.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[genutiltypes.ModuleName] = v040Codec.MustMarshalJSON(migrateGenutil(genutilGenState))
	}

	// Migrate x/cvm
	if appState[cvmtypes.ModuleName] != nil {
		var cvmGenState CVMGenesisState
		v039Codec.MustUnmarshalJSON(appState[cvmtypes.ModuleName], &cvmGenState)

		delete(appState, cvmtypes.ModuleName)

		appState[cvmtypes.ModuleName] = v040Codec.MustMarshalJSON(migrateCVM(cvmGenState))
	}

	// Migrate x/cert
	if appState[certtypes.ModuleName] != nil {
		var certGenState CertGenesisState
		v039Codec.MustUnmarshalJSON(appState[certtypes.ModuleName], &certGenState)

		delete(appState, certtypes.ModuleName)

		appState[certtypes.ModuleName] = v040Codec.MustMarshalJSON(migrateCert(certGenState))
	}

	// Migrate x/oracle
	if appState[oracletypes.ModuleName] != nil {
		var oracleGenState OracleGenesisState
		v039Codec.MustUnmarshalJSON(appState[oracletypes.ModuleName], &oracleGenState)

		delete(appState, oracletypes.ModuleName)

		appState[oracletypes.ModuleName] = v040Codec.MustMarshalJSON(migrateOracle(oracleGenState))
	}

	// Migrate x/shield
	if appState[shieldtypes.ModuleName] != nil {
		var shieldGenState ShieldGenesisState
		v039Codec.MustUnmarshalJSON(appState[shieldtypes.ModuleName], &shieldGenState)

		delete(appState, shieldtypes.ModuleName)

		appState[shieldtypes.ModuleName] = v040Codec.MustMarshalJSON(migrateShield(shieldGenState))
	}

	return appState
}

// MigrateGenesisCmd returns a command to execute genesis state migration.
func MigrateGenesisCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate [genesis-file]",
		Short: "Migrate 1.X.X genesis to 2.0.X",
		Long: fmt.Sprintf(`Migrate the source genesis into the target version and print to STDOUT.

Example:
$ %s migrate /path/to/genesis.json --chain-id=cosmoshub-4 --genesis-time=2019-04-22T17:00:00Z --initial-height=5000
`, version.AppName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			var err error

			importGenesis := args[0]

			jsonBlob, err := ioutil.ReadFile(importGenesis)

			if err != nil {
				return errors.Wrap(err, "failed to read provided genesis file")
			}

			genDoc, err := tmtypes.GenesisDocFromJSON(jsonBlob)
			if err != nil {
				return errors.Wrapf(err, "failed to read genesis document from file %s", importGenesis)
			}

			var initialState types.AppMap
			if err := json.Unmarshal(genDoc.AppState, &initialState); err != nil {
				return errors.Wrap(err, "failed to JSON unmarshal initial genesis state")
			}

			migrationFunc := Migrate

			// TODO: handler error from migrationFunc call
			newGenState := migrationFunc(initialState, clientCtx)

			var bankGenesis bank.GenesisState

			clientCtx.JSONMarshaler.MustUnmarshalJSON(newGenState[bank.ModuleName], &bankGenesis)

			bankGenesis.DenomMetadata = []bank.Metadata{
				{
					Description: "The native staking token of the CertiK Chain.",
					DenomUnits: []*bank.DenomUnit{
						{Denom: "uctk", Exponent: uint32(0), Aliases: []string{"microctk"}},
						{Denom: "mctk", Exponent: uint32(3), Aliases: []string{"millictk"}},
						{Denom: "ctk", Exponent: uint32(6), Aliases: []string{}},
					},
					Base:    "uctk",
					Display: "ctk",
				},
			}
			newGenState[bank.ModuleName] = clientCtx.JSONMarshaler.MustMarshalJSON(&bankGenesis)

			var stakingGenesis staking.GenesisState

			clientCtx.JSONMarshaler.MustUnmarshalJSON(newGenState[staking.ModuleName], &stakingGenesis)

			ibcTransferGenesis := ibcxfertypes.DefaultGenesisState()
			ibcCoreGenesis := ibccoretypes.DefaultGenesisState()
			capGenesis := captypes.DefaultGenesis()
			evGenesis := evtypes.DefaultGenesisState()

			ibcCoreGenesis.ClientGenesis.Params.AllowedClients = []string{exported.Tendermint}
			stakingGenesis.Params.HistoricalEntries = 10000

			newGenState[ibcxfertypes.ModuleName] = clientCtx.JSONMarshaler.MustMarshalJSON(ibcTransferGenesis)
			newGenState[host.ModuleName] = clientCtx.JSONMarshaler.MustMarshalJSON(ibcCoreGenesis)
			newGenState[captypes.ModuleName] = clientCtx.JSONMarshaler.MustMarshalJSON(capGenesis)
			newGenState[evtypes.ModuleName] = clientCtx.JSONMarshaler.MustMarshalJSON(evGenesis)
			newGenState[staking.ModuleName] = clientCtx.JSONMarshaler.MustMarshalJSON(&stakingGenesis)

			genDoc.AppState, err = json.Marshal(newGenState)
			if err != nil {
				return errors.Wrap(err, "failed to JSON marshal migrated genesis state")
			}

			genesisTime, _ := cmd.Flags().GetString(flagGenesisTime)
			if genesisTime != "" {
				var t time.Time

				err := t.UnmarshalText([]byte(genesisTime))
				if err != nil {
					return errors.Wrap(err, "failed to unmarshal genesis time")
				}

				genDoc.GenesisTime = t
			}

			chainID, _ := cmd.Flags().GetString(flags.FlagChainID)
			if chainID != "" {
				genDoc.ChainID = chainID
			}

			initialHeight, _ := cmd.Flags().GetInt(flagInitialHeight)

			genDoc.InitialHeight = int64(initialHeight)

			replacementKeys, _ := cmd.Flags().GetString(flagReplacementKeys)

			if replacementKeys != "" {
				genDoc = loadKeydataFromFile(clientCtx, replacementKeys, genDoc)
			}

			bz, err := tmjson.Marshal(genDoc)
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
	cmd.Flags().Int(flagInitialHeight, 0, "Set the starting height for the chain")
	cmd.Flags().String(flagReplacementKeys, "", "Proviide a JSON file to replace the consensus keys of validators")
	cmd.Flags().String(flags.FlagChainID, "", "override chain_id with this flag")
	cmd.Flags().Bool(flagNoProp29, false, "Do not implement fund recovery from prop29")

	return cmd
}
