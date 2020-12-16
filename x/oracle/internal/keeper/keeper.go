package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/oracle/internal/types"
)

type Keeper struct {
	cdc           codec.BinaryMarshaler
	storeKey      sdk.StoreKey
	authKeeper    types.AccountKeeper
	distrKeeper   types.DistrKeeper
	stakingKeeper types.StakingKeeper
	supplyKeeper  types.SupplyKeeper
	paramSpace    types.ParamSubspace
}

func NewKeeper(cdc codec.BinaryMarshaler, storeKey sdk.StoreKey, authKeeper types.AccountKeeper, distriKeeper types.DistrKeeper,
	stakingKeeper types.StakingKeeper, supplyKeeper types.SupplyKeeper, paramSpace types.ParamSubspace) Keeper {
	return Keeper{
		cdc:           cdc,
		paramSpace:    paramSpace,
		storeKey:      storeKey,
		authKeeper:    authKeeper,
		distrKeeper:   distriKeeper,
		stakingKeeper: stakingKeeper,
		supplyKeeper:  supplyKeeper,
	}
}

// GetAuthKeeper returns the auth keeper wrapped in module keeper.
func (k Keeper) GetAuthKeeper() types.AccountKeeper {
	return k.authKeeper
}
