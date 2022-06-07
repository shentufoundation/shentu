package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	cryptocodec "github.com/tendermint/tendermint/crypto/encoding"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmtypes "github.com/tendermint/tendermint/types"
	"io/ioutil"
	"log"
	"sort"

	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	slashing "github.com/cosmos/cosmos-sdk/x/slashing/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"
)

const (
	flagReplacementKeys = "replacement-cons-keys"
)

// RotateValKeysCmd returns a command to execute genesis state migration.
func RotateValKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "replace-validators [genesis-file] [replacement-cons-keys]",
		Short: "Replace validators in a given genesis with a set json",
		Long: fmt.Sprintf(`Migrate the source genesis into the target version and print to STDOUT.
Example:
$ %s migrate /path/to/genesis.json --chain-id=cosmoshub-4 --genesis-time=2019-04-22T17:00:00Z --initial-height=5000
`, version.AppName),
		Args: cobra.ExactArgs(2),
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

			replacementKeys := args[1]
			genDoc = loadKeydataFromFile(clientCtx, replacementKeys, genDoc)

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

	cmd.Flags().String(flagReplacementKeys, "", "Proviide a JSON file to replace the consensus keys of validators")

	return cmd
}

type replacementKeys []interface{}

func loadKeydataFromFile(clientCtx client.Context, replacementsJSON string, genDoc *tmtypes.GenesisDoc) *tmtypes.GenesisDoc {
	jsonReplacementBlob, err := ioutil.ReadFile(replacementsJSON)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "failed to read replacement keys from file %s", replacementsJSON))
	}

	var rks replacementKeys
	if err = json.Unmarshal(jsonReplacementBlob, &rks); err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	var state types.AppMap
	if err := json.Unmarshal(genDoc.AppState, &state); err != nil {
		log.Fatal(errors.Wrap(err, "failed to JSON unmarshal initial genesis state"))
	}

	var stakingGenesis staking.GenesisState
	var slashingGenesis slashing.GenesisState

	clientCtx.JSONCodec.MustUnmarshalJSON(state[staking.ModuleName], &stakingGenesis)
	clientCtx.JSONCodec.MustUnmarshalJSON(state[slashing.ModuleName], &slashingGenesis)

	// sort validators power descending
	sort.Slice(stakingGenesis.Validators, func(i, j int) bool {
		return stakingGenesis.Validators[i].BondedTokens().GT(stakingGenesis.Validators[j].GetTokens())
	})

	for i, _ := range rks {
		val := stakingGenesis.Validators[i]
		toReplaceValConsAddress, err := val.GetConsAddr()
		if err != nil {
			panic(err)
		}

		bz, err := json.Marshal(rks[i])
		if err != nil {
			panic(err)
		}
		var tmp codectypes.Any
		clientCtx.JSONCodec.MustUnmarshalJSON(bz, &tmp)
		var pk *codectypes.Any
		pk, err = codectypes.NewAnyWithValue(&tmp)
		if err != nil {
			panic(err)
		}

		stakingGenesis.Validators[i].ConsensusPubkey = pk
		if err != nil {
			panic(err)
		}

		replaceValConsAddress, err := val.GetConsAddr()
		if err != nil {
			panic(err)
		}
		protoReplaceValConsPubKey, err := val.TmConsPublicKey()
		if err != nil {
			panic(err)
		}
		replaceValConsPubKey, err := cryptocodec.PubKeyFromProto(protoReplaceValConsPubKey)
		if err != nil {
			panic(err)
		}

		for j, signingInfo := range slashingGenesis.SigningInfos {
			if signingInfo.Address == toReplaceValConsAddress.String() {
				slashingGenesis.SigningInfos[j].Address = replaceValConsAddress.String()
				slashingGenesis.SigningInfos[j].ValidatorSigningInfo.Address = replaceValConsAddress.String()
			}
		}

		for k, missedInfo := range slashingGenesis.MissedBlocks {
			if missedInfo.Address == toReplaceValConsAddress.String() {
				slashingGenesis.MissedBlocks[k].Address = replaceValConsAddress.String()
			}
		}

		for tmIdx, tmval := range genDoc.Validators {
			if bytes.Equal(tmval.Address.Bytes(), replaceValConsAddress.Bytes()) {
				fmt.Println("wow")
				genDoc.Validators[tmIdx].Address = replaceValConsAddress.Bytes()
				genDoc.Validators[tmIdx].PubKey = replaceValConsPubKey
			}
		}
		stakingGenesis.Validators[i].ConsensusPubkey
	}
	state[staking.ModuleName] = clientCtx.JSONCodec.MustMarshalJSON(&stakingGenesis)
	state[slashing.ModuleName] = clientCtx.JSONCodec.MustMarshalJSON(&slashingGenesis)

	genDoc.AppState, err = json.Marshal(state)

	if err != nil {
		log.Fatal("Could not marshal App State")
	}
	return genDoc
}
