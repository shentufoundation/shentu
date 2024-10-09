package keeper

import (
	"context"

	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

// SetTaskParams sets the current task params to the global param store.
func (k Keeper) SetTaskParams(ctx context.Context, taskParams types.TaskParams) {
	k.paramSpace.Set(ctx, types.ParamsStoreKeyTaskParams, &taskParams)
}

// GetTaskParams gets the current task params from the global param store.
func (k Keeper) GetTaskParams(ctx context.Context) types.TaskParams {
	var taskParams types.TaskParams
	k.paramSpace.Get(ctx, types.ParamsStoreKeyTaskParams, &taskParams)
	return taskParams
}

// SetLockedPoolParams sets the current locked pool params to the global param store.
func (k Keeper) SetLockedPoolParams(ctx context.Context, poolParams types.LockedPoolParams) {
	k.paramSpace.Set(ctx, types.ParamsStoreKeyPoolParams, &poolParams)
}

// GetLockedPoolParams gets the current locked pool params from the global param store.
func (k Keeper) GetLockedPoolParams(ctx context.Context) types.LockedPoolParams {
	var poolParams types.LockedPoolParams
	k.paramSpace.Get(ctx, types.ParamsStoreKeyPoolParams, &poolParams)
	return poolParams
}
