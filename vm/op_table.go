package vm

import (
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/execution/engine"
	. "github.com/hyperledger/burrow/execution/evm"
	. "github.com/hyperledger/burrow/execution/evm/asm"
)

// Dynamic gas calculation logic for certain instructions
type gasFunc func(engine.CallFrame, crypto.Address, *Stack, *gasMemory, uint64) (uint64, error)

// Can implement in the next iteration.
// opFunc func() ([]byte, error)

// make sure the default value, 0 is not assigned to any specific operation
// then we can assume memSize == 0 means no memory size
const (
	memorySha3 = 1 + iota
	memoryLog
	memoryReturn
	memoryRevert
	memoryCallDataCopy
	memoryReturnDataCopy
	memoryCodeCopy
	memoryExtCodeCopy
	memoryMLoad
	memoryMStore
	memoryMStore8
	memoryCreate
	memoryCreate2
	memoryCall
	memoryDelegateCall
	memoryStaticCall
)

type instruction struct {
	dynamicGas gasFunc
	staticGas  uint64

	// memsize enum
	memSize uint
}

var instructionSet = [256]instruction{
	STOP: {
		staticGas: 0,
	},
	ADD: {
		staticGas: GasVeryLow,
	},
	MUL: {
		staticGas: GasLow,
	},
	SUB: {
		staticGas: GasVeryLow,
	},
	DIV: {
		staticGas: GasLow,
	},
	SDIV: {
		staticGas: GasLow,
	},
	MOD: {
		staticGas: GasLow,
	},
	SMOD: {
		staticGas: GasLow,
	},
	ADDMOD: {
		staticGas: GasMid,
	},
	MULMOD: {
		staticGas: GasMid,
	},
	EXP: {
		dynamicGas: gasExp,
	},
	SIGNEXTEND: {
		staticGas: GasLow,
	},
	LT: {
		staticGas: GasVeryLow,
	},
	GT: {
		staticGas: GasVeryLow,
	},
	SLT: {
		staticGas: GasVeryLow,
	},
	SGT: {
		staticGas: GasVeryLow,
	},
	EQ: {
		staticGas: GasVeryLow,
	},
	ISZERO: {
		staticGas: GasVeryLow,
	},
	AND: {
		staticGas: GasVeryLow,
	},
	XOR: {
		staticGas: GasVeryLow,
	},
	OR: {
		staticGas: GasVeryLow,
	},
	NOT: {
		staticGas: GasVeryLow,
	},
	BYTE: {
		staticGas: GasVeryLow,
	},
	SHA3: {
		staticGas:  Sha3Gas,
		dynamicGas: gasSha3,
		memSize:    memorySha3,
	},
	ADDRESS: {
		staticGas: GasBase,
	},
	BALANCE: {
		staticGas: GasBalance,
	},
	ORIGIN: {
		staticGas: GasBase,
	},
	CALLER: {
		staticGas: GasBase,
	},
	CALLVALUE: {
		staticGas: GasBase,
	},
	CALLDATALOAD: {
		staticGas: GasVeryLow,
	},
	CALLDATASIZE: {
		staticGas: GasBase,
	},
	CALLDATACOPY: {
		staticGas:  GasVeryLow,
		dynamicGas: gasCallDataCopy,
		memSize:    memoryCallDataCopy,
	},
	CODESIZE: {
		staticGas: GasBase,
	},
	CODECOPY: {
		staticGas:  GasVeryLow,
		dynamicGas: gasCodeCopy,
		memSize:    memoryCodeCopy,
	},
	GASPRICE_DEPRECATED: {
		staticGas: GasBase,
	},
	EXTCODESIZE: {
		staticGas: GasExtcodeSize,
	},
	EXTCODECOPY: {
		staticGas:  ExtcodeCopyBase,
		dynamicGas: gasExtCodeCopy,
		memSize:    memoryExtCodeCopy,
	},
	BLOCKHASH: {
		staticGas: GasExtStep,
	},
	COINBASE: {
		staticGas: GasBase,
	},
	TIMESTAMP: {
		staticGas: GasBase,
	},
	BLOCKHEIGHT: {
		staticGas: GasBase,
	},
	DIFFICULTY_DEPRECATED: {
		staticGas: GasBase,
	},
	GASLIMIT: {
		staticGas: GasBase,
	},
	POP: {
		staticGas: GasBase,
	},
	MLOAD: {
		staticGas:  GasVeryLow,
		dynamicGas: gasMLoad,
		memSize:    memoryMLoad,
	},
	MSTORE: {
		staticGas:  GasVeryLow,
		dynamicGas: gasMStore,
		memSize:    memoryMStore,
	},
	MSTORE8: {
		staticGas:  GasVeryLow,
		dynamicGas: gasMStore8,
		memSize:    memoryMStore8,
	},
	SLOAD: {
		staticGas: GasSLoad,
	},
	SSTORE: {
		dynamicGas: gasSStore,
	},
	JUMP: {
		staticGas: GasMid,
	},
	JUMPI: {
		staticGas: GasHigh,
	},
	PC: {
		staticGas: GasBase,
	},
	MSIZE: {
		staticGas: GasBase,
	},
	GAS: {
		staticGas: GasBase,
	},
	JUMPDEST: {
		staticGas: JumpdestGas,
	},
	PUSH1: {
		staticGas: GasVeryLow,
	},
	PUSH2: {
		staticGas: GasVeryLow,
	},
	PUSH3: {
		staticGas: GasVeryLow,
	},
	PUSH4: {
		staticGas: GasVeryLow,
	},
	PUSH5: {
		staticGas: GasVeryLow,
	},
	PUSH6: {
		staticGas: GasVeryLow,
	},
	PUSH7: {
		staticGas: GasVeryLow,
	},
	PUSH8: {
		staticGas: GasVeryLow,
	},
	PUSH9: {
		staticGas: GasVeryLow,
	},
	PUSH10: {
		staticGas: GasVeryLow,
	},
	PUSH11: {
		staticGas: GasVeryLow,
	},
	PUSH12: {
		staticGas: GasVeryLow,
	},
	PUSH13: {
		staticGas: GasVeryLow,
	},
	PUSH14: {
		staticGas: GasVeryLow,
	},
	PUSH15: {
		staticGas: GasVeryLow,
	},
	PUSH16: {
		staticGas: GasVeryLow,
	},
	PUSH17: {
		staticGas: GasVeryLow,
	},
	PUSH18: {
		staticGas: GasVeryLow,
	},
	PUSH19: {
		staticGas: GasVeryLow,
	},
	PUSH20: {
		staticGas: GasVeryLow,
	},
	PUSH21: {
		staticGas: GasVeryLow,
	},
	PUSH22: {
		staticGas: GasVeryLow,
	},
	PUSH23: {
		staticGas: GasVeryLow,
	},
	PUSH24: {
		staticGas: GasVeryLow,
	},
	PUSH25: {
		staticGas: GasVeryLow,
	},
	PUSH26: {
		staticGas: GasVeryLow,
	},
	PUSH27: {
		staticGas: GasVeryLow,
	},
	PUSH28: {
		staticGas: GasVeryLow,
	},
	PUSH29: {
		staticGas: GasVeryLow,
	},
	PUSH30: {
		staticGas: GasVeryLow,
	},
	PUSH31: {
		staticGas: GasVeryLow,
	},
	PUSH32: {
		staticGas: GasVeryLow,
	},
	DUP1: {
		staticGas: GasVeryLow,
	},
	DUP2: {
		staticGas: GasVeryLow,
	},
	DUP3: {
		staticGas: GasVeryLow,
	},
	DUP4: {
		staticGas: GasVeryLow,
	},
	DUP5: {
		staticGas: GasVeryLow,
	},
	DUP6: {
		staticGas: GasVeryLow,
	},
	DUP7: {
		staticGas: GasVeryLow,
	},
	DUP8: {
		staticGas: GasVeryLow,
	},
	DUP9: {
		staticGas: GasVeryLow,
	},
	DUP10: {
		staticGas: GasVeryLow,
	},
	DUP11: {
		staticGas: GasVeryLow,
	},
	DUP12: {
		staticGas: GasVeryLow,
	},
	DUP13: {
		staticGas: GasVeryLow,
	},
	DUP14: {
		staticGas: GasVeryLow,
	},
	DUP15: {
		staticGas: GasVeryLow,
	},
	DUP16: {
		staticGas: GasVeryLow,
	},
	SWAP1: {
		staticGas: GasVeryLow,
	},
	SWAP2: {
		staticGas: GasVeryLow,
	},
	SWAP3: {
		staticGas: GasVeryLow,
	},
	SWAP4: {
		staticGas: GasVeryLow,
	},
	SWAP5: {
		staticGas: GasVeryLow,
	},
	SWAP6: {
		staticGas: GasVeryLow,
	},
	SWAP7: {
		staticGas: GasVeryLow,
	},
	SWAP8: {
		staticGas: GasVeryLow,
	},
	SWAP9: {
		staticGas: GasVeryLow,
	},
	SWAP10: {
		staticGas: GasVeryLow,
	},
	SWAP11: {
		staticGas: GasVeryLow,
	},
	SWAP12: {
		staticGas: GasVeryLow,
	},
	SWAP13: {
		staticGas: GasVeryLow,
	},
	SWAP14: {
		staticGas: GasVeryLow,
	},
	SWAP15: {
		staticGas: GasVeryLow,
	},
	SWAP16: {
		staticGas: GasVeryLow,
	},
	LOG0: {
		dynamicGas: makeGasLog(0),
		memSize:    memoryLog,
	},
	LOG1: {
		dynamicGas: makeGasLog(1),
		memSize:    memoryLog,
	},
	LOG2: {
		dynamicGas: makeGasLog(2),
		memSize:    memoryLog,
	},
	LOG3: {
		dynamicGas: makeGasLog(3),
		memSize:    memoryLog,
	},
	LOG4: {
		dynamicGas: makeGasLog(4),
		memSize:    memoryLog,
	},
	CREATE: {
		staticGas:  CreateGas,
		dynamicGas: gasCreate,
		memSize:    memoryCreate,
	},
	CALL: {
		staticGas:  CallGas,
		dynamicGas: gasCall,
		memSize:    memoryCall,
	},
	CALLCODE: {
		staticGas:  CallGas,
		dynamicGas: gasCallCode,
		memSize:    memoryCall,
	},
	RETURN: {
		dynamicGas: gasReturn,
		memSize:    memoryReturn,
	},
	SELFDESTRUCT: {
		dynamicGas: gasSelfdestruct,
	},
	// Ethereum Homestead instruction set
	DELEGATECALL: {
		dynamicGas: gasDelegateCall,
		memSize:    memoryDelegateCall,
	},
	// Ethereum Byzantium instruction set
	STATICCALL: {
		dynamicGas: gasStaticCall,
		memSize:    memoryStaticCall,
	},
	RETURNDATASIZE: {
		staticGas: GasBase,
	},
	RETURNDATACOPY: {
		dynamicGas: gasReturnDataCopy,
		memSize:    memoryReturnDataCopy,
	},
	REVERT: {
		dynamicGas: gasRevert,
		memSize:    memoryRevert,
	},
	// Ethereum Constantinople instruction set
	SHL: {
		staticGas: GasVeryLow,
	},
	SHR: {
		staticGas: GasVeryLow,
	},
	SAR: {
		staticGas: GasVeryLow,
	},
	EXTCODEHASH: {
		staticGas: GasExtcodeHash,
	},
	CREATE2: {
		dynamicGas: gasCreate2,
		memSize:    memoryCreate2,
	},
}
