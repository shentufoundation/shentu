package simulation

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/certikfoundation/shentu/v2/x/shield/types"
)

// NewDecodeStore unmarshals the KVPair's Value to the corresponding type of shield module.
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.ShieldAdminKey):
			return fmt.Sprintf("%v\n%v", sdk.AccAddress(kvA.Value), sdk.AccAddress(kvA.Value))

		case bytes.Equal(kvA.Key[:1], types.TotalCollateralKey),
			bytes.Equal(kvA.Key[:1], types.TotalShieldKey),
			bytes.Equal(kvA.Key[:1], types.TotalClaimedKey):
			var totalA, totalB sdk.IntProto
			cdc.MustUnmarshalLengthPrefixed(kvA.Value, &totalA)
			cdc.MustUnmarshalLengthPrefixed(kvB.Value, &totalB)
			return fmt.Sprintf("%v\n%v", totalA, totalB)

		case bytes.Equal(kvA.Key[:1], types.ServiceFeesKey),
			bytes.Equal(kvA.Key[:1], types.RemainingServiceFeesKey):
			var serviceFeesA, serviceFeesB types.ServiceFees
			cdc.MustUnmarshalLengthPrefixed(kvA.Value, &serviceFeesA)
			cdc.MustUnmarshalLengthPrefixed(kvB.Value, &serviceFeesB)
			return fmt.Sprintf("%v\n%v", serviceFeesA, serviceFeesB)

		case bytes.Equal(kvA.Key[:1], types.PoolKey):
			var poolA, poolB types.Pool
			cdc.MustUnmarshalLengthPrefixed(kvA.Value, &poolA)
			cdc.MustUnmarshalLengthPrefixed(kvB.Value, &poolB)
			return fmt.Sprintf("%v\n%v", poolA, poolB)

		case bytes.Equal(kvA.Key[:1], types.NextPoolIDKey),
			bytes.Equal(kvA.Key[:1], types.NextPurchaseIDKey):
			idA := binary.LittleEndian.Uint64(kvA.Value)
			idB := binary.LittleEndian.Uint64(kvB.Value)
			return fmt.Sprintf("%v\n%v", idA, idB)

		case bytes.Equal(kvA.Key[:1], types.ProviderKey):
			var providerA, providerB types.Provider
			cdc.MustUnmarshalLengthPrefixed(kvA.Value, &providerA)
			cdc.MustUnmarshalLengthPrefixed(kvB.Value, &providerB)
			return fmt.Sprintf("%v\n%v", providerA, providerB)

		case bytes.Equal(kvA.Key[:1], types.PurchaseKey):
			var sPA, spB types.Purchase
			cdc.MustUnmarshalLengthPrefixed(kvA.Value, &sPA)
			cdc.MustUnmarshalLengthPrefixed(kvB.Value, &spB)
			return fmt.Sprintf("%v\n%v", sPA, spB)

		case bytes.Equal(kvA.Key[:1], types.BlockServiceFeesKey):
			var blockFeesA, blockFeesB types.ServiceFees
			cdc.MustUnmarshalLengthPrefixed(kvA.Value, &blockFeesA)
			cdc.MustUnmarshalLengthPrefixed(kvB.Value, &blockFeesB)
			return fmt.Sprintf("%v\n%v", blockFeesA, blockFeesB)

		case bytes.Equal(kvA.Key[:1], types.DonationPoolKey):
			var donationPoolA, donationPoolB types.Pool
			cdc.MustUnmarshalLengthPrefixed(kvA.Value, &donationPoolA)
			cdc.MustUnmarshalLengthPrefixed(kvB.Value, &donationPoolB)
			return fmt.Sprintf("%v\n%v", donationPoolA, donationPoolB)

		default:
			panic(fmt.Sprintf("invalid %s key prefix %X", types.ModuleName, kvA.Key[:1]))
		}
	}
}
