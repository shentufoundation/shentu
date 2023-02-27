package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

// SetTxTask sets a tx task in KVStore.
func (k Keeper) SetTxTask(ctx sdk.Context, task types.TxTask) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.TxTaskStoreKey(task.TxHash), k.cdc.MustMarshalLengthPrefixed(&task))
}

// GetTxTask get a txTask form KVStore
func (k Keeper) GetTxTask(ctx sdk.Context, hash []byte) (types.TxTask, error) {
	store := ctx.KVStore(k.storeKey)
	txTaskData := store.Get(types.TxTaskStoreKey(hash))

	if txTaskData == nil {
		return types.TxTask{}, types.ErrTxTaskNotExists
	}
	var txTask types.TxTask
	k.cdc.MustUnmarshalLengthPrefixed(txTaskData, &txTask)
	return txTask, nil
}

func (k Keeper) DeleteTxTask(ctx sdk.Context, hash []byte) error {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.TxTaskStoreKey(hash))
	return nil
}

// CreateTxTask creates a new tx task.
func (k Keeper) CreateTxTask(ctx sdk.Context, creator string, bounty sdk.Coins, expirationTime time.Time, txHash []byte) error {
	txTaskData, err := k.GetTxTask(ctx, txHash)
	if err == nil {
		if expirationTime.Before(ctx.BlockTime()) {
			return types.ErrTxTaskExpirationTime
		}

		if txTaskData.Status != types.TaskStatusPending {
			return types.ErrTxTaskNotClosed
		} else if txTaskData.Creator != creator || txTaskData.Expiration != expirationTime {
			return types.ErrInvalidTxTask
		} else {
			for i, coin := range txTaskData.Bounty {
				if coin != bounty[i] {
					return types.ErrInvalidTxTask
				}
			}
		}
	}

	txTask := types.TxTask{
		Creator:    creator,
		TxHash:     txHash,
		Bounty:     bounty,
		Expiration: expirationTime,
		Status:     types.TaskStatusNil,
	}
	k.SetTxTask(ctx, txTask)
	return nil
}
