package simulation

import (
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	sim "github.com/cosmos/cosmos-sdk/x/simulation"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/certikfoundation/shentu/x/shield/types"
	"github.com/certikfoundation/shentu/x/staking"
)

const letters = "abcdefghijklmnopqrstuvwxyz"

// RandomizedGenState creates a random genesis state for module simulation.
func RandomizedGenState(simState *module.SimulationState) {
	r := simState.Rand

	gs := types.GenesisState{}
	simAccount, _ := sim.RandomAcc(r, simState.Accounts)
	gs.ShieldAdmin = simAccount.Address
	gs.NextPoolID = 1
	gs.PoolParams = GenPoolParams(r)
	gs.ClaimProposalParams = GenClaimProposalParams(r)

	var stakingGenState stakingTypes.GenesisState
	stakingGenStatebz := simState.GenState[staking.ModuleName]
	stakingTypes.ModuleCdc.MustUnmarshalJSON(stakingGenStatebz, &stakingGenState)
	gs.PoolParams.WithdrawPeriod = stakingGenState.Params.UnbondingTime

	gs.ClaimProposalParams.ClaimPeriod = time.Duration(sim.RandIntBetween(r,
		int(gs.PoolParams.WithdrawPeriod)/10, int(gs.PoolParams.WithdrawPeriod)))
	if gs.PoolParams.ProtectionPeriod >= gs.ClaimProposalParams.ClaimPeriod {
		gs.PoolParams.ProtectionPeriod = time.Duration(sim.RandIntBetween(r,
			int(gs.ClaimProposalParams.ClaimPeriod)/10, int(gs.ClaimProposalParams.ClaimPeriod)))
	}
	gs.ShieldStakingRate = GenShieldStakingRateParam(r)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(gs)
}

// GenPoolParams returns a randomized PoolParams object.
func GenPoolParams(r *rand.Rand) types.PoolParams {
	protectionPeriod := time.Duration(sim.RandIntBetween(r, 60*1, 60*60*24*2)) * time.Second
	withdrawPeriod := time.Duration(sim.RandIntBetween(r, 60*1, 60*60*24*3)) * time.Second
	shieldFeesRate := sdk.NewDecWithPrec(int64(sim.RandIntBetween(r, 0, 50)), 3)
	poolShieldLimit := sdk.NewDecWithPrec(int64(sim.RandIntBetween(r, 1, 20)), 2)

	return types.NewPoolParams(protectionPeriod, withdrawPeriod, shieldFeesRate, poolShieldLimit, sdk.Coins{})
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

// GenShieldStakingRateParam returns a randomized staking-shield rate.
func GenShieldStakingRateParam(r *rand.Rand) sdk.Dec {
	random := sim.RandomDecAmount(r, sdk.NewDec(10))
	if random.Equal(sdk.ZeroDec()) {
		return sdk.NewDec(2)
	}
	return random
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
