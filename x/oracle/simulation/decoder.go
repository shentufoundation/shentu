package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

// NewDecodeStore unmarshals the KVPair's Value to the corresponding type of oracle module.
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.OperatorStoreKeyPrefix):
			var operatorA, operatorB types.Operator
			cdc.MustUnmarshalLengthPrefixed(kvA.Value, &operatorA)
			cdc.MustUnmarshalLengthPrefixed(kvB.Value, &operatorB)
			return fmt.Sprintf("%v\n%v", operatorA, operatorB)

		case bytes.Equal(kvA.Key[:1], types.WithdrawStoreKeyPrefix):
			var withdrawA, withdrawB types.Withdraw
			cdc.MustUnmarshalLengthPrefixed(kvA.Value, &withdrawA)
			cdc.MustUnmarshalLengthPrefixed(kvB.Value, &withdrawB)
			return fmt.Sprintf("%v\n%v", withdrawA, withdrawB)

		case bytes.Equal(kvA.Key[:1], types.TotalCollateralKeyPrefix):
			var totalA, totalB types.CoinsProto
			cdc.MustUnmarshalLengthPrefixed(kvA.Value, &totalA)
			cdc.MustUnmarshalLengthPrefixed(kvB.Value, &totalB)
			return fmt.Sprintf("%v\n%v", totalA.Coins, totalB.Coins)

		case bytes.Equal(kvA.Key[:1], types.TaskStoreKeyPrefix):
			var taskA, taskB types.Task
			cdc.MustUnmarshalLengthPrefixed(kvA.Value, &taskA)
			cdc.MustUnmarshalLengthPrefixed(kvB.Value, &taskB)
			return fmt.Sprintf("%v\n%v", taskA, taskB)

		case bytes.Equal(kvA.Key[:1], types.ClosingTaskStoreKeyPrefix):
			var taskIDsA, taskIDsB types.TaskIDs
			cdc.MustUnmarshalLengthPrefixed(kvA.Value, &taskIDsA)
			cdc.MustUnmarshalLengthPrefixed(kvB.Value, &taskIDsB)
			return fmt.Sprintf("%v\n%v", taskIDsA.TaskIds, taskIDsB.TaskIds)

		default:
			panic(fmt.Sprintf("invalid %s key prefix %X", types.ModuleName, kvA.Key[:1]))
		}
	}
}
