package vm

import (
	"math"
	"math/big"
	"reflect"
	"testing"

	. "github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/execution/errors"
	. "github.com/hyperledger/burrow/execution/evm"
)

func TestStack_Get(t *testing.T) {
	type fields struct {
		Stack *Stack
	}
	type args struct {
		n int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *big.Int
	}{
		{"normal", fields{NewContentStack()}, args{1}, big.NewInt(10)},
		{"error", fields{NewContentStack()}, args{2}, big.NewInt(20)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := tt.fields.Stack
			if got := Get(st, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stack.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func NewContentStack() *Stack {
	err := new(errors.Maybe)
	gaz := uint64(math.MaxUint64)
	st := NewStack(err, 0, 0, &gaz)
	st.Push64(10)
	st.Push64(20)
	return st
}

func TestStack_GetWord256(t *testing.T) {
	type fields struct {
		Stack *Stack
	}
	type args struct {
		n int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Word256
	}{
		{"normal", fields{NewContentStack()}, args{1}, LeftPadWord256([]byte{0x0a})},
		{"error", fields{NewContentStack()}, args{2}, LeftPadWord256([]byte{0x14})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := tt.fields.Stack
			if got := GetWord256(st, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stack.GetWord256() = %v, want %v", got, tt.want)
			}
		})
	}
}

func newZeroBytes() []byte {
	var b = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	return b
}

func newWeirdBytes() []byte {
	var b = []byte{0x12, 0x34, 0x45, 0x56, 0x67, 0x48, 0xef, 0xff, 0xfa, 0xda, 0xba, 0x12, 0x11, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	return b
}

func TestSafeSub(t *testing.T) {
	type args struct {
		x uint64
		y uint64
	}
	tests := []struct {
		name  string
		args  args
		want  uint64
		want1 bool
	}{

		{"normal", args{32, 0}, 32, false},
		{"zero", args{1, 1}, 0, false},
		{"max", args{18446744073709551615, 18446744073709551614}, 1, false},
		{"overflow", args{1, 18446744073709551615}, 2, true},
		{"underflow", args{0, 1}, 18446744073709551615, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := SafeSub(tt.args.x, tt.args.y)
			if got != tt.want {
				t.Errorf("SafeSub() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SafeSub() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func newBool() *bool {
	t := false
	return &t
}

func TestSafeAdd(t *testing.T) {
	type args struct {
		x  uint64
		y  uint64
		of *bool
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{"one", args{1, 1, newBool()}, 2},
		{"overflow", args{18446744073709551615, 5, newBool()}, 4},    // overflow is set to True
		{"no overlfow", args{18446744073709551615, 1, newBool()}, 0}, // overflow is set to True
		{"normal", args{11, 11, newBool()}, 22},
		{"zero", args{123, 0, newBool()}, 123},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := SafeAdd(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("SafeAdd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSafeMul(t *testing.T) {
	type args struct {
		x  uint64
		y  uint64
		of *bool
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{"one", args{1, 1, newBool()}, 1},
		{"overflow", args{9223372036854775807, 5, newBool()}, 9223372036854775803}, // overflow is set to True
		{"no overlfow", args{9223372036854775807, 1, newBool()}, 9223372036854775807},
		{"normal", args{11, 11, newBool()}, 121},
		{"zero", args{123, 0, newBool()}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := SafeMul(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("SafeMul() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bigToUint64(t *testing.T) {
	type args struct {
		v *big.Int
	}
	tests := []struct {
		name  string
		args  args
		want  uint64
		want1 bool
	}{
		{"normal", args{big.NewInt(32)}, 32, false},
		{"zero", args{big.NewInt(0)}, 0, false},
		{"overflow", args{big.NewInt(-1)}, 1, true},
		{"max", args{big.NewInt(9223372036854775807)}, 9223372036854775807, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := bigToUint64(tt.args.v)
			if got != tt.want {
				t.Errorf("bigToUint64() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("bigToUint64() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_toWordSize(t *testing.T) {
	type args struct {
		size uint64
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{"exact", args{32}, 1},
		{"one over", args{33}, 2},
		{"one less", args{31}, 1},
		{"zero", args{0}, 0},
		{"one", args{1}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toWordSize(tt.args.size); got != tt.want {
				t.Errorf("toWordSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
