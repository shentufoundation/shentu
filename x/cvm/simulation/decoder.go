package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/certikfoundation/shentu/x/cvm/types"
)

// DecodeStore unmarshals the KVPair's value to the corresponding type of cvm module.
func NewDecodeStore(cdc codec.Marshaler) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.StorageStoreKeyPrefix):
			var valueA, valueB []byte
			valueA = kvA.Value
			valueB = kvB.Value
			return fmt.Sprintf("%b\n%b", valueA, valueB)

		case bytes.Equal(kvA.Key[:1], types.BlockHashStoreKeyPrefix):
			var hashA, hashB []byte
			hashA = kvA.Value
			hashB = kvB.Value
			return fmt.Sprintf("%b\n%b", hashA, hashB)

		case bytes.Equal(kvA.Key[:1], types.CodeStoreKeyPrefix):
			var evmCodeA, evmCodeB []byte
			evmCodeA = kvA.Value
			evmCodeB = kvB.Value
			return fmt.Sprintf("%b\n%b", evmCodeA, evmCodeB)

		case bytes.Equal(kvA.Key[:1], types.AbiStoreKeyPrefix):
			var abiA, abiB []byte
			abiA = kvA.Value
			abiB = kvB.Value
			return fmt.Sprintf("%b\n%b", abiA, abiB)

		case bytes.Equal(kvA.Key[:1], types.MetaHashStoreKeyPrefix):
			var metadataA, metadataB string
			metadataA = string(kvA.Value)
			metadataB = string(kvB.Value)
			return fmt.Sprintf("%s\n%s", metadataA, metadataB)

		case bytes.Equal(kvA.Key[:1], types.AddressMetaHashStoreKeyPrefix):
			var metadataA, metadataB types.ContractMetas
			cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &metadataA)
			cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &metadataB)
			return fmt.Sprintf("%v\n%v", metadataA, metadataB)

		default:
			panic(fmt.Sprintf("invalid %s key prefix %X", types.ModuleName, kvA.Key[:1]))
		}
	}
}
