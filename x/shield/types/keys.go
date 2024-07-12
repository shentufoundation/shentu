package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of this module.
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
	RemainingServiceFeesKey = []byte{0x06}
	ProviderKey             = []byte{0x0C}
	BlockServiceFeesKey     = []byte{0x12}
)

func GetBlockServiceFeesKey() []byte {
	return BlockServiceFeesKey
}

func GetRemainingServiceFeesKey() []byte {
	return RemainingServiceFeesKey
}

// GetProviderKey gets the key for the delegator's tracker.
func GetProviderKey(addr sdk.AccAddress) []byte {
	return append(ProviderKey, addr...)
}
