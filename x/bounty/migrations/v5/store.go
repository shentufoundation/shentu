package v5

import (
	"fmt"

	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/common"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

// MigrateStore migrates the bounty module state from version 5 to version 6.
// It copies the existing complexity_fee value into the new per-type fields:
// complexity_fee_rocq and complexity_fee_lean.
func MigrateStore(ctx sdk.Context, storeService corestoretypes.KVStoreService, cdc codec.BinaryCodec) error {
	store := storeService.OpenKVStore(ctx)

	// Get the params key bytes
	paramsKey := collections.NewPrefix(61)

	// Read existing params from raw store
	paramsBz, err := store.Get(paramsKey)
	if err != nil {
		return fmt.Errorf("failed to get params from store: %w", err)
	}

	ctx.Logger().Info("migrating bounty params v5->v6", "data_length", len(paramsBz))

	var oldParam Params
	err = cdc.Unmarshal(paramsBz, &oldParam)
	if err != nil {
		return fmt.Errorf("failed to unmarshal params (len=%d): %w", len(paramsBz), err)
	}

	// Copy the old complexity_fee into both per-type fee fields
	newParam := types.Params{
		MinGrant:              oldParam.MinGrant,
		MinDeposit:            oldParam.MinDeposit,
		TheoremMaxProofPeriod: oldParam.TheoremMaxProofPeriod,
		ProofMaxLockPeriod:    oldParam.ProofMaxLockPeriod,
		ComplexityFee:         oldParam.ComplexityFee,
		MaxComplexity:         oldParam.MaxComplexity,
		ComplexityFeeRocq:     oldParam.ComplexityFee,
		ComplexityFeeLean:     sdk.NewCoin(common.MicroCTKDenom, sdkmath.NewInt(800)),
	}

	// Marshal and save updated params
	bz, err := cdc.Marshal(&newParam)
	if err != nil {
		return err
	}

	return store.Set(paramsKey, bz)
}
