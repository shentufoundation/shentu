package simulation

import (
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/certikfoundation/shentu/v2/x/shield/types"
)

const letters = "abcdefghijklmnopqrstuvwxyz"

// RandomizedGenState creates a random genesis state for module simulation.
func RandomizedGenState(simState *module.SimulationState) {
	r := simState.Rand

	gs := types.GenesisState{}
	simAccount, _ := simtypes.RandomAcc(r, simState.Accounts)
	gs.ShieldAdmin = simAccount.Address.String()
	gs.NextPoolId = 1
	gs.PoolParams = GenPoolParams(r)
	gs.ClaimProposalParams = GenClaimProposalParams(r)
	gs.BlockRewardParams = GenBlockRewardParams(r)

	var stakingGenState stakingtypes.GenesisState
	stakingGenStatebz := simState.GenState[stakingtypes.ModuleName]
	simState.Cdc.MustUnmarshalJSON(stakingGenStatebz, &stakingGenState)
	gs.PoolParams.WithdrawPeriod = stakingGenState.Params.UnbondingTime

	gs.ClaimProposalParams.ClaimPeriod = time.Duration(simtypes.RandIntBetween(r,
		int(gs.PoolParams.WithdrawPeriod)/10, int(gs.PoolParams.WithdrawPeriod)))
	if gs.PoolParams.ProtectionPeriod >= gs.ClaimProposalParams.ClaimPeriod {
		gs.PoolParams.ProtectionPeriod = time.Duration(simtypes.RandIntBetween(r,
			int(gs.ClaimProposalParams.ClaimPeriod)/10, int(gs.ClaimProposalParams.ClaimPeriod)))
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&gs)
}

// GenPoolParams returns a randomized PoolParams object.
func GenPoolParams(r *rand.Rand) types.PoolParams {
	protectionPeriod := time.Duration(simtypes.RandIntBetween(r, 60*1, 60*60*24*2)) * time.Second
	withdrawPeriod := time.Duration(simtypes.RandIntBetween(r, 60*1, 60*60*24*3)) * time.Second
	shieldFeesRate := sdk.NewDecWithPrec(int64(simtypes.RandIntBetween(r, 0, 50)), 3)
	poolShieldLimit := sdk.NewDecWithPrec(int64(simtypes.RandIntBetween(r, 1, 20)), 2)
	withdrawFeesRate := sdk.NewDecWithPrec(int64(simtypes.RandIntBetween(r, 1, 50)), 3)
	cooldownPeriod := time.Duration(simtypes.RandIntBetween(r, 60*1, 60*60*24*3)) * time.Second
	return types.NewPoolParams(protectionPeriod, withdrawPeriod, cooldownPeriod, shieldFeesRate, withdrawFeesRate, poolShieldLimit, sdk.NewCoins())
}

// GenClaimProposalParams returns a randomized ClaimProposalParams object.
func GenClaimProposalParams(r *rand.Rand) types.ClaimProposalParams {
	claimPeriod := time.Duration(simtypes.RandIntBetween(r, 60*60*24, 60*60*24*2)) * time.Second
	payoutPeriod := time.Duration(simtypes.RandIntBetween(r, 60*60*24, 60*60*24*2)) * time.Second
	minDeposit := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(int64(simtypes.RandIntBetween(r, 5e7, 2e8)))))
	depositRate := sdk.NewDecWithPrec(int64(simtypes.RandIntBetween(r, 0, 100)), 3)
	feesRate := sdk.NewDecWithPrec(int64(simtypes.RandIntBetween(r, 0, 50)), 3)

	return types.NewClaimProposalParams(claimPeriod, payoutPeriod, minDeposit, depositRate, feesRate)
}

// GenShieldStakingRateParam returns a randomized staking-shield rate.
func GenShieldStakingRateParam(r *rand.Rand) sdk.Dec {
	random := simtypes.RandomDecAmount(r, sdk.NewDec(10))
	if random.Equal(sdk.ZeroDec()) {
		return sdk.NewDec(2)
	}
	return random
}

// GenBlockRewardParams returns a randomized BlockRewardParams object.
func GenBlockRewardParams(r *rand.Rand) types.BlockRewardParams {
	modelParamA := sdk.NewDecWithPrec(int64(simtypes.RandIntBetween(r, 0, 20)), 2)
	modelParamB := sdk.NewDecWithPrec(int64(simtypes.RandIntBetween(r, 20, 40)), 2)
	targetLeverage := sdk.NewDecWithPrec(int64(simtypes.RandIntBetween(r, 40, 60)), 1)
	return types.NewBlockRewardParams(modelParamA, modelParamB, targetLeverage)
}

// GetRandDenom generates a random coin denom.
func GetRandDenom(r *rand.Rand) string {
	length := simtypes.RandIntBetween(r, 3, 8)
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[simtypes.RandIntBetween(r, 0, len(letters))]
	}
	return string(b)
}
