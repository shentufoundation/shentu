package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shentufoundation/shentu/v2/x/oracle/types"
	"time"
)

// SetPrecogTask sets a precog task in KVStore.
func (k Keeper) SetPrecogTask(ctx sdk.Context, task types.PrecogTask) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PrecogTaskStoreKey(task.BusinessTxHash), k.cdc.MustMarshalLengthPrefixed(&task))
}

// GetPrecogTask get a score form KVStore
func (k Keeper) GetPrecogTask(ctx sdk.Context, hash string) (types.PrecogTask, error) {
	store := ctx.KVStore(k.storeKey)
	precogTaskData := store.Get(types.PrecogTaskStoreKey(hash))

	if precogTaskData == nil {
		return types.PrecogTask{}, types.ErrPrecogTaskNotExists
	}
	var precogTask types.PrecogTask
	k.cdc.MustUnmarshalLengthPrefixed(precogTaskData, &precogTask)
	return precogTask, nil
}

func (k Keeper) DeletePrecogTask(ctx sdk.Context, hash string) error {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.PrecogTaskStoreKey(hash))
	return nil
}

// CreatePrecogTask creates a new task.
func (k Keeper) CreatePrecogTask(ctx sdk.Context, creator, chainID string, bounty sdk.Coins, scoringWaitTime uint64, usageExpirationTime time.Time, businessTxHsh string) error {
	precogTaskData, err := k.GetPrecogTask(ctx, businessTxHsh)
	if err == nil {
		if usageExpirationTime.After(ctx.BlockTime()) {
			return types.ErrPrecogTaskNotClosed
		}

		if precogTaskData.Status == types.PrecogTaskStatusPending {
			precogTaskData.Status = types.PrecogTaskStatusCreated
		}
		if precogTaskData.Status == types.PrecogTaskStatusCreated {
			return types.ErrPrecogTaskNotClosed
		}

		if precogTaskData.Status == types.PrecogTaskStatusExpired || precogTaskData.Status == types.PrecogTaskStatusDone {
			k.DeletePrecogTask(ctx, businessTxHsh)
		}
	}

	precogTask := types.PrecogTask{
		Creator:             creator,
		ChainId:             chainID,
		Bounty:              bounty,
		ScoringWaitTime:     scoringWaitTime,
		CreateTime:          ctx.BlockTime(),
		UsageExpirationTime: usageExpirationTime,
		BusinessTxHash:      businessTxHsh,
		Status:              types.PrecogTaskStatusCreated,
	}
	k.SetPrecogTask(ctx, precogTask)
	return nil
}
