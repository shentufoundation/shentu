package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

type Keeper struct {
	keeper.Keeper
	storeKey storetypes.StoreKey
}

func NewKeeper(cdc codec.BinaryCodec, key storetypes.StoreKey, ak types.AccountKeeper, bk types.BankKeeper, paramstore paramtypes.Subspace) Keeper {
	return Keeper{
		Keeper:   keeper.NewKeeper(cdc, key, ak, bk, paramstore),
		storeKey: key,
	}
}

func (k Keeper) RemoveUBDQueue(ctx sdk.Context, timestamp time.Time) {
	ctx.KVStore(k.storeKey).Delete(types.GetUnbondingDelegationTimeKey(timestamp))
}
