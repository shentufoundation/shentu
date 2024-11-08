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
	txTasks := data.TxTasks

	for _, operator := range operators {
		if err := k.SetOperator(ctx, operator); err != nil {
			panic(err)
		}
	}

	if err := k.SetTotalCollateral(ctx, totalCollateral); err != nil {
		panic(err)
	}
	k.SetLockedPoolParams(ctx, *poolParams)
	k.SetTaskParams(ctx, *taskParams)

	for _, withdraw := range withdraws {
		withdraw.DueBlock += ctx.BlockHeight()
		if err := k.SetWithdraw(ctx, withdraw); err != nil {
			panic(err)
		}
	}

	for i := range tasks {
		if err := k.UpdateAndSetTask(ctx, &tasks[i]); err != nil {
			panic(err)
		}
	}
	for i := range txTasks {
		if err := k.SetTxTask(ctx, &txTasks[i]); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis extracts all data from store to genesis state.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	operators := k.GetAllOperators(ctx)
	totalCollateral, _ := k.GetTotalCollateral(ctx)

	poolParams := k.GetLockedPoolParams(ctx)
	taskParams := k.GetTaskParams(ctx)
	withdraws := k.GetAllWithdrawsForExport(ctx)
	tasks, txTasks := k.UpdateAndGetAllTasks(ctx)

	return types.NewGenesisState(operators, totalCollateral, poolParams, taskParams, withdraws, tasks, txTasks)
}
