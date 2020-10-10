package simulation

import (
	"bytes"
	"fmt"

	tmkv "github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding staking type
func DecodeStore(cdc *codec.Codec, kvA, kvB tmkv.Pair) string {
	switch {
	case bytes.Equal(kvA.Key[:1], types.LastTotalPowerKey):
		var powerA, powerB sdk.Int
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &powerA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &powerB)
		return fmt.Sprintf("%v\n%v", powerA, powerB)

	case bytes.Equal(kvA.Key[:1], types.ValidatorsKey):
		var validatorA, validatorB types.Validator
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &validatorA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &validatorB)
		return fmt.Sprintf("%v\n%v", validatorA, validatorB)

	case bytes.Equal(kvA.Key[:1], types.LastValidatorPowerKey),
		bytes.Equal(kvA.Key[:1], types.ValidatorsByConsAddrKey),
		bytes.Equal(kvA.Key[:1], types.ValidatorsByPowerIndexKey):
		return fmt.Sprintf("%v\n%v", sdk.ValAddress(kvA.Value), sdk.ValAddress(kvB.Value))

	case bytes.Equal(kvA.Key[:1], types.DelegationKey):
		var delegationA, delegationB types.Delegation
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &delegationA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &delegationB)
		return fmt.Sprintf("%v\n%v", delegationA, delegationB)

	case bytes.Equal(kvA.Key[:1], types.UnbondingDelegationKey),
		bytes.Equal(kvA.Key[:1], types.UnbondingDelegationByValIndexKey):
		var ubdA, ubdB types.UnbondingDelegation
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &ubdA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &ubdB)
		return fmt.Sprintf("%v\n%v", ubdA, ubdB)

	case bytes.Equal(kvA.Key[:1], types.RedelegationKey),
		bytes.Equal(kvA.Key[:1], types.RedelegationByValSrcIndexKey):
		var redA, redB types.Redelegation
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &redA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &redB)
		return fmt.Sprintf("%v\n%v", redA, redB)

	case bytes.Equal(kvA.Key[:1], types.UnbondingQueueKey):
		var ubdqA, ubdqB []types.DVPair
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &ubdqA)
		fmt.Printf(">> DEBUG UnbondingQueueKey: A %v\n", ubdqA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &ubdqB)
		fmt.Printf(">> DEBUG UnbondingQueueKey: B %v\n", ubdqB)
		return fmt.Sprintf("%v\n%v", ubdqA, ubdqB)

	case bytes.Equal(kvA.Key[:1], types.RedelegationQueueKey):
		var redqA, redqB []types.DVVTriplet
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &redqA)
		fmt.Printf(">> DEBUG RedelegationQueueKey: A %v\n", redqA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &redqB)
		fmt.Printf(">> DEBUG RedelegationQueueKey: B %v\n", redqB)
		return fmt.Sprintf("%v\n%v", redqA, redqB)

	default:
		panic(fmt.Sprintf("invalid staking key prefix %X", kvA.Key[:1]))
	}
}
