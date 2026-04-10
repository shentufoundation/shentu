package simulation

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

// NewDecodeStore returns a decoder function for the cert module's KV store.
func NewDecodeStore(cdc codec.BinaryCodec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.CertifiersStoreKey()):
			var certifierA, certifierB types.Certifier
			cdc.MustUnmarshalLengthPrefixed(kvA.Value, &certifierA)
			cdc.MustUnmarshalLengthPrefixed(kvB.Value, &certifierB)
			return fmt.Sprintf("%v\n%v", certifierA, certifierB)

		case bytes.Equal(kvA.Key[:1], types.CertificatesStoreKey()):
			var certA, certB types.Certificate
			cdc.MustUnmarshal(kvA.Value, &certA)
			cdc.MustUnmarshal(kvB.Value, &certB)
			return fmt.Sprintf("%v\n%v", certA, certB)

		case bytes.Equal(kvA.Key, types.NextCertificateIDStoreKey()):
			idA := binary.LittleEndian.Uint64(kvA.Value)
			idB := binary.LittleEndian.Uint64(kvB.Value)
			return fmt.Sprintf("NextCertificateID: %d\n%d", idA, idB)

		default:
			panic(fmt.Sprintf("invalid %s key prefix %X", types.ModuleName, kvA.Key[:1]))
		}
	}
}
