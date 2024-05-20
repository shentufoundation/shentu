package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/shentufoundation/shentu/v2/x/auth/types"
)

type Keeper struct {
	cdc codec.BinaryCodec
	key storetypes.StoreKey
	ak  types.AccountKeeper
	ck  types.CertKeeper
}

func NewKeeper(cdc codec.BinaryCodec, key storetypes.StoreKey, ak types.AccountKeeper, ck types.CertKeeper) Keeper {
	return Keeper{
		cdc: cdc,
		key: key,
		ak:  ak,
		ck:  ck,
	}
}
