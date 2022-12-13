package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

// Keeper - crisis keeper
type Keeper struct {
	storeKey sdk.StoreKey
	cdc codec.BinaryCodec
	paramSpace paramtypes.Subspace

	bk types.BankKeeper
}

// NewKeeper creates a new Keeper object
func NewKeeper(
	paramSpace paramtypes.Subspace, bankKeeper types.BankKeeper,
) Keeper {

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		paramSpace: paramSpace,
		bk:         bankKeeper,
	}
}
