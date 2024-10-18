package keeper

import (
	"cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/shentufoundation/shentu/v2/x/auth/types"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	storeService store.KVStoreService
	ak           types.AccountKeeper
	ck           types.CertKeeper
}

func NewKeeper(cdc codec.BinaryCodec, storeService store.KVStoreService, ak types.AccountKeeper, ck types.CertKeeper) Keeper {
	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		ak:           ak,
		ck:           ck,
	}
}
