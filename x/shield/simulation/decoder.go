package simulation

import (
	"bytes"
	"fmt"

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

	case bytes.Equal(kvA.Key[:1], types.TotalCollateralKey):
		var totalCollateralA, totalCollateralB sdk.Int
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &totalCollateralA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &totalCollateralB)
		return fmt.Sprintf("%v\n%v", totalCollateralA, totalCollateralB)

	case bytes.Equal(kvA.Key[:1], types.TotalShieldKey):
		var totalShieldA, totalShieldB sdk.Int
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &totalShieldA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &totalShieldB)
		return fmt.Sprintf("%v\n%v", totalShieldA, totalShieldB)

	case bytes.Equal(kvA.Key[:1], types.TotalLockedKey):
		var totalLockedA, totalLockedB sdk.Int
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &totalLockedA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &totalLockedB)
		return fmt.Sprintf("%v\n%v", totalLockedA, totalLockedB)

	case bytes.Equal(kvA.Key[:1], types.ServiceFeesKey):
		var serviceFeesA, serviceFeesB types.MixedDecCoins
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &serviceFeesA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &serviceFeesB)
		return fmt.Sprintf("%v\n%v", serviceFeesA, serviceFeesB)

	case bytes.Equal(kvA.Key[:1], types.PoolKey):
		var poolA, poolB types.Pool
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &poolA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &poolB)
		return fmt.Sprintf("%v\n%v", poolA, poolB)

	case bytes.Equal(kvA.Key[:1], types.NextPoolIDKey):
		var poolIDA, poolIDB uint64
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &poolIDA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &poolIDB)
		return fmt.Sprintf("%v\n%v", poolIDA, poolIDB)

	case bytes.Equal(kvA.Key[:1], types.NextPurchaseIDKey):
		var purchaseIDA, purchaseIDB uint64
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &purchaseIDA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &purchaseIDB)
		return fmt.Sprintf("%v\n%v", purchaseIDA, purchaseIDB)

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

	default:
		panic(fmt.Sprintf("invalid %s key prefix %X", types.ModuleName, kvA.Key[:1]))
	}
}
