package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/x/oracle/types"
)

// SetTaskParams sets the current task params to the global param store.
func (k Keeper) SetTaskParams(ctx sdk.Context, taskParams types.TaskParams) {
	k.paramSpace.Set(ctx, types.ParamsStoreKeyTaskParams, &taskParams)
}

// GetTaskParams gets the current task params from the global param store.
func (k Keeper) GetTaskParams(ctx sdk.Context) types.TaskParams {
	var taskParams types.TaskParams
	k.paramSpace.Get(ctx, types.ParamsStoreKeyTaskParams, &taskParams)
	return taskParams
}

// SetLockedPoolParams sets the current locked pool params to the global param store.
func (k Keeper) SetLockedPoolParams(ctx sdk.Context, poolParams types.LockedPoolParams) {
	k.paramSpace.Set(ctx, types.ParamsStoreKeyPoolParams, &poolParams)
}

// GetLockedPoolParams gets the current locked pool params from the global param store.
func (k Keeper) GetLockedPoolParams(ctx sdk.Context) types.LockedPoolParams {
	var poolParams types.LockedPoolParams
	k.paramSpace.Get(ctx, types.ParamsStoreKeyPoolParams, &poolParams)
	return poolParams
}
