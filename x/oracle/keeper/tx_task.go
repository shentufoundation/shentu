package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

// CreateTxTask creates a new tx task.
func (k Keeper) CreateTxTask(ctx sdk.Context, task types.TaskI) error {
	txTaskData, err := k.GetTask(ctx, task.GetID())
	if err == nil {
		if !task.IsValid(ctx) {
			return types.ErrInvalidTxTask
		}

		if txTaskData.GetStatus() != types.TaskStatusPending {
			return types.ErrTxTaskNotClosed
		} else if txTaskData.GetCreator() != task.GetCreator() || txTaskData.GetValidTime() != task.GetValidTime() {
			return types.ErrInvalidTxTask
		} else {
			for i, coin := range txTaskData.GetBounty() {
				if coin != task.GetBounty()[i] {
					return types.ErrInvalidTxTask
				}
			}
		}
	}

	k.SetTask(ctx, task)
	return nil
}
