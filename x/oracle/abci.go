package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/oracle/internal/keeper"
)

func BeginBlocker(ctx sdk.Context, k Keeper) {
	k.FinalizeMatureWithdraws(ctx)
}

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	closingTaskIDs := k.GetClosingTaskIDs(ctx, ctx.BlockHeight())
	for _, taskID := range closingTaskIDs {
		err := k.Aggregate(ctx, taskID.Contract, taskID.Function)
		if err != nil {
			continue
		}
		task, err := k.GetTask(ctx, taskID.Contract, taskID.Function)
		if err != nil {
			continue
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				"aggregate_task",
				sdk.NewAttribute("contract", task.Contract),
				sdk.NewAttribute("function", task.Function),
				sdk.NewAttribute("result", task.Result.String()),
				sdk.NewAttribute("expiration", task.Expiration.String()),
			),
		)

		if err := k.DistributeBounty(ctx, task); err != nil {
			// TODO
			continue
		}
	}
	k.DeleteClosingTaskIDs(ctx, ctx.BlockHeight())
}
