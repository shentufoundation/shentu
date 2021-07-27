// Package keeper specifies the keeper for the cert module.
package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/cert/types"
)

// Keeper manages certifier & security council related logics.
type Keeper struct {
	storeKey       sdk.StoreKey
	cdc            codec.BinaryMarshaler
	slashingKeeper types.SlashingKeeper
	stakingKeeper  types.StakingKeeper
}

// NewKeeper creates a new instance of the certifier keeper.
func NewKeeper(cdc codec.BinaryMarshaler, storeKey sdk.StoreKey, slashingKeeper types.SlashingKeeper, stakingKeeper types.StakingKeeper) Keeper {
	return Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		slashingKeeper: slashingKeeper,
		stakingKeeper:  stakingKeeper,
	}
}
