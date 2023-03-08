package oracle

import (
	"encoding/base64"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/oracle/keeper"
	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	k.FinalizeMatureWithdraws(ctx)
}

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	closingTaskIDs := k.GetClosingTaskIDs(ctx, nil)
	for _, taskID := range closingTaskIDs {
		err := k.Aggregate(ctx, taskID.Tid)
		if err != nil {
			continue
		}
		task, err := k.GetTask(ctx, taskID.Tid)
		if err != nil {
			continue
		}

		if err := k.DistributeBounty(ctx, task); err != nil {
			// TODO
			continue
		}

		switch task := task.(type) {
		case *types.Task:
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
		case *types.TxTask:
			//implement me
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					"aggregate_task",
					//shall be base64 encoded to be aligned with proto json encoding
					sdk.NewAttribute("tx_hash", base64.StdEncoding.EncodeToString(task.TxHash)),
				),
			)
		}
	}
	k.DeleteClosingTaskIDs(ctx)
}
