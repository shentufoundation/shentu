package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// SetOperator sets the Shield Operator account address.
func (k Keeper) SetOperator(ctx sdk.Context, operator sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	bz := operator
	store.Set(types.GetShieldOperatorKey(), bz)
}

// GetOperator gets the Shield Operator account address.
func (k Keeper) GetOperator(ctx sdk.Context) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	return store.Get(types.GetShieldOperatorKey())
}
