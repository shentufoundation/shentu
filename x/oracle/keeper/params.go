package keeper

import (
	"context"

	"github.com/shentufoundation/shentu/v2/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetTaskParams sets the current task params to the global param store.
func (k Keeper) SetTaskParams(ctx context.Context, taskParams types.TaskParams) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	k.paramSpace.Set(sdkCtx, types.ParamsStoreKeyTaskParams, &taskParams)
}

// GetTaskParams gets the current task params from the global param store.
func (k Keeper) GetTaskParams(ctx context.Context) types.TaskParams {
	var taskParams types.TaskParams
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	k.paramSpace.Get(sdkCtx, types.ParamsStoreKeyTaskParams, &taskParams)
	return taskParams
}

// SetLockedPoolParams sets the current locked pool params to the global param store.
func (k Keeper) SetLockedPoolParams(ctx context.Context, poolParams types.LockedPoolParams) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	k.paramSpace.Set(sdkCtx, types.ParamsStoreKeyPoolParams, &poolParams)
}

// GetLockedPoolParams gets the current locked pool params from the global param store.
func (k Keeper) GetLockedPoolParams(ctx context.Context) types.LockedPoolParams {
	var poolParams types.LockedPoolParams
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	k.paramSpace.Get(sdkCtx, types.ParamsStoreKeyPoolParams, &poolParams)
	return poolParams
}
