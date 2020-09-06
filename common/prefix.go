package common

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// Bech32MainPrefix is the common prefix of all prefixes.
	Bech32MainPrefix = "certik"

	// Bech32PrefixAccAddr is the prefix of account addresses.
	Bech32PrefixAccAddr = Bech32MainPrefix
	// Bech32PrefixAccPub is the prefix of account public keys.
	Bech32PrefixAccPub = Bech32MainPrefix + sdk.PrefixPublic
	// Bech32PrefixValAddr is the prefix of validator operator addresses.
	Bech32PrefixValAddr = Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixOperator
	// Bech32PrefixValPub is the prefix of validator operator public keys.
	Bech32PrefixValPub = Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	// Bech32PrefixConsAddr is the prefix of consensus node addresses.
	Bech32PrefixConsAddr = Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixConsensus
	// Bech32PrefixConsPub is the prefix of consensus node public keys.
	Bech32PrefixConsPub = Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
)
