package keeper

import (
	"github.com/certikfoundation/shentu/v2/x/auth/types"
)

type Keeper struct {
	ak types.AccountKeeper
	ck types.CertKeeper
}

func NewKeeper(ak types.AccountKeeper, ck types.CertKeeper) Keeper {
	return Keeper{
		ak: ak,
		ck: ck,
	}
}
