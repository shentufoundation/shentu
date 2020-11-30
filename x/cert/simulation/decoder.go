package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/certikfoundation/shentu/x/cert/internal/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding type of cert module.
func DecodeStore(cdc *codec.Codec, kvA, kvB tmkv.Pair) string {
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

	case bytes.Equal(kvA.Key[:1], types.PlatformsStoreKey()):
		var platformA, platformB types.Platform
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &platformA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &platformB)
		return fmt.Sprintf("%v\n%v", platformA, platformB)

	case bytes.Equal(kvA.Key[:1], types.CertificatesStoreKey()):
		var certificateA, certificateB types.Certificate
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &certificateA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &certificateB)
		return fmt.Sprintf("%v\n%v", certificateA, certificateB)

	case bytes.Equal(kvA.Key[:1], types.LibrariesStoreKey()):
		var libraryA, libraryB types.Library
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &libraryA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &libraryB)
		return fmt.Sprintf("%v\n%v", libraryA, libraryB)

	case bytes.Equal(kvA.Key[:1], types.CertifierAliasesStoreKey()):
		var certifierA, certifierB types.Certifier
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &certifierA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &certifierB)
		return fmt.Sprintf("%v\n%v", certifierA, certifierB)

	default:
		panic(fmt.Sprintf("invalid %s key prefix %X", types.ModuleName, kvA.Key[:1]))
	}
}
