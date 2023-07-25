package v280

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/common"
	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

func MigrateAllTaskStore(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.TaskStoreKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var oldTask types.TaskI

		err := cdc.UnmarshalInterface(iterator.Value(), &oldTask)
		if err != nil {
			return err
		}

		switch task := oldTask.(type) {
		case *types.Task:
			if err = MigrateTaskStore(task, store, iterator.Key(), cdc); err != nil {
				return err
			}
		case *types.TxTask:
			if err = MigrateTxTaskStore(task, store, iterator.Key(), cdc); err != nil {
				return err
			}
		default:
			return fmt.Errorf("err kvstore")
		}

	}
	return nil
}

func MigrateTaskStore(task *types.Task, store store.KVStore, key []byte, cdc codec.BinaryCodec) error {
	shentuAddr, err := common.PrefixToShentu(task.Creator)
	if err != nil {
		return err
	}

	newTask := types.Task{
		Contract:      task.Contract,
		Function:      task.Function,
		BeginBlock:    task.BeginBlock,
		Bounty:        task.Bounty,
		Description:   task.Description,
		Expiration:    task.Expiration,
		Creator:       shentuAddr,
		Responses:     nil,
		Result:        task.Result,
		ExpireHeight:  task.ExpireHeight,
		WaitingBlocks: task.WaitingBlocks,
		Status:        task.Status,
	}

	for _, response := range task.Responses {
		operator, err := common.PrefixToShentu(response.Operator)
		if err != nil {
			return err
		}
		newResponse := types.Response{
			Operator: operator,
			Score:    response.Score,
			Weight:   response.Weight,
			Reward:   response.Reward,
		}
		newTask.Responses = append(newTask.Responses, newResponse)
	}
	// delete old task
	store.Delete(key)
	// set task
	bz, err := cdc.MarshalInterface(&newTask)
	if err != nil {
		return err
	}
	store.Set(types.TaskStoreKey(newTask.GetID()), bz)
	return nil
}

func MigrateTxTaskStore(task *types.TxTask, store store.KVStore, key []byte, cdc codec.BinaryCodec) error {
	shentuAddr, err := common.PrefixToShentu(task.Creator)
	if err != nil {
		return err
	}

	newTask := types.TxTask{
		AtxHash:    task.AtxHash,
		Creator:    shentuAddr,
		Bounty:     task.Bounty,
		ValidTime:  task.ValidTime,
		Expiration: task.Expiration,
		Responses:  nil,
		Score:      task.Score,
		Status:     task.Status,
	}

	for _, response := range task.Responses {
		operator, err := common.PrefixToShentu(response.Operator)
		if err != nil {
			return err
		}
		newResponse := types.Response{
			Operator: operator,
			Score:    response.Score,
			Weight:   response.Weight,
			Reward:   response.Reward,
		}
		newTask.Responses = append(newTask.Responses, newResponse)
	}
	// delete old task
	store.Delete(key)
	// set task
	bz, err := cdc.MarshalInterface(&newTask)
	if err != nil {
		return err
	}
	store.Set(types.TaskStoreKey(newTask.GetID()), bz)
	return nil
}

func MigrateOperatorStore(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.OperatorStoreKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var operator types.Operator
		cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &operator)

		shentuOperatorAddress, err := common.PrefixToShentu(operator.Address)
		if err != nil {
			return err
		}
		shentuProposal, err := common.PrefixToShentu(operator.Proposer)
		if err != nil {
			return err
		}

		operator.Address = shentuOperatorAddress
		operator.Proposer = shentuProposal

		bz := cdc.MustMarshalLengthPrefixed(&operator)
		addr := sdk.MustAccAddressFromBech32(shentuOperatorAddress)
		store.Set(types.OperatorStoreKey(addr), bz)
	}
	return nil
}
