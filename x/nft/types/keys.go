package types

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of this module
	ModuleName = "nft"

	// RouterKey is used to route messages.
	RouterKey = ModuleName

	// StoreKey is the prefix under which we store this module's data.
	StoreKey = ModuleName

	// QuerierRoute is used to handle abci_query requests.
	QuerierRoute = ModuleName
)

var (
	AdminKeyPrefix             = []byte{0x10}
	CertificateStoreKeyPrefix  = []byte{0x11}
	NextCertificateIDKeyPrefix = []byte{0x12}
)

func AdminKey(addr sdk.AccAddress) []byte {
	return append(AdminKeyPrefix, addr...)
}

// CertificateStoreKey returns the kv-store key for accessing a given certificate (ID).
func CertificateStoreKey(id uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, id)
	return append(CertificateStoreKeyPrefix, bz...)
}

// CertificatesStoreKey returns the kv-store key for accessing all certificates.
func CertificatesStoreKey() []byte {
	return CertificateStoreKeyPrefix
}

// NextCertificateIDStoreKey returns the kv-store key for next certificate ID to assign.
func NextCertificateIDStoreKey() []byte {
	return NextCertificateIDKeyPrefix
}
