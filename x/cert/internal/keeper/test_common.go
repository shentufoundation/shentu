package keeper

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// RandomValidator returns a random validator given access to the keeper and ctx.
func RandomValidator(r *rand.Rand, k Keeper, ctx sdk.Context) (staking.Validator, bool) {
	vals := k.stakingKeeper.GetAllValidators(ctx)
	if len(vals) == 0 {
		return staking.Validator{}, false
	}

	i := r.Intn(len(vals))
	return vals[i], true
}

// RandomDelegator returns a random delegator address given access to the keeper and ctx.
func RandomDelegator(r *rand.Rand, k Keeper, ctx sdk.Context) (sdk.AccAddress, bool) {
	val, ok := RandomValidator(r, k, ctx)
	if !ok {
		return nil, false
	}

	valAddr, err := sdk.ValAddressFromBech32(val.OperatorAddress)
	if err != nil {
		panic(err)
	}
	dels := k.stakingKeeper.GetValidatorDelegations(ctx, valAddr)

	i := r.Intn(len(dels))
	delAddr, err := sdk.AccAddressFromBech32(dels[i].DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return delAddr, true
}
