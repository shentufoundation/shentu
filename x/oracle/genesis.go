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
	leftBounties := data.LeftBounties

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
		k.UpdateAndSetTask(ctx, &tasks[i])
	}
	for i := range txTasks {
		k.SetTxTask(ctx, &txTasks[i])
	}

	for _, bounty := range leftBounties {
		k.SetCreatorLeftBounty(ctx, bounty)
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
	leftBounties := k.GetAllLeftBounties(ctx)
	return types.NewGenesisState(operators, totalCollateral, poolParams, taskParams, withdraws, tasks, txTasks, leftBounties)
}
