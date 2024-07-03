// Package keeper implements custom bank keeper through CVM.
package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/shentufoundation/shentu/v2/x/bank/types"
)

// Keeper is a wrapper of the basekeeper with CVM keeper.
type Keeper struct {
	bankKeeper.BaseKeeper
	ak types.AccountKeeper
}

// NewKeeper returns a new Keeper.
func NewKeeper(cdc codec.BinaryCodec, storeKey storetypes.StoreKey, ak types.AccountKeeper, paramSpace paramsTypes.Subspace,
	blockedAddrs map[string]bool) Keeper {
	bk := bankKeeper.NewBaseKeeper(cdc, storeKey, ak, paramSpace, blockedAddrs)
	return Keeper{
		BaseKeeper: bk,
		ak:         ak,
	}
}
