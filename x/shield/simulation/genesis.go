package simulation

import (
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	sim "github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/certikfoundation/shentu/x/shield/types"
)

const letters = "abcdefghijklmnopqrstuvwxyz"

// RandomizedGenState creates a random genesis state for module simulation.
func RandomizedGenState(simState *module.SimulationState) {
	r := simState.Rand

	//numPools := uint64(sim.RandIntBetween(r, 10, 30))

	//gs := types.DefaultGenesisState()
	gs := types.GenesisState{}
	//gs.ShieldAdmin =
	gs.NextPoolID = 1
	gs.PoolParams = GenPoolParams(r)
	gs.ClaimProposalParams = GenClaimProposalParams(r)

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(gs)
}

// GenPoolParams returns a randomized PoolParams object.
func GenPoolParams(r *rand.Rand) types.PoolParams {
	protectionPeriod := time.Duration(sim.RandIntBetween(r, 60*60*24, 60*60*24*2)) * time.Second
	minPoolLife := time.Duration(sim.RandIntBetween(r, 60*60*24, 60*60*24*5)) * time.Second
	shieldFeesRate := sdk.NewDecWithPrec(int64(sim.RandIntBetween(r, 0, 50)), 3)
	withdrawalPeriod := time.Duration(sim.RandIntBetween(r, 60*60*24, 60*60*24*3)) * time.Second

	return types.NewPoolParams(protectionPeriod, minPoolLife, withdrawalPeriod, shieldFeesRate)
}

// GenClaimProposalParams returns a randomized ClaimProposalParams object.
func GenClaimProposalParams(r *rand.Rand) types.ClaimProposalParams {
	claimPeriod := time.Duration(sim.RandIntBetween(r, 60*60*24, 60*60*24*2)) * time.Second
	payoutPeriod := time.Duration(sim.RandIntBetween(r, 60*60*24, 60*60*24*2)) * time.Second
	minDeposit := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(int64(sim.RandIntBetween(r, 5e7, 2e8)))))
	depositRate := sdk.NewDecWithPrec(int64(sim.RandIntBetween(r, 0, 100)), 3)
	feesRate := sdk.NewDecWithPrec(int64(sim.RandIntBetween(r, 0, 50)), 3)

	return types.NewClaimProposalParams(claimPeriod, payoutPeriod, minDeposit, depositRate, feesRate)
}

// GetRandDenom generates a random coin denom.
func GetRandDenom(r *rand.Rand) string {
	length := sim.RandIntBetween(r, 3, 8)
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[sim.RandIntBetween(r, 0, len(letters))]
	}
	return string(b)
}
