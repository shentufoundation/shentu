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
	case bytes.Equal(kvA.Key[:1], types.PoolKey):
		var poolA, poolB types.Pool
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &poolA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &poolB)
		return fmt.Sprintf("%v\n%v", poolA, poolB)

	case bytes.Equal(kvA.Key[:1], types.ShieldAdminKey):
		return fmt.Sprintf("%v\n%v", sdk.AccAddress(kvA.Value), sdk.AccAddress(kvA.Value))

	case bytes.Equal(kvA.Key[:1], types.NextPoolIDKey):
		var poolIDA, poolIDB uint64
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &poolIDA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &poolIDB)
		return fmt.Sprintf("%v\n%v", poolIDA, poolIDB)

	case bytes.Equal(kvA.Key[:1], types.PurchaseListKey):
		var purchaseA, purchaseB types.Purchase
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &purchaseA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &purchaseB)
		return fmt.Sprintf("%v\n%v", purchaseA, purchaseB)

	case bytes.Equal(kvA.Key[:1], types.ReimbursementKey):
		var reimbursementA, reimbursementB types.Reimbursement
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &reimbursementA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &reimbursementB)
		return fmt.Sprintf("%v\n%v", reimbursementA, reimbursementB)

	case bytes.Equal(kvA.Key[:1], types.CollateralKey):
		var collateralA, collateralB types.Collateral
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &collateralA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &collateralB)
		return fmt.Sprintf("%v\n%v", collateralA, collateralB)

	case bytes.Equal(kvA.Key[:1], types.ProviderKey):
		var providerA, providerB types.Provider
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &providerA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &providerB)
		return fmt.Sprintf("%v\n%v", providerA, providerB)

	case bytes.Equal(kvA.Key[:1], types.PurchaseQueueKey):
		var purchasesA, purchasesB []types.Purchase
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &purchasesA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &purchasesB)
		return fmt.Sprintf("%v\n%v", purchasesA, purchasesB)

	default:
		panic(fmt.Sprintf("invalid %s key prefix %X", types.ModuleName, kvA.Key[:1]))
	}
}
