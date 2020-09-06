package vm

import (
	"math"
	"math/big"

	. "github.com/hyperledger/burrow/binary"
	. "github.com/hyperledger/burrow/execution/evm"
)

// Get returns the n-th element from the stack
func Get(st *Stack, n int) *big.Int {
	st.Dup(n + 1)
	tmp := st.Pop()
	return new(big.Int).SetBytes(tmp[:])
}

// GetWord256 is similar with Get, but returns with Word256 type (default type in the Stack).
func GetWord256(st *Stack, n int) Word256 {
	st.Dup(n + 1)
	return st.Pop()
}

// SafeSub Subtracts y from x, and sets the overflow flag.
// of==true means an overflow has occurred.
func SafeSub(x, y uint64) (uint64, bool) {
	return x - y, x < y
}

// SafeAdd add x and y and sets the overflow flag.
// of==true means an overflow has occurred.
func SafeAdd(x, y uint64) (uint64, bool) {
	return x + y, y > math.MaxUint64-x
}

// SafeMul multiplies x by y and sets the overflow flag.
// of==true means an overflow has occurred.
func SafeMul(x, y uint64) (uint64, bool) {
	if x == 0 || y == 0 {
		return 0, false
	}
	return x * y, y > math.MaxUint64/x
}

// bigToUint64 returns the integer casted to a uint64 and returns whether it
// overflowed in the process.
func bigToUint64(v *big.Int) (uint64, bool) {
	return v.Uint64(), !v.IsUint64()
}

// toWordSize returns the word size given a byte size.
func toWordSize(size uint64) uint64 {
	if size > math.MaxUint64-31 {
		return math.MaxUint64/32 + 1
	}
	return (size + 31) / 32
}

// Min returns the smaller of two uint64 integers.
func Min(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}

// GetRefund returns the refund counter of the vm.
func (vm *CVM) GetRefund() uint64 {
	return vm.refund
}
