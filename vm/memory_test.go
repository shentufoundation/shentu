package vm

import (
	"math"
	"math/big"
	"testing"

	"github.com/hyperledger/burrow/execution/errors"

	"github.com/stretchr/testify/assert"

	"github.com/hyperledger/burrow/execution/evm"
)

func NewFakeStack() *evm.Stack {
	var dummyGas uint64 = math.MaxUint64 // set it to MAX(uint64) since we don't care
	stack := evm.NewStack(new(errors.Maybe), 1024, 100000, big.NewInt(int64(dummyGas)))
	return stack
}

func TestCalcMemSize(t *testing.T) {
	op := memoryCreate
	st := NewFakeStack()
	st.Push64(1)
	st.Push64(13)
	st.Push64(14)
	v, of := calcMemSize(uint(op), st)
	v2, of2 := mem64(1, 2, st)
	assert.Equal(t, v, v2)
	assert.Equal(t, of, of2)
}

func TestMem64(t *testing.T) {
	st := NewFakeStack()
	st.Push64(18446744073709551615)
	st.Push64(1)
	r, of := mem64(0, 1, st)
	assert.True(t, of)
	assert.Equal(t, uint64(0), r)
}

func TestMemUint64(t *testing.T) {
	st := NewFakeStack()
	st.Push64(312)

	r, of := memUint64(0, 32, st)
	assert.False(t, of)
	assert.Equal(t, uint64(312+32), r, "equals offset + length")

	r, of = memUint64(0, 1, st)
	assert.False(t, of)
	assert.Equal(t, uint64(312+1), r, "equals offset + length")

	r, of = memUint64(0, 0, st)
	assert.False(t, of)
	assert.Equal(t, uint64(0), r, "equals 0 if length == 0")

	st.Push64(18446744073709551615)
	st.Push64(1)
	r, of = mem64(0, 1, st)
	assert.True(t, of)
	assert.Equal(t, uint64(0), r)
}

func TestMem64Comp(t *testing.T) {
	st := NewFakeStack()
	st.Push64(123)
	st.Push64(321)
	st.Push64(1000)
	st.Push64(20000)

	r, of := mem64Comp(0, 1, 2, 3, st)
	assert.False(t, of)
	assert.Equal(t, uint64(21000), r)

	r, of = mem64Comp(2, 3, 0, 1, st)
	assert.False(t, of)
	assert.Equal(t, uint64(21000), r)

	r, of = mem64Comp(1, 0, 3, 2, st)
	assert.False(t, of)
	assert.Equal(t, uint64(21000), r)
}
