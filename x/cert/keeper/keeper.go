// Package keeper specifies the keeper for the cert module.
package keeper

import (
	"fmt"

	"cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper manages certifier & certificate related logics.
type Keeper struct {
	storeService store.KVStoreService
	cdc          codec.BinaryCodec
	authority    string
}

// NewKeeper creates a new instance of the certifier keeper.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	authority string,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		authority:    authority,
	}
}
