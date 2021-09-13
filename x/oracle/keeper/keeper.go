package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/x/oracle/types"
)

type Keeper struct {
	cdc           codec.BinaryCodec
	storeKey      sdk.StoreKey
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	distrKeeper   types.DistrKeeper
	stakingKeeper types.StakingKeeper
	paramSpace    types.ParamSubspace
}

func NewKeeper(cdc codec.BinaryCodec, storeKey sdk.StoreKey, authKeeper types.AccountKeeper, distriKeeper types.DistrKeeper,
	stakingKeeper types.StakingKeeper, bankKeeper types.BankKeeper, paramSpace types.ParamSubspace) Keeper {
	return Keeper{
		cdc:           cdc,
		paramSpace:    paramSpace,
		storeKey:      storeKey,
		accountKeeper: authKeeper,
		distrKeeper:   distriKeeper,
		stakingKeeper: stakingKeeper,
		bankKeeper:    bankKeeper,
	}
}

// GetAccountKeeper returns the auth keeper wrapped in module keeper.
func (k Keeper) GetAccountKeeper() types.AccountKeeper {
	return k.accountKeeper
}
