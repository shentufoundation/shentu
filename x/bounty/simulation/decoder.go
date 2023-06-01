package simulation

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

// BytesToUint64List trans byte array to an uint64
func BytesToUint64List(list []byte) []uint64 {
	buf := bytes.NewBuffer(list)
	r64 := make([]uint64, (len(list)+7)/8)
	binary.Read(buf, binary.LittleEndian, &r64)
	return r64
}

// NewDecodeStore unmarshals the KVPair's Value to the corresponding type of bounty module.
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.ProgramsKey):
			var programA, programB types.Program
			cdc.MustUnmarshalLengthPrefixed(kvA.Value, &programA)
			cdc.MustUnmarshalLengthPrefixed(kvB.Value, &programB)
			return fmt.Sprintf("%v\n%v", programA, programB)

		case bytes.Equal(kvA.Key[:1], types.NextProgramIDKey),
			bytes.Equal(kvA.Key[:1], types.NextFindingIDKey):
			idA := binary.LittleEndian.Uint64(kvA.Value)
			idB := binary.LittleEndian.Uint64(kvB.Value)
			return fmt.Sprintf("%v\n%v", idA, idB)

		case bytes.Equal(kvA.Key[:1], types.FindingKey):
			var findingA, findingB types.Finding
			cdc.MustUnmarshalLengthPrefixed(kvA.Value, &findingA)
			cdc.MustUnmarshalLengthPrefixed(kvB.Value, &findingB)
			return fmt.Sprintf("%v\n%v", findingA, findingB)

		case bytes.Equal(kvA.Key[:1], types.ProgramIDFindingListKey):
			listA := BytesToUint64List(kvA.Value)
			listB := BytesToUint64List(kvB.Value)
			return fmt.Sprintf("%v\n%v", listA, listB)

		default:
			panic(fmt.Sprintf("invalid %s key prefix %X", types.ModuleName, kvA.Key[:1]))
		}
	}
}
