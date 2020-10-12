package keeper

import (
	"time"

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
