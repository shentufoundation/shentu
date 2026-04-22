package v4

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

// MigrateStore migrates the bounty module state from version 4 to version 5.
// It reads params from raw KVStore to handle protobuf schema changes.
// Protobuf automatically handles removed fields (CheckerRate) and new fields get default values.
func MigrateStore(ctx sdk.Context, storeService corestoretypes.KVStoreService, cdc codec.BinaryCodec) error {
	store := storeService.OpenKVStore(ctx)

	// Get the params key bytes
	paramsKey := collections.NewPrefix(61)

	// Read existing params from raw store
	// Protobuf will automatically ignore the old CheckerRate field
	paramsBz, err := store.Get(paramsKey)
	if err != nil {
		return fmt.Errorf("failed to get params from store: %w", err)
	}

	ctx.Logger().Info("migrating bounty params", "data_length", len(paramsBz))

	var oldParam Params
	err = cdc.Unmarshal(paramsBz, &oldParam)
	if err != nil {
		return fmt.Errorf("failed to unmarshal old params (len=%d): %w", len(paramsBz), err)
	}

	// Create new params with the new fields set
	// ComplexityFee and MaxComplexity are new fields added in this version
	newParam := types.Params{
		MinGrant:              oldParam.MinGrant,
		MinDeposit:            oldParam.MinDeposit,
		TheoremMaxProofPeriod: oldParam.TheoremMaxProofPeriod,
		ProofMaxLockPeriod:    oldParam.ProofMaxLockPeriod,
		ComplexityFee:         sdk.NewCoin(common.MicroCTKDenom, sdkmath.NewInt(10000)),
		MaxComplexity:         int64(1000000),
	}

	// Marshal and save updated params
	bz, err := cdc.Marshal(&newParam)
	if err != nil {
		return err
	}

	return store.Set(paramsKey, bz)
}
