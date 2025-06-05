package keeper

import (
	"cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/codec"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/shentufoundation/shentu/v2/x/shield/types"
)

// Keeper implements the shield keeper.
type Keeper struct {
	storeService store.KVStoreService
	cdc          codec.BinaryCodec
	ak           types.AccountKeeper
	bk           types.BankKeeper
	paramSpace   paramtypes.Subspace
}

// NewKeeper creates a shield keeper.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	paramSpace paramtypes.Subspace,
) Keeper {
	return Keeper{
		storeService: storeService,
		cdc:          cdc,
		ak:           ak,
		bk:           bk,
		paramSpace:   paramSpace,
	}
}
