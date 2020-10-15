package keeper

import (
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

type Keeper struct {
	staking.Keeper
	storeKey sdk.StoreKey
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, supplyKeeper types.SupplyKeeper, paramstore params.Subspace) Keeper {
	return Keeper{
		Keeper:   staking.NewKeeper(cdc, key, supplyKeeper, paramstore),
		storeKey: key,
	}
}

func (k Keeper) RemoveUBDQueue(ctx sdk.Context, timestamp time.Time) {
	ctx.KVStore(k.storeKey).Delete(staking.GetUnbondingDelegationTimeKey(timestamp))
}

// UpdateValidatorCommission attempts to update a validator's commission rate.
// An error is returned if the new commission rate is invalid.
func (k Keeper) UpdateValidatorCommission(ctx sdk.Context,
	validator types.Validator, newRate sdk.Dec) (types.Commission, error) {
	if newRate.LT(sdk.OneDec()) {
		return types.Commission{}, sdkerrors.Register(types.ModuleName, 100, "commission cannot be less than 100%")
	}

	commission := validator.Commission
	blockTime := ctx.BlockHeader().Time

	if err := commission.ValidateNewRate(newRate, blockTime); err != nil {
		return commission, err
	}

	commission.Rate = newRate
	commission.UpdateTime = blockTime

	return commission, nil
}
