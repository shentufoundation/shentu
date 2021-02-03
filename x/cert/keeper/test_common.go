package keeper

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
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

	dels := k.stakingKeeper.GetValidatorDelegations(ctx, val.OperatorAddress)

	i := r.Intn(len(dels))
	return dels[i].DelegatorAddress, true
}
