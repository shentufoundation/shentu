package vm

import (
	"math/big"

	. "github.com/hyperledger/burrow/execution/evm"
)

// calcMemSize calculates memory size depending on the operation.
// default means the memory function is not specified, so it returns 0
func calcMemSize(fn uint, stack *Stack) (uint64, bool) {
	switch fn {
	case memorySha3, memoryLog, memoryReturn, memoryRevert:
		return mem64(0, 1, stack)
	case memoryCallDataCopy, memoryCodeCopy, memoryReturnDataCopy:
		return mem64(0, 2, stack)
	case memoryExtCodeCopy:
		return mem64(1, 3, stack)
	case memoryMLoad, memoryMStore:
		return memUint64(0, 32, stack)
	case memoryMStore8:
		return memUint64(0, 1, stack)
	case memoryCreate, memoryCreate2:
		return mem64(1, 2, stack)
	case memoryCall:
		return mem64Comp(5, 6, 3, 4, stack)
	case memoryDelegateCall, memoryStaticCall:
		return mem64Comp(4, 5, 2, 3, stack)
	default:
		return 0, false
	}
}

// calcMemSize64 takes two values and passes them to calcMemSize64WithUint.
func calcMemSize64(off, l *big.Int) (uint64, bool) {
	if !l.IsUint64() {
		return 0, true
	}
	return calcMemSize64WithUint(off, l.Uint64())
}

// calcMemSize64WithUint calculates the required memory size, and returns
// the size and whether the result overflowed uint64
// Identical to calcMemSize64, but length is a uint64
func calcMemSize64WithUint(off *big.Int, length64 uint64) (uint64, bool) {
	// if length is zero, memsize is always zero, regardless of offset
	if length64 == 0 {
		return 0, false
	}
	// Check that offset doesn't overflow
	if !off.IsUint64() {
		return 0, true
	}
	offset64 := off.Uint64()
	val := offset64 + length64
	// if value < either of it's parts, then it overflowed
	return val, val < offset64
}

// mem64 peeks a-th and b-th value from the top, and passes them to the generic function calcMemSize64.
func mem64(a, b int, stack *Stack) (uint64, bool) {
	return calcMemSize64(Get(stack, a), Get(stack, b))
}

// memUint64 peeks a-th value from the top, and takes a length b, and passes them to the generic function calcMemSize64WithUint.
func memUint64(a int, b uint64, stack *Stack) (uint64, bool) {
	return calcMemSize64WithUint(Get(stack, a), b)
}

// mem64Comp calculates two memories, and returns the bigger one.
func mem64Comp(a, b, c, d int, stack *Stack) (uint64, bool) {
	x, overflow := calcMemSize64(Get(stack, a), Get(stack, b))
	if overflow {
		return 0, true
	}
	y, overflow := calcMemSize64(Get(stack, c), Get(stack, d))
	if overflow {
		return 0, true
	}
	if x > y {
		return x, false
	}
	return y, false
}
