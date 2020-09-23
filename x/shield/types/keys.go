package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName is the name of this module
	ModuleName = "shield"

	// RouterKey is used to route messages.
	RouterKey = ModuleName

	// StoreKey is the prefix under which we store this module's data.
	StoreKey = ModuleName

	// QuerierRoute is used to handle abci_query requests.
	QuerierRoute = ModuleName

	// DefaultParamspace is the default name for parameter store.
	DefaultParamspace = ModuleName
)

var (
	PoolKey = []byte{0x0}
)

// gets the key for the pool with address
func GetPoolKey(accAddr sdk.AccAddress) []byte {
	return append(PoolKey, accAddr.Bytes()...)
}
