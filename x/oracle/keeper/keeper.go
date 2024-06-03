package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

type Keeper struct {
	cdc           codec.BinaryCodec
	storeKey      storetypes.StoreKey
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	distrKeeper   types.DistrKeeper
	stakingKeeper types.StakingKeeper
	CertKeeper    types.CertKeeper
	paramSpace    types.ParamSubspace
}

func NewKeeper(cdc codec.BinaryCodec, storeKey storetypes.StoreKey, authKeeper types.AccountKeeper, distriKeeper types.DistrKeeper,
	stakingKeeper types.StakingKeeper, bankKeeper types.BankKeeper, certKeeper types.CertKeeper, paramSpace types.ParamSubspace) Keeper {
	return Keeper{
		cdc:           cdc,
		paramSpace:    paramSpace,
		storeKey:      storeKey,
		accountKeeper: authKeeper,
		distrKeeper:   distriKeeper,
		stakingKeeper: stakingKeeper,
		bankKeeper:    bankKeeper,
		CertKeeper:    certKeeper,
	}
}

// GetAccountKeeper returns the auth keeper wrapped in module keeper.
func (k Keeper) GetAccountKeeper() types.AccountKeeper {
	return k.accountKeeper
}
