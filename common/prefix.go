package common

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

const (
	// Bech32MainPrefix is the common prefix of all prefixes.
	Bech32MainPrefix = "shentu"

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

// PrefixToCertik convert shentu prefix address to certik prefix address
func PrefixToCertik(address string) (string, error) {
	hrp, data, err := bech32.DecodeAndConvert(address)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(hrp, "shentu") {
		return "", fmt.Errorf("invalid address:%s", address)
	}

	newhrp := strings.Replace(hrp, "shentu", "certik", 1)
	certikAddr, err := bech32.ConvertAndEncode(newhrp, data)
	if err != nil {
		return "", err
	}
	return certikAddr, nil
}

// PrefixToShentu convert certik prefix address to shentu prefix address
func PrefixToShentu(address string) (string, error) {
	hrp, data, err := bech32.DecodeAndConvert(address)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(hrp, "certik") {
		return "", fmt.Errorf("invalid address:%s", address)
	}
	newhrp := strings.Replace(hrp, "certik", "shentu", 1)
	shentuAddr, err := bech32.ConvertAndEncode(newhrp, data)
	if err != nil {
		return "", err
	}
	return shentuAddr, nil
}
