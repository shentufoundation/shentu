package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/certikfoundation/shentu/x/cert/types"
)

// NewDecodeStore unmarshals the KVPair's Value to the corresponding type of cert module.
func NewDecodeStore(cdc codec.Marshaler) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.CertifiersStoreKey()):
			var certifierA, certifierB types.Certifier
			cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &certifierA)
			cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &certifierB)
			return fmt.Sprintf("%v\n%v", certifierA, certifierB)

		case bytes.Equal(kvA.Key[:1], types.ValidatorsStoreKey()):
			var validatorA, validatorB types.Validator
			cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &validatorA)
			cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &validatorB)
			return fmt.Sprintf("%v\n%v", validatorA, validatorB)

		case bytes.Equal(kvA.Key[:1], types.CertifierAliasesStoreKey()):
			var certifierA, certifierB types.Certifier
			cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &certifierA)
			cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &certifierB)
			return fmt.Sprintf("%v\n%v", certifierA, certifierB)

		default:
			panic(fmt.Sprintf("invalid %s key prefix %X", types.ModuleName, kvA.Key[:1]))
		}
	}
}
