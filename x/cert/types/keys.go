package types

import (
	"encoding/binary"

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

	// platformStoreKeyPrefix is the prefix of validator host platform kv-store keys.
	platformStoreKeyPrefix = []byte{0x2}

	// certificateStoreKeyPrefix is the prefix of certificate kv-store keys.
	certificateStoreKeyPrefix = []byte{0x5}

	// libraryStoreKeyPrefix is the prefix of library kv-store keys.
	libraryStoreKeyPrefix = []byte{0x6}

	// certifierAliasStoreKeyPrefix is the prefix of certifier alias kv-store keys.
	certifierAliasStoreKeyPrefix = []byte{0x7}

	// nextCertificateIDKeyPrefix is the prefix of the next certificate ID to assign.
	nextCertificateIDKeyPrefix = []byte{0x08}
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

// CertificateStoreKey returns the kv-store key for accessing a given certificate (ID).
func CertificateStoreKey(id uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, id)
	return concat(certificateStoreKeyPrefix, bz)
}

// CertificatesStoreKey returns the kv-store key for accessing all certificates.
func CertificatesStoreKey() []byte {
	return certificateStoreKeyPrefix
}

// LibraryStoreKey returns the kv-store key for accessing certificate library address.
func LibraryStoreKey(library sdk.AccAddress) []byte {
	return concat(libraryStoreKeyPrefix, library.Bytes())
}

// LibrariesStoreKey returns the kv-store key for accessing all certificate library addresses.
func LibrariesStoreKey() []byte {
	return libraryStoreKeyPrefix
}

// PlatformStoreKey returns the kv-store key for the validator host platform certificate.
func PlatformStoreKey(validator cryptotypes.PubKey) []byte {
	return append(platformStoreKeyPrefix, validator.Bytes()...)
}

// PlatformsStoreKey returns the kv-store key for accessing all platform certificates.
func PlatformsStoreKey() []byte {
	return platformStoreKeyPrefix
}

// NextCertificateIDStoreKey returns the kv-store key for next certificate ID to assign.
func NextCertificateIDStoreKey() []byte {
	return nextCertificateIDKeyPrefix
}
