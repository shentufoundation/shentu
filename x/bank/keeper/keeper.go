// Package keeper implements custom bank keeper through CVM.
package keeper

import (
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/codec"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	"github.com/shentufoundation/shentu/v2/x/bank/types"
)

// Keeper is a wrapper of the basekeeper with CVM keeper.
type Keeper struct {
	bankKeeper.BaseKeeper
	ak types.AccountKeeper
}

// NewKeeper returns a new Keeper.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	ak types.AccountKeeper,
	blockedAddrs map[string]bool,
	authority string,
	logger log.Logger,
) Keeper {
	bk := bankKeeper.NewBaseKeeper(cdc, storeService, ak, blockedAddrs, authority, logger)
	return Keeper{
		BaseKeeper: bk,
		ak:         ak,
	}
}
