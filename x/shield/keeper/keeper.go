package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/shentufoundation/shentu/v2/x/shield/types"
)

// Keeper implements the shield keeper.
type Keeper struct {
	storeKey   storetypes.StoreKey
	cdc        codec.BinaryCodec
	ak         types.AccountKeeper
	bk         types.BankKeeper
	paramSpace types.ParamSubspace
}

// NewKeeper creates a shield keeper.
func NewKeeper(cdc codec.BinaryCodec, shieldStoreKey storetypes.StoreKey, ak types.AccountKeeper, bk types.BankKeeper,
	paramSpace types.ParamSubspace) Keeper {
	return Keeper{
		storeKey:   shieldStoreKey,
		cdc:        cdc,
		ak:         ak,
		bk:         bk,
		paramSpace: paramSpace,
	}
}
