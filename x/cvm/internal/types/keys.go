package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/hyperledger/burrow/acm/acmstate"
	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/crypto"
)

const (
	// ModuleName is the name of the staking module.
	ModuleName = "cvm"

	// StoreKey is the string store representation.
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the staking module.
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the staking module.
	RouterKey = ModuleName

	// DefaultParamspace is the default name for parameter store.
	DefaultParamspace = ModuleName
)

var (
	// StorageStoreKeyPrefix is the prefix of CVM storage kv-store keys.
	StorageStoreKeyPrefix = []byte{0x00}

	// BlockHashStoreKeyPrefix is the prefix of block hash kv-store keys.
	BlockHashStoreKeyPrefix = []byte{0x01}

	// CodeStoreKeyPrefix is the prefix of code kv-store keys.
	CodeStoreKeyPrefix = []byte{0x02}

	// AbiStoreKeyPrefix is the prefix of code ABI kv-store keys.
	AbiStoreKeyPrefix = []byte{0x03}

	// MetaHashStoreKeyPrefix is the prefix of contract metadata hash kv-store keys.
	MetaHashStoreKeyPrefix = []byte{0x04}

	// AddressMetaHashStoreKeyPrefix is the prefix of contract metadata hash kv-store keys.
	AddressMetaHashStoreKeyPrefix = []byte{0x5}
)

// StorageStoreKey returns the kv-store key for the contract's storage key.
func StorageStoreKey(addr crypto.Address, key binary.Word256) []byte {
	return append(append(StorageStoreKeyPrefix, addr.Bytes()...), key.Bytes()...)
}

// BlockHashStoreKey returns the kv-store key for the chain's block hashes.
func BlockHashStoreKey(height int64) []byte {
	return append(BlockHashStoreKeyPrefix, sdk.NewInt(height).BigInt().Bytes()...)
}

// CodeStoreKey returns the kv-store key for the contract's storage key.
func CodeStoreKey(addr crypto.Address) []byte {
	return append(CodeStoreKeyPrefix, addr.Bytes()...)
}

// AbiStoreKey returns the kv-store key for the contract's code ABI key.
func AbiStoreKey(addr crypto.Address) []byte {
	return append(AbiStoreKeyPrefix, addr.Bytes()...)
}

// metaHashStoreKey returns the kv-store key for the metahash key.
func MetaHashStoreKey(metahash acmstate.MetadataHash) []byte {
	return append(MetaHashStoreKeyPrefix, metahash.Bytes()...)
}

// AddressMetaStoreKey returns the kv-store key for the address-metahash key.
func AddressMetaStoreKey(addr crypto.Address) []byte {
	return append(AddressMetaHashStoreKeyPrefix, addr.Bytes()...)
}
