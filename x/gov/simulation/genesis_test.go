package simulation_test

import (
	"encoding/json"
	"math/rand"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/shentufoundation/shentu/v2/x/gov/simulation"
	typesv1 "github.com/shentufoundation/shentu/v2/x/gov/types/v1"
)

// TestRandomizedGenState exercises the happy path of RandomizedGenState:
// valid genesis bytes, populated Params, and — the shentu-specific
// invariant — a CertifierUpdateSecurityVoteTally derived from the same
// randomized quorum/threshold/veto as Params. If this test ever has
// to assert a non-nil stake tally, the v7 single-round model has
// regressed.
func TestRandomizedGenState(t *testing.T) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	r := rand.New(rand.NewSource(1))
	simState := module.SimulationState{
		AppParams:    make(simtypes.AppParams),
		Cdc:          cdc,
		Rand:         r,
		NumBonded:    3,
		BondDenom:    sdk.DefaultBondDenom,
		Accounts:     simtypes.RandomAccounts(r, 3),
		InitialStake: sdkmath.NewInt(1000),
		GenState:     make(map[string]json.RawMessage),
	}

	simulation.RandomizedGenState(&simState)

	raw, ok := simState.GenState[govtypes.ModuleName]
	require.True(t, ok, "gov genesis must be written into GenState")

	var govGenesis typesv1.GenesisState
	require.NoError(t, cdc.UnmarshalJSON(raw, &govGenesis))

	require.NotNil(t, govGenesis.Params, "Params must be populated")
	require.NotNil(t, govGenesis.CustomParams, "CustomParams must be populated")
	require.NotNil(t, govGenesis.CustomParams.CertifierUpdateSecurityVoteTally,
		"security tally must be set — it drives the certifier round")

	// The security tally inherits quorum/threshold/veto from Params —
	// they share the same randomized dec values by construction.
	secTally := govGenesis.CustomParams.CertifierUpdateSecurityVoteTally
	require.Equal(t, govGenesis.Params.Quorum, secTally.Quorum)
	require.Equal(t, govGenesis.Params.Threshold, secTally.Threshold)
	require.Equal(t, govGenesis.Params.VetoThreshold, secTally.VetoThreshold)

	require.NoError(t, typesv1.ValidateGenesis(&govGenesis),
		"randomized genesis must satisfy ValidateGenesis")
}

// TestRandomizedGenState_MissingFieldsPanic confirms RandomizedGenState
// treats an under-constructed SimulationState as a programmer error
// (matches upstream gov's contract). The simulator never calls us with
// an empty simState — if it ever does, panicking is the right failure
// mode so the regression is caught immediately.
func TestRandomizedGenState_MissingFieldsPanic(t *testing.T) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)
	r := rand.New(rand.NewSource(1))

	cases := []struct {
		name     string
		simState module.SimulationState
	}{
		{"fully empty", module.SimulationState{}},
		{"no GenState map", module.SimulationState{
			AppParams: make(simtypes.AppParams),
			Cdc:       cdc,
			Rand:      r,
		}},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Panics(t, func() { simulation.RandomizedGenState(&tc.simState) })
		})
	}
}
