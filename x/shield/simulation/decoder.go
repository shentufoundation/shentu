package simulation

import (
	"bytes"
	"encoding/binary"
	"fmt"

	gogotypes "github.com/gogo/protobuf/types"

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

		case bytes.Equal(kvA.Key[:1], types.FeesKey),
			bytes.Equal(kvA.Key[:1], types.RemainingFeesKey):
			var feesA, feesB types.Fees
			cdc.MustUnmarshalLengthPrefixed(kvA.Value, &feesA)
			cdc.MustUnmarshalLengthPrefixed(kvB.Value, &feesB)
			return fmt.Sprintf("%v\n%v", feesA, feesB)

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

		case bytes.Equal(kvA.Key[:1], types.PurchaseListKey):
			var purchaseA, purchaseB types.PurchaseList
			cdc.MustUnmarshalLengthPrefixed(kvA.Value, &purchaseA)
			cdc.MustUnmarshalLengthPrefixed(kvB.Value, &purchaseB)
			return fmt.Sprintf("%v\n%v", purchaseA, purchaseB)

		case bytes.Equal(kvA.Key[:1], types.ProviderKey):
			var providerA, providerB types.Provider
			cdc.MustUnmarshalLengthPrefixed(kvA.Value, &providerA)
			cdc.MustUnmarshalLengthPrefixed(kvB.Value, &providerB)
			return fmt.Sprintf("%v\n%v", providerA, providerB)

		case bytes.Equal(kvA.Key[:1], types.LastUpdateTimeKey):
			var timeA, timeB gogotypes.Timestamp
			cdc.MustUnmarshalLengthPrefixed(kvA.Value, &timeA)
			cdc.MustUnmarshalLengthPrefixed(kvB.Value, &timeB)
			return fmt.Sprintf("%v\n%v", timeA, timeB)

		case bytes.Equal(kvA.Key[:1], types.StakeForShieldKey):
			var sPA, spB types.ShieldStaking
			cdc.MustUnmarshalLengthPrefixed(kvA.Value, &sPA)
			cdc.MustUnmarshalLengthPrefixed(kvB.Value, &spB)
			return fmt.Sprintf("%v\n%v", sPA, spB)

		case bytes.Equal(kvA.Key[:1], types.BlockFeesKey):
			var blockFeesA, blockFeesB types.Fees
			cdc.MustUnmarshalLengthPrefixed(kvA.Value, &blockFeesA)
			cdc.MustUnmarshalLengthPrefixed(kvB.Value, &blockFeesB)
			return fmt.Sprintf("%v\n%v", blockFeesA, blockFeesB)

		case bytes.Equal(kvA.Key[:1], types.OriginalStakingKey):
			var rateA, rateB sdk.IntProto
			cdc.MustUnmarshalLengthPrefixed(kvA.Value, &rateA)
			cdc.MustUnmarshalLengthPrefixed(kvB.Value, &rateB)
			return fmt.Sprintf("%v\n%v", rateA, rateB)

		default:
			panic(fmt.Sprintf("invalid %s key prefix %X", types.ModuleName, kvA.Key[:1]))
		}
	}
}
