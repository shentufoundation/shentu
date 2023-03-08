package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/oracle/keeper"
	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

// InitGenesis puts all data from genesis state into store.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	operators := data.Operators
	totalCollateral := data.TotalCollateral
	poolParams := data.PoolParams
	taskParams := data.TaskParams
	withdraws := data.Withdraws
	tasks := data.Tasks

	for _, operator := range operators {
		k.SetOperator(ctx, operator)
	}

	k.SetTotalCollateral(ctx, totalCollateral)
	k.SetLockedPoolParams(ctx, *poolParams)
	k.SetTaskParams(ctx, *taskParams)

	for _, withdraw := range withdraws {
		withdraw.DueBlock += ctx.BlockHeight()
		k.SetWithdraw(ctx, withdraw)
	}

	for i := range tasks {
		task := tasks[i]
		task.ClosingBlock = ctx.BlockHeight() + task.WaitingBlocks
		k.UpdateAndSetTask(ctx, &task)
	}
}

// ExportGenesis extracts all data from store to genesis state.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	operators := k.GetAllOperators(ctx)
	totalCollateral, _ := k.GetTotalCollateral(ctx)

	poolParams := k.GetLockedPoolParams(ctx)
	taskParams := k.GetTaskParams(ctx)
	withdraws := k.GetAllWithdrawsForExport(ctx)

	tasks := k.UpdateAndGetAllTasks(ctx)

	//TODO: reimplement this to take both Task and TxTask
	var smartContractTasks []types.Task
	for _, t := range tasks {
		if sct, ok := t.(*types.Task); ok {
			smartContractTasks = append(smartContractTasks, *sct)
		}
	}
	return types.NewGenesisState(operators, totalCollateral, poolParams, taskParams, withdraws, smartContractTasks)
}
