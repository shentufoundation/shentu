package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
	sk       types.StakingKeeper
}

// NewKeeper creates a slashing keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, sk types.StakingKeeper) Keeper {
	return Keeper{
		storeKey: key,
		cdc:      cdc,
		sk:       sk,
	}
}
