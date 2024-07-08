package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

type Keeper struct {
	keeper.Keeper
	storeKey storetypes.StoreKey
}

func NewKeeper(cdc codec.BinaryCodec, key storetypes.StoreKey, ak types.AccountKeeper, bk types.BankKeeper, authority string) *Keeper {
	return &Keeper{
		Keeper:   *keeper.NewKeeper(cdc, key, ak, bk, authority),
		storeKey: key,
	}
}
