package simulation

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	stakingSim "github.com/cosmos/cosmos-sdk/x/staking/simulation"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/certikfoundation/shentu/common"
)

// RandomizedGenState generates a random GenesisState for staking.
func RandomizedGenState(simState *module.SimulationState) {
	// params
	var unbondTime time.Duration
	simState.AppParams.GetOrGenerate(
		simState.Cdc, stakingSim.UnbondingTime, &unbondTime, simState.Rand,
		func(r *rand.Rand) { unbondTime = stakingSim.GenUnbondingTime(r) },
	)

	var maxValidators uint16
	simState.AppParams.GetOrGenerate(
		simState.Cdc, stakingSim.MaxValidators, &maxValidators, simState.Rand,
		func(r *rand.Rand) { maxValidators = stakingSim.GenMaxValidators(r) },
	)

	// NOTE: the slashing module need to be defined after the staking module on the
	// NewSimulationManager constructor for this to work.
	simState.UnbondTime = unbondTime

	params := stakingTypes.NewParams(simState.UnbondTime, maxValidators, 7, 3, common.MicroCTKDenom)

	// validators & delegations
	var (
		validators  []stakingTypes.Validator
		delegations []stakingTypes.Delegation
	)

	valAddrs := make([]sdk.ValAddress, simState.NumBonded)
	for i := 0; i < int(simState.NumBonded); i++ {
		valAddr := sdk.ValAddress(simState.Accounts[i].Address)
		valAddrs[i] = valAddr

		maxCommission := sdk.NewDecWithPrec(int64(simulation.RandIntBetween(simState.Rand, 1, 100)), 2)
		commission := stakingTypes.NewCommission(
			simulation.RandomDecAmount(simState.Rand, maxCommission),
			maxCommission,
			simulation.RandomDecAmount(simState.Rand, maxCommission),
		)

		validator := stakingTypes.NewValidator(valAddr, simState.Accounts[i].PubKey, stakingTypes.Description{})
		validator.Tokens = sdk.NewInt(simState.InitialStake)
		validator.DelegatorShares = sdk.NewDec(simState.InitialStake)
		validator.Commission = commission

		delegation := stakingTypes.NewDelegation(simState.Accounts[i].Address, valAddr, sdk.NewDec(simState.InitialStake))
		validators = append(validators, validator)
		delegations = append(delegations, delegation)
	}

	stakingGenesis := stakingTypes.NewGenesisState(params, validators, delegations)

	fmt.Printf("Selected randomly generated staking parameters:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, stakingGenesis.Params))
	simState.GenState[stakingTypes.ModuleName] = simState.Cdc.MustMarshalJSON(stakingGenesis)
}
