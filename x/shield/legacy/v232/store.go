package v232

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/x/shield/types"
)

// MigrateStore performs in-place store migrations from v2.3.1/v2.3.2 to v2.4.0.
// The migration includes:
//
// - Setting the MinCommissionRate param in the paramstore
func MigrateStore(ctx sdk.Context, paramstore types.ParamSubspace) {
	migrateParamsStore(ctx, paramstore)
}

func migrateParamsStore(ctx sdk.Context, paramstore types.ParamSubspace) {
	ctx.Logger().Info("Adding Additional Shield Params..")
	blockRewardParams := types.DefaultDistributionParams()
	paramstore.Set(ctx, types.ParamStoreKeyDistribution, &blockRewardParams)
}
