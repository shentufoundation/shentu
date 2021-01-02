package oracle

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/oracle/keeper"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
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

		if err := k.DistributeBounty(ctx, task); err != nil {
			// TODO
			continue
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				"aggregate_task",
				sdk.NewAttribute("contract", task.Contract),
				sdk.NewAttribute("function", task.Function),
				sdk.NewAttribute("begin_block_height", strconv.FormatInt(task.BeginBlock, 10)),
				sdk.NewAttribute("bounty", task.Bounty.String()),
				sdk.NewAttribute("description", task.Description),
				sdk.NewAttribute("expiration", task.Expiration.String()),
				sdk.NewAttribute("creator", task.Creator),
				sdk.NewAttribute("responses", task.Responses.String()),
				sdk.NewAttribute("result", task.Result.String()),
				sdk.NewAttribute("end_block_height", strconv.FormatInt(task.ClosingBlock, 10)),
				sdk.NewAttribute("status", task.Status.String()),
			),
		)
	}
	k.DeleteClosingTaskIDs(ctx, ctx.BlockHeight())
}
