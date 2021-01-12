package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/gov/types"
)

// GetDepositParams returns the current DepositParams from the global param store.
func (k Keeper) GetDepositParams(ctx sdk.Context) types.DepositParams {
	var depositParams types.DepositParams
	k.paramSpace.Get(ctx, types.ParamStoreKeyDepositParams, &depositParams)
	return depositParams
}

// GetTallyParams returns the current TallyParams from the global param store.
func (k Keeper) GetTallyParams(ctx sdk.Context) types.TallyParams {
	var tallyParams types.TallyParams
	k.paramSpace.Get(ctx, types.ParamStoreKeyTallyParams, &tallyParams)
	return tallyParams
}

// SetDepositParams sets parameters space for deposits.
func (k Keeper) SetDepositParams(ctx sdk.Context, depositParams types.DepositParams) {
	k.paramSpace.Set(ctx, types.ParamStoreKeyDepositParams, &depositParams)
}

// SetTallyParams sets parameters space for tally period.
func (k Keeper) SetTallyParams(ctx sdk.Context, tallyParams types.TallyParams) {
	k.paramSpace.Set(ctx, types.ParamStoreKeyTallyParams, &tallyParams)
}
