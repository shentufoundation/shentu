// Package keeper implements custom bank keeper through CVM.
package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	"github.com/shentufoundation/shentu/v2/x/bank/types"
)

// Keeper is a wrapper of the basekeeper with CVM keeper.
type Keeper struct {
	bankKeeper.BaseKeeper
	ak types.AccountKeeper
}

// NewKeeper returns a new Keeper.
func NewKeeper(cdc codec.BinaryCodec, storeKey storetypes.StoreKey, ak types.AccountKeeper,
	blockedAddrs map[string]bool, authority string) Keeper {
	bk := bankKeeper.NewBaseKeeper(cdc, storeKey, ak, blockedAddrs, authority)
	return Keeper{
		BaseKeeper: bk,
		ak:         ak,
	}
}
