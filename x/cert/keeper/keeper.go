// Package keeper specifies the keeper for the cert module.
package keeper

import (
	"fmt"

	"cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

// Keeper manages certifier & security council related logics.
type Keeper struct {
	storeService   store.KVStoreService
	cdc            codec.BinaryCodec
	slashingKeeper types.SlashingKeeper
	stakingKeeper  types.StakingKeeper
	authority      string
}

// NewKeeper creates a new instance of the certifier keeper.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	slashingKeeper types.SlashingKeeper,
	stakingKeeper types.StakingKeeper,
	authority string,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:            cdc,
		storeService:   storeService,
		slashingKeeper: slashingKeeper,
		stakingKeeper:  stakingKeeper,
		authority:      authority,
	}
}
