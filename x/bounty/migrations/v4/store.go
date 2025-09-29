package v4

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func MigrateStore(ctx sdk.Context, storeService store.KVStoreService, cdc codec.BinaryCodec) error {
	sb := collections.NewSchemaBuilder(storeService)

	params := collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc))
	oldParam, err := params.Get(ctx)
	if err != nil {
		return err
	}
	newParam := types.Params{
		MinGrant:              oldParam.MinGrant,
		MinDeposit:            oldParam.MinDeposit,
		TheoremMaxProofPeriod: oldParam.TheoremMaxProofPeriod,
		ProofMaxLockPeriod:    oldParam.ProofMaxLockPeriod,
		ComplexityFee:         sdk.NewCoin("uctk", sdkmath.NewInt(10000)),
	}
	err = params.Set(ctx, newParam)
	if err != nil {
		return err
	}

	return nil
}
