package types

import (
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func concat(bytes ...[]byte) []byte {
	a := []byte{}
	for _, b := range bytes {
		a = append(a, b...)
	}
	return a
}

var (
	// certifierStoreKeyPrefix is the prefix of certifier kv-store keys.
	certifierStoreKeyPrefix = []byte{0x0}

	// validatorStoreKeyPrefix is the prefix of certified validator kv-store keys.
	validatorStoreKeyPrefix = []byte{0x1}

	// certifierAliasStoreKeyPrefix is the prefix of certifier alias kv-store keys.
	certifierAliasStoreKeyPrefix = []byte{0x7}
)

// CertifierStoreKey returns the kv-store key for the certifier registration.
func CertifierStoreKey(certifier sdk.AccAddress) []byte {
	return concat(certifierStoreKeyPrefix, certifier.Bytes())
}

// CertifiersStoreKey returns the kv-store key for accessing all current certifiers in the security council.
func CertifiersStoreKey() []byte {
	return certifierStoreKeyPrefix
}

// CertifierAliasStoreKey returns the kv-store key for the certifier alias.
func CertifierAliasStoreKey(alias string) []byte {
	return concat(certifierAliasStoreKeyPrefix, []byte(alias))
}

// CertifierAliasesStoreKey returns the kv-store key for accessing aliases of all current certifiers.
func CertifierAliasesStoreKey() []byte {
	return certifierAliasStoreKeyPrefix
}

// ValidatorStoreKey returns the kv-store key for the validator node certification.
func ValidatorStoreKey(validator cryptotypes.PubKey) []byte {
	return concat(validatorStoreKeyPrefix, validator.Bytes())
}

// ValidatorsStoreKey returns the kv-store key for accessing all validator node certifications.
func ValidatorsStoreKey() []byte {
	return validatorStoreKeyPrefix
}
