package types

import (
	"encoding/binary"
	"encoding/hex"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of this module
	ModuleName = "oracle"

	// RouterKey is used to route messages.
	RouterKey = ModuleName

	// StoreKey is the prefix under which we store this module's data.
	StoreKey = ModuleName

	// QuerierRoute is used to handle abci_query requests.
	QuerierRoute = ModuleName
)

var (
	OperatorStoreKeyPrefix         = []byte{0x01}
	WithdrawStoreKeyPrefix         = []byte{0x02}
	TotalCollateralKeyPrefix       = []byte{0x03}
	TaskStoreKeyPrefix             = []byte{0x04}
	ClosingTaskStoreKeyPrefix      = []byte{0x05}
	ClosingTaskStoreKeyTimedPrefix = []byte{0x06}
	ExpireTaskStoreKeyPrefix       = []byte{0x07}
	ShortcutTasksKeyPrefix         = []byte{0x08}
)

func OperatorStoreKey(operator sdk.AccAddress) []byte {
	return append(OperatorStoreKeyPrefix, operator.Bytes()...)
}

func WithdrawStoreKey(address sdk.AccAddress, dueBlock int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(dueBlock))
	return append(append(WithdrawStoreKeyPrefix, b...), address.Bytes()...)
}

func TotalCollateralKey() []byte {
	return TotalCollateralKeyPrefix
}

func TaskStoreKey(taskID []byte) []byte {
	return append(TaskStoreKeyPrefix, taskID...)
}

func ClosingTaskIDsStoreKey(blockHeight int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(blockHeight))
	return append(ClosingTaskStoreKeyPrefix, b...)
}

func TimeStoreKey(prefix []byte, theTime time.Time) []byte {
	b := sdk.FormatTimeBytes(theTime)
	return append(prefix, b...)
}

func NewTaskID(contract, function string) []byte {
	return append([]byte(contract), []byte(function)...)
}

func NewAtxTaskID(txHash string) ([]byte, error) {
	return hex.DecodeString(txHash)
}
