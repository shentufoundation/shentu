package keeper

import (
	"cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/codec"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

type Keeper struct {
	cdc           codec.BinaryCodec
	storeService  store.KVStoreService
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	distrKeeper   types.DistrKeeper
	stakingKeeper types.StakingKeeper
	CertKeeper    types.CertKeeper
	paramSpace    paramtypes.Subspace
}

func NewKeeper(cdc codec.BinaryCodec, storeService store.KVStoreService, authKeeper types.AccountKeeper, distriKeeper types.DistrKeeper,
	stakingKeeper types.StakingKeeper, bankKeeper types.BankKeeper, certKeeper types.CertKeeper, paramSpace paramtypes.Subspace) Keeper {
	return Keeper{
		cdc:           cdc,
		paramSpace:    paramSpace,
		storeService:  storeService,
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
