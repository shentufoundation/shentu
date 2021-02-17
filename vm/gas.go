package vm

import (
	"bytes"
	"errors"

	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/execution/engine"
	. "github.com/hyperledger/burrow/execution/evm"
)

// These are based on Ethereum yellow paper.
const (
	// Basic gas units
	GasBase    uint64 = 2
	GasVeryLow uint64 = 3
	GasLow     uint64 = 5
	GasMid     uint64 = 8
	GasHigh    uint64 = 10
	GasExtStep uint64 = 20 // homestead extcode cost

	// Constantinople
	GasExtcodeSize          uint64 = 700
	GasExtcodeCopy          uint64 = 700
	GasExtcodeHash          uint64 = 700
	GasBalance              uint64 = 700
	GasSLoad                uint64 = 800
	GasCalls                uint64 = 700
	GasSelfdestruct         uint64 = 5000
	GasExpByte              uint64 = 50
	GasCreateBySelfdestruct uint64 = 25000

	CallValueTransferGas uint64 = 9000  // Paid for CALL when the value transfer is non-zero.
	CallNewAccountGas    uint64 = 25000 // Paid for CALL when the destination address didn't exist prior.
	QuadCoeffDiv         uint64 = 512   // Divisor for the quadratic particle of the memory cost equation.
	LogDataGas           uint64 = 8     // Per byte in a LOG* operation's data.

	Sha3Gas     uint64 = 30 // Once per SHA3 operation.
	Sha3WordGas uint64 = 6  // Once per word of the SHA3 operation's data.

	SstoreSetGas    uint64 = 20000 // Once per SLOAD operation.
	SstoreResetGas  uint64 = 5000  // Once per SSTORE operation if the zeroness changes from zero.
	SstoreClearGas  uint64 = 5000  // Once per SSTORE operation if the zeroness doesn't change.
	SstoreRefundGas uint64 = 15000 // Once per SSTORE operation if the zeroness changes to zero.

	NetSstoreNoopGas uint64 = 200 // Once per SSTORE operation if the value doesn't change.

	JumpdestGas           uint64 = 1     // Once per JUMPDEST operation.
	CallGas               uint64 = 40    // Once per CALL operation & message call transaction.
	ExpGas                uint64 = 10    // Once per EXP instruction
	LogGas                uint64 = 375   // Per LOG* operation.
	CopyGas               uint64 = 3     //
	LogTopicGas           uint64 = 375   // Multiplied by the * of the LOG*, per LOG transaction. e.g. LOG0 incurs 0 * c_txLogTopicGas, LOG4 incurs 4 * c_txLogTopicGas.
	CreateGas             uint64 = 32000 // Once per CREATE operation & contract-creation transaction.
	SelfdestructRefundGas uint64 = 24000 // Refunded following a suicide operation.
	MemoryGas             uint64 = 3     // Times the address of the (highest referenced byte in memory + 1). NOTE: referencing happens on read, write and in instructions such as RETURN and CALL.

	ExtcodeCopyBase = 20
)

const (
	uint64Length = 8
)

var errGasUintOverflow = errors.New("gas uint64 overflow")

// gasMemory is a custom memory type, with track of last gas cost and refund counter.
type gasMemory struct {
	engine.Memory
	lastGasCost uint64
	refund      uint64
}

// memGasCost calculates the additional gas cost based on memory usage.
// It is reused in multiple of the subsequent functions.
func memGasCost(mem *gasMemory, newMemSize uint64) (uint64, error) {
	if newMemSize == 0 {
		return 0, nil
	}
	// The maximum that will fit in a uint64 is max_word_count - 1. Anything above
	// that will result in an overflow. Additionally, a newMemSize which results in
	// a newMemSizeWords larger than 0xFFFFFFFF will cause the square operation to
	// overflow. The constant 0x1FFFFFFFE0 is the highest number that can be used
	// without overflowing the gas calculation.
	if newMemSize > 0x1FFFFFFFE0 {
		return 0, errGasUintOverflow
	}
	newMemSizeWords := toWordSize(newMemSize)
	newMemSize = newMemSizeWords * 32

	if newMemSize <= mem.Capacity().Uint64() {
		return 0, nil
	}
	square := newMemSizeWords * newMemSizeWords
	linCoef := newMemSizeWords * MemoryGas
	quadCoef := square / QuadCoeffDiv
	newTotalFee := linCoef + quadCoef

	fee := newTotalFee - mem.lastGasCost
	mem.lastGasCost = newTotalFee

	return fee, nil
}

// memoryGas handles all simple dynamic gas calculation that depends on memory size of the operation.
func memoryGas(mem *gasMemory, memorySize uint64, gasAdd uint64) (uint64, error) {
	gas, err := memGasCost(mem, memorySize)
	if err != nil {
		return gas, err
	}

	gas, of := SafeAdd(gas, gasAdd)
	if of {
		return 0, errGasUintOverflow
	}

	return gas, nil
}

// copyGas handles all dynamic gas calculation that copies data.
func copyGas(stack *Stack, mem *gasMemory, memorySize uint64, n int, gasAdd uint64) (uint64, error) {
	gas, err := memGasCost(mem, memorySize)
	if err != nil {
		return 0, err
	}

	var words uint64
	var of bool
	if words, of = bigToUint64(Get(stack, n)); of {
		return 0, errGasUintOverflow
	}
	if words, of = SafeMul(toWordSize(words), gasAdd); of {
		return 0, errGasUintOverflow
	}
	if gas, of = SafeAdd(gas, words); of {
		return 0, errGasUintOverflow
	}
	return gas, nil
}

func onlyMemGasCost(st engine.CallFrame, address crypto.Address, stack *Stack, mem *gasMemory, memorySize uint64) (uint64, error) {
	return memGasCost(mem, memorySize)
}

func onlyMemoryGas(addGas uint64) gasFunc {
	return func(st engine.CallFrame, address crypto.Address, stack *Stack, mem *gasMemory, memorySize uint64) (uint64, error) {
		return memoryGas(mem, memorySize, addGas)
	}
}

func onlyCopyGas(stackPos int, gasBase uint64, gasAdd uint64) gasFunc {
	return func(st engine.CallFrame, address crypto.Address, stack *Stack, mem *gasMemory, memorySize uint64) (uint64, error) {
		return copyGas(stack, mem, memorySize, stackPos, gasAdd)
	}
}

var (
	gasMLoad   = onlyMemGasCost
	gasReturn  = onlyMemGasCost
	gasRevert  = onlyMemGasCost
	gasMStore8 = onlyMemGasCost
	gasMStore  = onlyMemGasCost
	gasCreate  = onlyMemGasCost

	gasDelegateCall = onlyMemoryGas(GasCalls)
	gasStaticCall   = onlyMemoryGas(GasCalls)

	gasCreate2        = onlyCopyGas(2, CreateGas, Sha3WordGas)
	gasCallDataCopy   = onlyCopyGas(2, GasVeryLow, CopyGas)
	gasReturnDataCopy = onlyCopyGas(2, GasVeryLow, CopyGas)
	gasCodeCopy       = onlyCopyGas(2, GasVeryLow, CopyGas)
	gasExtCodeCopy    = onlyCopyGas(3, GasExtcodeCopy, CopyGas)
	gasSha3           = onlyCopyGas(1, Sha3Gas, Sha3WordGas)
)

func makeGasLog(n uint64) gasFunc {
	return func(st engine.CallFrame, address crypto.Address, stack *Stack, mem *gasMemory, memorySize uint64) (uint64, error) {
		requestedSize, of := bigToUint64(Get(stack, 1))
		if of {
			return 0, errGasUintOverflow
		}

		gas, err := memGasCost(mem, memorySize)
		if err != nil {
			return 0, err
		}
		if gas, of = SafeAdd(gas, LogGas); of {
			return 0, errGasUintOverflow
		}
		if gas, of = SafeAdd(gas, n*LogTopicGas); of {
			return 0, errGasUintOverflow
		}

		memorySizeGas, of := SafeMul(requestedSize, LogDataGas)
		if of {
			return 0, errGasUintOverflow
		}
		if gas, of = SafeAdd(gas, memorySizeGas); of {
			return 0, errGasUintOverflow
		}
		return gas, nil
	}
}

func gasExp(st engine.CallFrame, address crypto.Address, stack *Stack, mem *gasMemory, memorySize uint64) (uint64, error) {
	expByteLen := uint64(len(Get(stack, 1).Bytes()))
	gas := expByteLen * GasExpByte // no overflow check required. Max is 256 * GasExpByte gas

	gas, of := SafeAdd(gas, ExpGas)
	if of {
		return 0, errGasUintOverflow
	}
	return gas, nil
}

func gasCall(st engine.CallFrame, address crypto.Address, stack *Stack, mem *gasMemory, memorySize uint64) (uint64, error) {
	var (
		gas            = GasCalls
		transfersValue = Get(stack, 2).Sign() != 0
		address2       = crypto.AddressFromWord256(GetWord256(stack, 1))
		// eip158         = vm.ChainConfig().IsEIP158(vm.BlockNumber)
	)
	acc, _ := st.GetAccount(address2)
	if transfersValue && acc == nil {
		gas += CallNewAccountGas
	}
	if transfersValue {
		gas += CallValueTransferGas
	}
	return memoryGas(mem, memorySize, gas)
}

func gasCallCode(st engine.CallFrame, address crypto.Address, stack *Stack, mem *gasMemory, memorySize uint64) (uint64, error) {
	gas := GasCalls
	if Get(stack, 2).Sign() != 0 {
		gas += CallValueTransferGas
	}
	return memoryGas(mem, memorySize, gas)
}

func gasSStore(st engine.CallFrame, address crypto.Address, stack *Stack, mem *gasMemory, memorySize uint64) (uint64, error) {
	var (
		x, y       = GetWord256(stack, 0), GetWord256(stack, 1)
		current, _ = st.GetStorage(address, x)
	)
	// The SStore takes some situations of the change of the stored value's state.
	// The resulting gas price is decided by the change of the state.
	//	1. From a zero-value address to a non-zero value         (NEW VALUE)
	//	2. From a non-zero value address to a zero-value address (DELETE)
	//	3. From a non-zero to a non-zero                         (CHANGE)
	//  4. (Additional) Value doesn't change                     (NOOP)

	switch {
	case bytes.Equal(current, binary.Zero256.Bytes()) && y != binary.Zero256: // 0 => non 0
		return SstoreSetGas, nil
	case !bytes.Equal(current, binary.Zero256.Bytes()) && y == binary.Zero256: // non 0 => 0
		mem.refund += SstoreRefundGas
		return SstoreClearGas, nil

	// Not too much complications with dirty state, but state noop added.
	case bytes.Equal(current, y.Bytes()):
		return NetSstoreNoopGas, nil
	default: // non 0 => non 0 (or 0 => 0)
		return SstoreResetGas, nil
	}
}

// gasSelfdestruct is called when contract self-destructs, freeing CVM memory.
// When contract successfully kills itself, some amount of gas is refunded.
func gasSelfdestruct(st engine.CallFrame, address crypto.Address, stack *Stack, mem *gasMemory, memorySize uint64) (uint64, error) {
	gas := GasSelfdestruct

	address2 := crypto.AddressFromWord256(GetWord256(stack, 0))

	// if empty and transfers value
	acc, err := st.GetAccount(address2)
	if acc == nil {
		gas += GasCreateBySelfdestruct
	}
	if err != nil {
		return gas, err
	}
	refund := mem.refund
	mem.refund += SelfdestructRefundGas
	if mem.refund < refund {
		return gas, errGasUintOverflow
	}
	return gas, nil
}
