package simulation

import (
	"bytes"
	"fmt"
	"time"

	tmkv "github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding type of shield module.
func DecodeStore(cdc *codec.Codec, kvA, kvB tmkv.Pair) string {
	switch {
	case bytes.Equal(kvA.Key[:1], types.ShieldAdminKey):
		return fmt.Sprintf("%v\n%v", sdk.AccAddress(kvA.Value), sdk.AccAddress(kvA.Value))

	case bytes.Equal(kvA.Key[:1], types.TotalCollateralKey),
		bytes.Equal(kvA.Key[:1], types.TotalShieldKey),
		bytes.Equal(kvA.Key[:1], types.TotalClaimedKey):
		var totalA, totalB sdk.Int
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &totalA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &totalB)
		return fmt.Sprintf("%v\n%v", totalA, totalB)

	case bytes.Equal(kvA.Key[:1], types.ServiceFeesKey),
		bytes.Equal(kvA.Key[:1], types.RemainingServiceFeesKey):
		var serviceFeesA, serviceFeesB types.MixedDecCoins
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &serviceFeesA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &serviceFeesB)
		return fmt.Sprintf("%v\n%v", serviceFeesA, serviceFeesB)

	case bytes.Equal(kvA.Key[:1], types.PoolKey):
		var poolA, poolB types.Pool
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &poolA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &poolB)
		return fmt.Sprintf("%v\n%v", poolA, poolB)

	case bytes.Equal(kvA.Key[:1], types.NextPoolIDKey),
		bytes.Equal(kvA.Key[:1], types.NextPurchaseIDKey):
		var idA, idB uint64
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &idA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &idB)
		return fmt.Sprintf("%v\n%v", idA, idB)

	case bytes.Equal(kvA.Key[:1], types.PurchaseListKey):
		var purchaseA, purchaseB types.PurchaseList
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &purchaseA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &purchaseB)
		return fmt.Sprintf("%v\n%v", purchaseA, purchaseB)

	case bytes.Equal(kvA.Key[:1], types.ProviderKey):
		var providerA, providerB types.Provider
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &providerA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &providerB)
		return fmt.Sprintf("%v\n%v", providerA, providerB)

	case bytes.Equal(kvA.Key[:1], types.LastUpdateTimeKey):
		var timeA, timeB time.Time
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &timeA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &timeB)
		return fmt.Sprintf("%v\n%v", timeA, timeB)

	case bytes.Equal(kvA.Key[:1], types.StakeForShieldKey):
		var sPA, spB types.ShieldStaking
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &sPA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &spB)
		return fmt.Sprintf("%v\n%v", sPA, spB)

	case bytes.Equal(kvA.Key[:1], types.BlockServiceFeesKey):
		var blockFeesA, blockFeesB types.MixedDecCoins
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &blockFeesA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &blockFeesB)
		return fmt.Sprintf("%v\n%v", blockFeesA, blockFeesB)

	case bytes.Equal(kvA.Key[:1], types.OriginalStakingKey):
		var rateA, rateB sdk.Int
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &rateA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &rateB)
		return fmt.Sprintf("%v\n%v", rateA, rateB)

	default:
		panic(fmt.Sprintf("invalid %s key prefix %X", types.ModuleName, kvA.Key[:1]))
	}
}
