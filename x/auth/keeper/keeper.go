package keeper

import (
	"github.com/certikfoundation/shentu/x/auth/types"
)

type Keeper struct {
	ck types.CertKeeper
}

func NewKeeper(ck types.CertKeeper) Keeper {
	return Keeper{
		ck: ck,
	}
}
